package querykit

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/jimsmart/schema"
	"github.com/segmentio/parquet-go"
)

func Generate(db *sql.DB, table string) (*parquet.Schema, []interface{}, error) {
	ct, err := schema.ColumnTypes(db, "", table)
	if err != nil {
		return nil, nil, err
	}

	ptrValues := make([]interface{}, len(ct))

	var fields []reflect.StructField

	for i, c := range ct {
		f := reflect.StructField{
			Name: fmt.Sprintf("Field%d", i),
			Tag:  reflect.StructTag(fmt.Sprintf(`parquet:"%s"`, c.Name())),
		}
		tt := reflect.New(c.ScanType()).Interface()
		switch tt.(type) {
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
		case *sql.RawBytes:
			var t []byte
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

	v := reflect.New(typ).Elem()
	for i := 0; i < v.NumField(); i++ {
		ptrValues[i] = v.Field(i).Addr().Interface()
	}

	for i := 0; i < v.NumField(); i++ {
		ptrValues[i] = v.Field(i).Addr().Interface()
	}

	return parquet.SchemaOf(v.Interface()), ptrValues, nil
}
