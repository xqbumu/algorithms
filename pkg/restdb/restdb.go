// Package restdb turns any SQL database into a RESTful API.
// Currently supports MySQL, MariaDB, PostgreSQL, and SQLite.
package restdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

var format = "json"

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

func (app App) Handler(w http.ResponseWriter, r *http.Request) {
	if app.Web.Cors {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	}

	if r.Method == "OPTIONS" {
		w.Header().Set("Allow", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	urlPath := r.URL.Path
	urlPrefix := FromPrefix(r.Context())
	if len(urlPrefix) > 0 {
		urlPath = urlPath[len(urlPrefix):]
	}
	urlParts := strings.Split(strings.TrimPrefix(urlPath, "/"), "/")
	databaseId := urlParts[0]
	database := app.Databases[databaseId]
	if database == nil {
		fmt.Fprintf(w, `{"error":"database %v not found"}`, urlParts[0])
		return
	}
	objectId := urlParts[1]

	authHeader := r.Header.Get("authorization")

	methodUpper := strings.ToUpper(r.Method)

	authorized, err := app.authorize(methodUpper, authHeader, databaseId, objectId)
	if !authorized {
		fmt.Fprintf(w, `{"error":"%v"}`, err.Error())
		return
	}

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	var bodyData map[string]any
	json.Unmarshal(body, &bodyData)

	paramValues, _ := url.ParseQuery(r.URL.RawQuery)
	params := valuesToMap(false, paramValues)
	for k, v := range bodyData {
		params[k] = v
	}

	var result any

	if methodUpper == "EXEC" || methodUpper == http.MethodPatch {
		script := app.Scripts[objectId]
		if script == nil {
			fmt.Fprintf(w, `{"error":"script %v not found"}`, objectId)
			return
		}
		script.SQL = strings.TrimSpace(script.SQL)
		script.Path = strings.TrimSpace(script.Path)

		if os.Getenv("env") == "dev" {
			script.built = false
		}

		if !script.built {
			if script.SQL == "" && script.Path == "" {
				fmt.Fprintf(w, `{"error":"script %v is empty"}`, objectId)
				return
			}

			if script.Path != "" {
				f, err := os.ReadFile(script.Path)
				if err != nil {
					fmt.Fprintf(w, `{"error":"%v"}`, err.Error())
					return
				}
				script.SQL = string(f)
			}

			err = BuildStatements(script, database.IsPg())
			if err != nil {
				fmt.Fprintf(w, `{"error":"%v"}`, err.Error())
				return
			}
			app.Scripts[objectId] = script
		}

		sepIndex := strings.LastIndex(r.RemoteAddr, ":")
		clientIP := r.RemoteAddr[0:sepIndex]
		clientIP = strings.ReplaceAll(strings.ReplaceAll(clientIP, "[", ""), "]", "")
		params["__client_ip"] = clientIP

		result, err = runExec(database, script, params)
		if err != nil {
			result = map[string]any{
				"error": err.Error(),
			}
		}
	} else {
		dataId := ""
		if len(urlParts) > 2 {
			dataId = urlParts[2]
		}
		table := app.Tables[objectId]
		if table == nil {
			fmt.Fprintf(w, `{"error":"table %v not found"}`, objectId)
			return
		}
		result, err = runTable(methodUpper, database, table, dataId, params)
		if err != nil {
			result = map[string]any{
				"error": err.Error(),
			}
		}
	}

	jsonData, err := json.Marshal(result)
	jsonString := string(jsonData)
	fmt.Fprintln(w, jsonString)

	if err != nil {
		fmt.Fprintf(w, `{"error":"%v"}`, err.Error())
	}
}

func (app App) authorize(methodUpper string, authHeader string, databaseId string, objectId string) (bool, error) {
	// if object is not found, return false
	// if object is found, check if it is public
	// if object is not public, return true regardless of token
	// if database is not specified in object, the object is shared across all databases
	if methodUpper == "EXEC" || methodUpper == http.MethodPatch {
		script := app.Scripts[objectId]
		if script == nil || (script.Database != "" && script.Database != databaseId) {
			return false, fmt.Errorf("script %v not found", objectId)
		}
		if script.Public {
			return true, nil
		}
	} else {
		table := app.Tables[objectId]
		if table == nil || (table.Database != "" && table.Database != databaseId) {
			return false, fmt.Errorf("table %v not found", objectId)
		}
		if table.PublicRead && methodUpper == http.MethodGet {
			return true, nil
		}
		if table.PublicWrite && (methodUpper == http.MethodPost || methodUpper == http.MethodPut || methodUpper == http.MethodDelete) {
			return true, nil
		}
	}

	// object is not public, check token
	// if token doesn't have any access, return false
	accesses := app.Tokens[authHeader]
	if accesses == nil || len(*accesses) == 0 {
		return false, fmt.Errorf("access denied")
	}

	// when token has access, check if any access is allowed for database and object
	for _, access := range *accesses {
		if access.Database == databaseId && slices.Contains(access.Objects, objectId) {
			switch methodUpper {
			case "EXEC", http.MethodPatch:
				if access.Exec {
					return true, nil
				}
			case http.MethodGet:
				if access.Read {
					return true, nil
				}
			case http.MethodPost, http.MethodPut, http.MethodDelete:
				if access.Write {
					return true, nil
				}
			}
		}
	}
	return false, fmt.Errorf("access token not allowed for database %v and object %v", databaseId, objectId)
}

func runTable(method string, database *Database, table *Table, dataId any, params map[string]any) (any, error) {
	db, err := database.Open()
	if err != nil {
		return nil, err
	}
	sqlSafe(&table.Name)
	switch method {
	case http.MethodGet:
		if dataId == "" {
			where, values, err := mapForSqlWhere(params, database.IsPg())
			if err != nil {
				return nil, err
			}
			return QueryToMap(db, Lower, fmt.Sprintf(`SELECT * FROM %v WHERE TRUE %v`, table.Name, where), values...)
		} else {
			r, err := QueryToMap(db, Lower, fmt.Sprintf(`SELECT * FROM %v WHERE id=%v`, table.Name, database.GetPlaceHolder(0)), dataId)
			if err != nil {
				return nil, err
			}
			if len(r) == 0 {
				return nil, nil
			} else {
				return r[0], nil
			}
		}
	case http.MethodPost:
		// should return the id of the new record?
		qms, keys, values, err := mapForSqlInsert(params, database.IsPg())
		if err != nil {
			return nil, err
		}
		return Exec(db, fmt.Sprintf(`INSERT INTO %v (%v) VALUES (%v)`, table.Name, keys, qms), values...)
	case http.MethodPut:
		set, values, err := mapForSqlUpdate(params, database.IsPg())
		if err != nil {
			return nil, err
		}
		return Exec(db, fmt.Sprintf(`UPDATE %v SET %v WHERE ID=%v`, table.Name, set, database.GetPlaceHolder(len(params))), append(values, dataId)...)
	case http.MethodDelete:
		return Exec(db, fmt.Sprintf(`DELETE FROM %v WHERE ID=%v`, table.Name, database.GetPlaceHolder(0)), dataId)
	}
	return nil, fmt.Errorf("Method %v not supported.", method)
}

func runExec(database *Database, script *Script, params map[string]any) (any, error) {
	db, err := database.Open()
	if err != nil {
		return nil, err
	}
	exportedResults := map[string]any{}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	for _, statement := range script.Statements {
		if len(statement.SQL) == 0 {
			continue
		}

		// double underscore
		scriptParams := ExtractScriptParamsFromMap(params)
		for k, v := range scriptParams {
			statement.SQL = strings.ReplaceAll(statement.SQL, k, v.(string))
		}

		var result any
		sqlParams := []any{}
		for _, param := range statement.Params {
			if val, ok := params[param]; ok {
				sqlParams = append(sqlParams, val)
			} else {
				tx.Rollback()
				return nil, fmt.Errorf("Parameter %v not provided.", param)
			}
		}

		if statement.Query {
			if format == "array" {
				header, data, err := QueryToArray(tx, Lower, statement.SQL, sqlParams...)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				if statement.Export {
					exportedResults[statement.Label] = map[string]any{
						"header": header,
						"data":   data,
					}
				}
			} else {
				result, err = QueryToMap(tx, Lower, statement.SQL, sqlParams...)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				if statement.Export {
					exportedResults[statement.Label] = result
				}
			}
		} else {
			result, err = Exec(tx, statement.SQL, sqlParams...)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			if statement.Export {
				exportedResults[statement.Label] = result
			}
		}

	}

	tx.Commit()
	if len(exportedResults) == 0 {
		return nil, nil
	}
	if len(exportedResults) == 1 && exportedResults["0"] != nil {
		return exportedResults["0"], nil
	}
	return exportedResults, nil
}
