package main

import (
	"bondbaas/handlers"
	"bondbaas/storage"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
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

	http.HandleFunc("/", tableHandler)
	http.HandleFunc("/admin/tables", adminHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), nil))
}

func tableHandler(w http.ResponseWriter, r *http.Request) {
	table, id := extractPath(r.URL.Path)

	if tableNameIsEmpty(table, w) || tableNameIsAReservedWord(table, w) || tableNameNotFound(table, w) {
		return
	}

	gateway := storage.TableGateway{DB: db, Table: table}
	handler := handlers.TableHandler{
		Request:      r,
		Response:     w,
		TableGateway: gateway,
	}

	if r.Method == "GET" {
		handler.Get(id)
	}

	if r.Method == "POST" {
		handler.Create()
	}

	if r.Method == "PUT" {
		handler.Update(id)
	}

	if r.Method == "DELETE" {
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
func adminHandler(w http.ResponseWriter, r *http.Request) {
	gateway := storage.AdminGateway{DB: db}
	handler := handlers.AdminHandler{
		Request:      r,
		Response:     w,
		AdminGateway: gateway,
	}

	if r.Method == "GET" {
		handler.Get()
	}

	if r.Method == "POST" {
		handler.Create()
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
