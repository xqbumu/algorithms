package main

import (
	"algorithms/pkg/querykit"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-module/carbon/v2"
	"github.com/segmentio/parquet-go"
	"github.com/segmentio/parquet-go/compress/zstd"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	cmdTable := &cobra.Command{
		Use:   "table",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			years, err := cmd.Flags().GetIntSlice("years")
			if err != nil {
				panic(err)
			}

			for _, year := range years {
				runTable(year)
			}
		},
	}

	cmdTable.Flags().IntSliceP("years", "", []int{carbon.Now().Year()}, "")
	rootCmd.AddCommand(cmdTable)
}

func runTable(year int) {
	table := "scat_report_simple_no_promote"
	pk := []string{"dw_date", "id"}

	date := carbon.Now().SetYear(year)

	dsn = fmt.Sprintf("%s?charset=utf8mb4&parseTime=True&loc=Local", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Create a new Parquet file writer
	fn := fmt.Sprintf("./%s/%s_%d.parquet", outputDir, table, date.Year())
	fnTmp := fmt.Sprintf("%s.%d", fn, carbon.Now().Timestamp())
	w, err := os.Create(fnTmp)
	if err != nil {
		log.Println("Can't create local file", err)
		return
	}
	defer w.Close()

	rv, pVals, err := querykit.TableSchema(db, table)
	if err != nil {
		panic(err)
	}

	schema := parquet.SchemaOf(rv.Interface())
	dest := parquet.NewWriter(
		w,
		parquet.Compression(&zstd.Codec{}),
		schema,
	)
	defer dest.Close()

	page := querykit.Pagination[any]{
		Size: pageSize,
		Page: 1,
		Sort: strings.Join(pk, ","),
	}

	round := 0
	for {
		round += 1
		cnt, err := ExportTable(db, table, dest, page, rv, pVals, func(tx *gorm.DB) {
			tx.Where(
				"dw_date BETWEEN ? AND ?",
				date.StartOfYear().Format("Y-m-d"), date.EndOfYear().Format("Y-m-d"),
			)
		})
		if err != nil {
			panic(err)
		}
		if err := dest.Flush(); err != nil {
			panic(err)
		}
		if rowLimit > 0 && page.Size*page.Page > rowLimit {
			break
		}
		page.Page += 1
		log.Printf("year: %d, round: %d", year, round)
		if cnt < page.Size {
			break
		}
	}

	os.Rename(fnTmp, fn)
}
