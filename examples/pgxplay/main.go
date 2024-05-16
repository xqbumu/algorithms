package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var (
		ctx       = context.Background()
		empNo     int64
		birthDate time.Time
		firstName string
		gender    byte
	)

	rows, err := conn.Query(
		ctx,
		"select emp_no, birth_date, first_name, gender from employees where emp_no=$1",
		10001,
	)

	// // load extra data types
	// oids := make([]uint32, 0, 8)
	// for _, fieldOpt := range rows.FieldDescriptions() {
	// 	oids = append(oids, fieldOpt.DataTypeOID)
	// }
	// RegisterDataTypesByOID(ctx, conn, oids)

	for rows.Next() {
		if rows.Scan(&empNo, &birthDate, &firstName, &gender); err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println(empNo, birthDate, firstName, gender)
}

func RegisterDataTypesByOID(ctx context.Context, conn *pgx.Conn, oids []uint32) error {
	var typeName string
	for _, typeOID := range oids {
		err := conn.QueryRow(ctx, "SELECT typname FROM pg_type WHERE oid = $1;", typeOID).Scan(&typeName)
		if err != nil {
			return err
		}
		dataType, err := conn.LoadType(ctx, typeName)
		if err != nil {
			return err
		}
		conn.TypeMap().RegisterType(dataType)
	}

	return nil
}

func RegisterDataTypes(ctx context.Context, conn *pgx.Conn) error {
	dataTypeNames := []string{
		"employees_employees_gender_enum",
	}

	for _, typeName := range dataTypeNames {
		dataType, err := conn.LoadType(ctx, typeName)
		if err != nil {
			return err
		}
		conn.TypeMap().RegisterType(dataType)
	}

	return nil
}
