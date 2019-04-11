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
)

var db *sql.DB

func main() {
	db = storage.InitDB()
	defer db.Close()

	http.HandleFunc("/", tableHaResource)
	http.HandleFunc("/admin/tables", adminHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), nil))
}

func tableHaResource(w http.ResponseWriter, r *http.Request) {
	adminStorage := storage.AdminStorage{DB: db}
	middleware := handlers.Middleware{
		Request:      r,
		Response:     w,
		AdminStorage: adminStorage,
	}
	tableName, ID, invalid := middleware.Route()
	if invalid {
		return
	}

	ResourceStorage := storage.ResourceStorage{
		DB:    db,
		Table: tableName,
	}
	handler := handlers.ResourceService{
		Request:         r,
		Response:        w,
		ResourceStorage: ResourceStorage,
	}

	if r.Method == "GET" {
		handler.Get(ID)
	}

	if r.Method == "POST" {
		handler.Create()
	}

	if r.Method == "PUT" {
		handler.Update(ID)
	}

	if r.Method == "DELETE" {
		handler.Delete(ID)
	}
}

/// Create a table schema
/*
POST
/admin/tables
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
	adminStorage := storage.AdminStorage{DB: db}
	handler := handlers.AdminHandler{
		Request:      r,
		Response:     w,
		AdminStorage: adminStorage,
	}

	if r.Method == "GET" {
		handler.Get()
	}

	if r.Method == "POST" {
		handler.Create()
	}
}
