package main

import (
	"algorithms/pkg/querykit"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/segmentio/parquet-go"
	"github.com/segmentio/parquet-go/compress/zstd"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := ""
	flag.StringVar(&dsn, "dsn", "", "example: root:@tcp(127.0.0.1:3306)/db_name")
	flag.Parse()

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s?charset=utf8mb4&parseTime=True&loc=Local", dsn)), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	repo := GameRepo{db: db}

	// Create a new Parquet file writer
	w, err := os.Create("Game.parquet")
	if err != nil {
		log.Println("Can't create local file", err)
		return
	}
	defer w.Close()

	pw := parquet.NewGenericWriter[Game](
		w, parquet.Compression(&zstd.Codec{}),
		// parquet.ColumnPageBuffers(parquet.NewFileBufferPool("", "buffers.*")),
	)
	if err != nil {
		panic(err)
	}

	page := querykit.Pagination[Game]{
		Size: 100000,
		Page: 1,
		Sort: "game_id ASC",
	}

	ch := make(chan []Game, 1)
	go func() {
		for {
			p, err := repo.List(page)
			if err != nil {
				panic(err)
			}
			log.Printf("Page: %d", page.Page)
			ch <- p.Rows
			page.Page += 1
			if len(p.Rows) < page.Size {
				close(ch)
				break
			}
		}
	}()

	for {
		rows, ok := <-ch // read data from channel
		if !ok {
			break // channel has been closed
		}
		if _, err := pw.Write(rows); err != nil {
			panic(err)
		}
		if err := pw.Flush(); err != nil {
			panic(err)
		}
	}

	if err = pw.Close(); err != nil {
		log.Println("WriteStop error", err)
		return
	}

	log.Println("Write Finished")
}
