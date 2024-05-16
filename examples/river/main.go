package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"golang.org/x/exp/constraints"
)

var (
	// Maximum number of connections for the connection pool. This is the same
	// default that pgxpool uses (the larger of 4 or number of CPUs), but made a
	// variable here so that we can reference it from the test suite and not
	// rely on implicit knowledge of pgxpool implementation details that could
	// change in the future. If changing this value, also change the number of
	// databases to create in `testdbman`.
	dbPoolMaxConns = int32(max(4, runtime.NumCPU())) //nolint:gochecknoglobals
)

type SortArgs struct {
	// Strings is a slice of strings to sort.
	Strings []string `json:"strings"`
}

func (SortArgs) Kind() string { return "sort" }

type SortWorker struct {
	river.WorkerDefaults[SortArgs]
}

func (w *SortWorker) Work(ctx context.Context, job *river.Job[SortArgs]) error {
	sort.Strings(job.Args.Strings)
	fmt.Printf("Sorted strings: %+v\n", job.Args.Strings)
	return nil
}

// ValOrDefault returns the given value if it's non-zero, and otherwise returns
// the default.
func ValOrDefault[T constraints.Integer | string](val, defaultVal T) T {
	var zero T
	if val != zero {
		return val
	}
	return defaultVal
}

func DatabaseConfig(databaseName string) *pgxpool.Config {
	databaseURL := ValOrDefault(os.Getenv("RIVER_DATABASE_URL"), "postgres:///river_testdb?sslmode=disable")

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		panic(fmt.Sprintf("error parsing database URL: %v", err))
	}
	config.MaxConns = dbPoolMaxConns
	config.ConnConfig.ConnectTimeout = 10 * time.Second
	config.ConnConfig.Database = databaseName
	config.ConnConfig.RuntimeParams["timezone"] = "UTC"
	return config
}

// TruncateRiverTables truncates River tables in the target database. This is
// for test cleanup and should obviously only be used in tests.
func TruncateRiverTables(ctx context.Context, pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tables := []string{"river_job", "river_leader"}

	for _, table := range tables {
		if _, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s;", table)); err != nil {
			return fmt.Errorf("error truncating %q: %w", table, err)
		}
	}

	return nil
}

// WaitTimeout returns a duration broadly appropriate for waiting on an expected
// event in a test, and which is used for `TestSignal.WaitOrTimeout` and
// `riverinternaltest.WaitOrTimeout`. It's main purpose is to allow a little
// extra leeway in GitHub Actions where we occasionally seem to observe subpar
// performance which leads to timeouts and test intermittency, while still
// keeping a tight a timeout for local test runs where this is never a problem.
func WaitTimeout() time.Duration {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return 10 * time.Second
	}

	return 3 * time.Second
}

// Wait on the given subscription channel for numJobs. Times out with a panic if
// jobs take too long to be received.
func waitForNJobs(subscribeChan <-chan *river.Event, numJobs int) {
	var (
		timeout  = WaitTimeout()
		deadline = time.Now().Add(timeout)
		events   = make([]*river.Event, 0, numJobs)
	)

	for {
		select {
		case event := <-subscribeChan:
			events = append(events, event)

			if len(events) >= numJobs {
				return
			}

		case <-time.After(time.Until(deadline)):
			panic(fmt.Sprintf("WaitOrTimeout timed out after waiting %s (received %d job(s), wanted %d)",
				timeout, len(events), numJobs))
		}
	}
}

// Example_insertAndWork demonstrates how to register job workers, start a
// client, and insert a job on it to be worked.
func main() {
	ctx := context.Background()

	dbPool, err := pgxpool.NewWithConfig(ctx, DatabaseConfig("river_testdb_example"))
	if err != nil {
		panic(err)
	}
	defer dbPool.Close()

	// Required for the purpose of this test, but not necessary in real usage.
	if err := TruncateRiverTables(ctx, dbPool); err != nil {
		panic(err)
	}

	workers := river.NewWorkers()
	river.AddWorker(workers, &SortWorker{})

	riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})),
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		Workers: workers,
	})
	if err != nil {
		panic(err)
	}

	// Out of example scope, but used to make wait until a job is worked.
	subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobCompleted)
	defer subscribeCancel()

	if err := riverClient.Start(ctx); err != nil {
		panic(err)
	}

	// Start a transaction to insert a job. It's also possible to insert a job
	// outside a transaction, but this usage is recommended to ensure that all
	// data a job needs to run is available by the time it starts. Because of
	// snapshot visibility guarantees across transactions, the job will not be
	// worked until the transaction has committed.
	tx, err := dbPool.Begin(ctx)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback(ctx)

	_, err = riverClient.InsertTx(ctx, tx, SortArgs{
		Strings: []string{
			"whale", "tiger", "bear",
		},
	}, nil)
	if err != nil {
		panic(err)
	}

	if err := tx.Commit(ctx); err != nil {
		panic(err)
	}

	waitForNJobs(subscribeChan, 1)

	if err := riverClient.Stop(ctx); err != nil {
		panic(err)
	}

}
