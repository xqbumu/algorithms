package querykit

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/jimsmart/schema"
	"gorm.io/gorm"
)

func TableSchema(db *gorm.DB, table string) (reflect.Value, []interface{}, error) {
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
			switch c.DatabaseTypeName() {
			case "DECIMAL":
				var t *float64
				f.Type = reflect.TypeOf(t)
			default:
				var t *string
				f.Type = reflect.TypeOf(t)
			}
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
