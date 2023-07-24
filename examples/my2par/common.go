package main

import (
	"algorithms/pkg/querykit"
	"reflect"

	"github.com/segmentio/parquet-go"
	"gorm.io/gorm"
)

func ExportTable(
	db *gorm.DB, table string, w *parquet.Writer, page querykit.Pagination[any],
	rv reflect.Value, ptrVals []any, txCallback func(tx *gorm.DB),
) (int, error) {
	tx := db.
		Scopes(func(db *gorm.DB) *gorm.DB {
			return db.Offset(page.GetOffset()).Limit(page.GetLimit()).Order(page.GetSort())
		}).
		Table(table)

	if txCallback != nil {
		txCallback(tx)
	}

	rows, err := tx.Rows()
	if err != nil {
		panic(err)
	}

	cnt := 0
	for rows.Next() {
		cnt += 1
		if err := rows.Scan(ptrVals...); err != nil {
			return cnt, err
		}
		err = w.Write(rv.Interface())
		if err != nil {
			return cnt, err
		}
	}

	return cnt, rows.Err()
}
