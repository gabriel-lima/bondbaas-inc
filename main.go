package main

import (
	"bondbaas/handlers"
	"bondbaas/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var db *sql.DB

func main() {
	db = storage.InitDB()
	defer db.Close()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/admin/tables", adminTablesHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	table, id := extractPath(r.URL.Path)

	if tableNameIsEmpty(table, w) || tableNameIsAReservedWord(table, w) || tableNameNotFound(table, w) {
		return
	}

	var err error
	gateway := storage.TableGateway{DB: db, Table: table}

	handler := handlers.TableHandler{
		Response:     w,
		TableGateway: gateway,
	}

	if r.Method == "GET" {
		handler.Get(id)
	}

	if r.Method == "POST" {
		var js map[string]interface{}
		var body []byte
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			responseMalformed(w, err)
			return
		}
		err = json.Unmarshal([]byte(body), &js)
		if err != nil {
			responseMalformed(w, err)
			return
		}

		handler.Create(js)
	}

	if r.Method == "PUT" {
		var js map[string]interface{}
		var body []byte
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			responseMalformed(w, err)
			return
		}
		err = json.Unmarshal([]byte(body), &js)
		if err != nil {
			responseMalformed(w, err)
			return
		}

		handler.Update(id, js)
	}

	if r.Method == "DELETE" {
		if id == 0 {
			http.Error(w, "Id must to be an integer", 422)
			return
		}

		handler.Delete(id)
	}
}

func extractPath(path string) (tableName string, ID int) {
	paths := strings.Split(path, "/")
	deleteEmptyPaths(paths)

	tableName = strings.ToLower(strings.Join(paths[1:], ""))
	ID, _ = strconv.Atoi(strings.Join(paths[2:], ""))

	return tableName, ID
}

func deleteEmptyPaths(s []string) {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
}

func tableNameIsEmpty(tableName string, w http.ResponseWriter) bool {
	if tableName == "" {
		http.Error(w, "Please inform a table", 404)
		return true
	}
	return false
}

func tableNameIsAReservedWord(tableName string, w http.ResponseWriter) bool {
	if tableName == "admin" {
		http.Error(w, "admin is a reserved name.", 422)
		return true
	}
	return false
}

func tableNameNotFound(tableName string, w http.ResponseWriter) bool {
	gateway := storage.AdminGateway{DB: db}

	hasTable, err := gateway.HasTable(tableName)
	if err != nil {
		responseInternalError(w, err)
		return true
	}

	if !hasTable {
		http.Error(w, fmt.Sprintf("Table %s not found.", tableName), 404)
		return true
	}

	return false
}

/// Create a table schema
/*
POST
/tables
{
	"name": "products",
	"columns": [
		{
			"name": "brand",
			"type": "VARCHAR(50)",
			"constraint": "NULL"
		}
	]
}
*/
func adminTablesHandler(w http.ResponseWriter, r *http.Request) {
	gateway := storage.AdminGateway{DB: db}
	handler := handlers.AdminHandler{Response: w, AdminGateway: gateway}

	if r.Method == "GET" {
		handler.Get()
	}

	if r.Method == "POST" {
		table := storage.Table{}
		err := json.NewDecoder(r.Body).Decode(&table)

		if err != nil {
			responseMalformed(w, err)
			return
		}

		handler.Create(table)
	}
}

func responseMalformed(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 422)
	log.Println(err)
}

func responseInternalError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 500)
	log.Println(err)
}
