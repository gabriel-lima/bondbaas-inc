package main

import (
	"bondbaas/handlers"
	"bondbaas/presenter"
	"bondbaas/service"
	"bondbaas/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
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

func getPayload(request *http.Request) (payload map[string]interface{}, err error) {
	var raw []byte
	raw, _ = ioutil.ReadAll(request.Body)

	err = json.Unmarshal([]byte(raw), &payload)

	return payload, err
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

	resourceStorage := storage.ResourceStorage{
		DB:        db,
		TableName: tableName,
	}
	resourcePresenter := presenter.ResourcePresenter{
		Response: w,
	}
	handler := service.ResourceService{
		ResourceStorage:   resourceStorage,
		ResourcePresenter: resourcePresenter,
	}

	if r.Method == "GET" {
		handler.Get(ID)
	}

	if r.Method == "POST" {
		payload, err := getPayload(r)
		if err != nil {
			// s.ResourcePresenter.Fail(422, err.Error())
			return
		}

		handler.Create(payload)
	}

	if r.Method == "PUT" {
		payload, err := getPayload(r)
		if err != nil {
			// s.ResourcePresenter.Fail(422, err.Error())
			return
		}

		handler.Update(ID, payload)
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
