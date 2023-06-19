package main

import (
	"algorithms/pkg/querykit"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/jimsmart/schema"
	"github.com/segmentio/parquet-go"
	"github.com/segmentio/parquet-go/compress/zstd"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := ""
	table := ""
	pk := ""
	limit := -1
	paegSize := 100000
	flag.StringVar(&dsn, "dsn", "", "example: root:@tcp(127.0.0.1:3306)/db_name")
	flag.StringVar(&table, "table", "", "")
	flag.StringVar(&pk, "pk", "", "")
	flag.IntVar(&limit, "limit", -1, "")
	flag.IntVar(&paegSize, "page-szie", 100000, "")
	flag.Parse()

	dsn = fmt.Sprintf("%s?charset=utf8mb4&parseTime=True&loc=Local", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Create a new Parquet file writer
	w, err := os.Create(fmt.Sprintf("./output/%s.parquet", table))
	if err != nil {
		log.Println("Can't create local file", err)
		return
	}
	defer w.Close()

	rv, pVals, err := getTableSchema(db, table)
	if err != nil {
		panic(err)
	}

	dest := parquet.NewWriter(
		w,
		parquet.Compression(&zstd.Codec{}),
		parquet.SchemaOf(rv.Interface()),
	)
	defer dest.Close()

	page := querykit.Pagination[any]{
		Size: paegSize,
		Page: 1,
		Sort: fmt.Sprintf("%s ASC", pk),
	}

	round := 0
	for {
		round += 1
		cnt, err := ExportTable(db, table, dest, page, rv, pVals)
		if err != nil {
			panic(err)
		}
		if err := dest.Flush(); err != nil {
			panic(err)
		}
		page.Page += 1
		log.Printf("round: %d", round)
		if cnt < page.Size {
			break
		}
	}
}

func ExportTable(
	db *gorm.DB, table string, w *parquet.Writer, page querykit.Pagination[any],
	rv reflect.Value, ptrVals []any,
) (int, error) {
	rows, err := db.Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Offset(page.GetOffset()).Limit(page.GetLimit()).Order(page.GetSort())
	}).Table(table).Rows()
	if err != nil {
		panic(err)
	}

	cnt := 0
	for rows.Next() {
		cnt += 1
		if err := rows.Scan(ptrVals...); err != nil {
			return cnt, err
		}
		err := w.Write(rv.Interface())
		if err != nil {
			return cnt, err
		}
	}

	return cnt, rows.Err()
}

func getTableSchema(db *gorm.DB, table string) (reflect.Value, []interface{}, error) {
	conn, err := db.DB()
	if err != nil {
		return reflect.Value{}, nil, err
	}
	ct, err := schema.ColumnTypes(conn, "", table)
	if err != nil {
		return reflect.Value{}, nil, err
	}

	var fields []reflect.StructField

	for i, c := range ct {
		f := reflect.StructField{
			Name: fmt.Sprintf("Field%d", i),
			Tag:  reflect.StructTag(fmt.Sprintf(`parquet:"%s"`, c.Name())),
		}
		tt := reflect.New(c.ScanType()).Interface()
		switch tt.(type) {
		case *int8, *int32, *uint32:
			var t *int32
			f.Type = reflect.TypeOf(t)
		case *int64:
			var t *int64
			f.Type = reflect.TypeOf(t)
		case *sql.NullBool:
			var t *bool
			f.Type = reflect.TypeOf(t)
		case *sql.NullFloat64:
			var t *float64
			f.Type = reflect.TypeOf(t)
		case *sql.NullInt64:
			var t *int64
			f.Type = reflect.TypeOf(t)
		case *sql.NullString:
			var t *string
			f.Type = reflect.TypeOf(t)
		case *sql.NullTime:
			var t *time.Time
			f.Type = reflect.TypeOf(t)
		case *sql.RawBytes:
			var t *string
			f.Type = reflect.TypeOf(t)
		case **interface{}:
			var t *string
			f.Type = reflect.TypeOf(t)
		default:
			panic(fmt.Sprintf("Unknown type for field %s %T", c.Name(), tt))
		}

		fields = append(fields, f)
	}

	typ := reflect.StructOf(fields)

	pVals := make([]interface{}, len(ct))

	rv := reflect.New(typ).Elem()
	for i := 0; i < rv.NumField(); i++ {
		pVals[i] = rv.Field(i).Addr().Interface()
	}

	for i := 0; i < rv.NumField(); i++ {
		pVals[i] = rv.Field(i).Addr().Interface()
	}

	return rv, pVals, nil
}
