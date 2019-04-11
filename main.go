package main

import (
	"bondbaas/model"
	"bondbaas/presenter"
	"bondbaas/router"
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

func getPayloadSchema(request *http.Request) (model.TableModel, error) {
	table := model.TableModel{}
	err := json.NewDecoder(request.Body).Decode(&table)
	return table, err
}

func tableHaResource(w http.ResponseWriter, r *http.Request) {
	adminStorage := storage.AdminStorage{DB: db}
	resourcePresenter := presenter.GenericPresenter{
		Response: w,
	}
	resourceRouter := router.ResourceRouter{
		AdminStorage:      adminStorage,
		ResourcePresenter: resourcePresenter,
	}
	tableName, ID, invalidRoute := resourceRouter.Route(r.URL.Path)
	if invalidRoute {
		return
	}

	resourceStorage := storage.ResourceStorage{
		DB:        db,
		TableName: tableName,
	}
	genericPresenter := presenter.GenericPresenter{
		Response: w,
	}
	handler := service.ResourceService{
		ResourceStorage:   resourceStorage,
		ResourcePresenter: genericPresenter,
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
	adminStorage := storage.AdminStorage{
		DB: db,
	}
	adminPresenter := presenter.GenericPresenter{
		Response: w,
	}
	adminService := service.AdminService{
		AdminStorage:   adminStorage,
		AdminPresenter: adminPresenter,
	}

	if r.Method == "GET" {
		adminService.Get()
	}

	if r.Method == "POST" {
		payload, err := getPayloadSchema(r)
		if err != nil {
			// fail(h.Response, 422, err.Error())
			return
		}

		adminService.Create(payload)
	}
}
