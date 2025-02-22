package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type CtxCount struct{}

type Count struct {
	sync.Mutex

	now   time.Time
	sum   int64
	round int32
}

func NewCount() *Count {
	r := &Count{
		now: time.Now(),
	}
	return r
}

func (r *Count) AddVal(ctx context.Context, val int64) (string, error) {
	if _, ok := ctx.Value(CtxCount{}).(*Count); ok {
		panic("counter is not allowed")
	}

	r.round++
	r.sum += val

	ret := strconv.Itoa(int(r.sum))

	return ret, nil
}

func (r *Count) Check() {
	if len(r.String()) == 0 {
		panic("r is empty")
	}
}

func (r *Count) String() string {
	return fmt.Sprintf("round: %d, sum: %d", r.round, r.sum)
}

func OneShot(ctx context.Context) string {
	r := NewCount()
	_, err := r.AddVal(ctx, 1)
	if err != nil {
		panic(err)
	}
	return r.String()
}

func main() {
	var err error
	var r = NewCount()

	ctx := context.TODO()

	// Round: 1
	ret, err := r.AddVal(ctx, 3)
	if err != nil {
		panic(err)
	}
	if len(ret) == 0 {
		panic("ret is empty")
	}

	now := time.Now()
	ts := now.UnixNano()
	ret, err = r.AddVal(ctx, ts+r.sum)
	if err != nil {
		panic(err)
	}
	if len(ret) == 0 {
		panic("ret is empty")
	}

	// Round: 2
	r = NewCount()
	_, err = r.AddVal(ctx, 3)
	if err != nil {
		panic(err)
	}
	r.Check()

	now = time.Now()
	ts = now.UnixNano()
	_, err = r.AddVal(ctx, ts)
	if err != nil {
		panic(err)
	}
	r.Check()
}
