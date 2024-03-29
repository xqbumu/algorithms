package main

import (
	"algorithms/pkg/restdb"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/sijms/go-ora/v2"
	_ "modernc.org/sqlite"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var app *restdb.App

func main() {

	confPath := flag.String("c", "restdb.json", "configration file path")
	flag.Parse()
	confBytes, err := os.ReadFile(*confPath)
	if err != nil {
		log.Fatal(err)
	}

	app, err = restdb.NewApp(confBytes)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", app.Handler)

	if app.Web.HttpAddr != "" {
		srv := &http.Server{
			Addr:         app.Web.HttpAddr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		go func() {
			fmt.Println(fmt.Sprint("Listening on http://", app.Web.HttpAddr, "/"))
			log.Fatal(srv.ListenAndServe())
		}()
	}

	if app.Web.HttpsAddr != "" {
		srvs := &http.Server{
			Addr:         app.Web.HttpsAddr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		go func() {
			fmt.Println(fmt.Sprint("Listening on https://", app.Web.HttpsAddr, "/"))
			log.Fatal(srvs.ListenAndServeTLS(app.Web.CertFile, app.Web.KeyFile))
		}()
	}

	restdb.Hook(nil)
}
