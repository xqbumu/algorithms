package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type App struct {
	Web       *Web                 `json:"web"`
	Databases map[string]*Database `json:"databases"`
	Scripts   map[string]*Script   `json:"scripts"`
	Tables    map[string]*Table    `json:"tables"`
	Tokens    map[string]*[]Access `json:"tokens"`
	Opt       map[string]any       `json:"opt"`
}

func NewApp(confBytes []byte) (*App, error) {
	var app *App
	err := json.Unmarshal(confBytes, &app)
	return app, err
}

type Web struct {
	HttpAddr  string `json:"http_addr"`
	HttpsAddr string `json:"https_addr"`
	CertFile  string `json:"cert_file"`
	KeyFile   string `json:"key_file"`
	Cors      bool   `json:"cors"`
}

type Database struct {
	Type string `json:"type"`
	Url  string `json:"url"`
	Conn *sql.DB
}

func (db *Database) Open() (*sql.DB, error) {
	if db.Conn != nil {
		return db.Conn, nil
	}
	var err error
	db.Conn, err = sql.Open(db.Type, db.Url)
	return db.Conn, err
}

func (db *Database) GetPlaceHolder(index int) string {
	if db.IsPg() {
		return fmt.Sprintf("$%d", index+1)
	} else {
		return "?"
	}
}

func (db *Database) IsPg() bool {
	return db.Type == "pgx"
}

type Access struct {
	Database string `json:"database"`
	Object   string `json:"object"`
	Read     bool   `json:"read"`
	Write    bool   `json:"write"`
	Exec     bool   `json:"exec"`
}

type Statement struct {
	Index  int
	Label  string
	Text   string
	Params []string
	Query  bool
	Export bool
	Script *Script
}

type Script struct {
	Database   string `json:"database"`
	Text       string `json:"text"`
	Path       string `json:"path"`
	PublicExec bool   `json:"public_exec"`
	Statements []*Statement
}

type Table struct {
	Database    string `json:"database"`
	Name        string `json:"name"`
	PublicRead  bool   `json:"public_read"`
	PublicWrite bool   `json:"public_write"`
}