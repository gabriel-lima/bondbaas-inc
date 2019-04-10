package main

import (
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
	var err error
	var table string
	var id int

	uri := deleteEmptyStrings(strings.Split(r.URL.Path, "/"))

	if len(uri) == 0 {
		http.Error(w, "Root is not a valid route.", 404)
		return
	}

	table = strings.ToLower(uri[0])
	if table == "admin" {
		http.Error(w, "Admin is a reserved name.", 422)
		return
	}
	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM information_schema.tables 
			WHERE table_name = $1
		)`, table).Scan(&exists)
	if err != nil {
		responseInternalError(w, err)
		return
	}
	if !exists {
		http.Error(w, "Table not found.", 404)
		return
	}

	if len(uri) == 2 {
		id, err = strconv.Atoi(uri[1])
		if err != nil {
			responseMalformed(w, err)
			return
		}
	}

	gateway := storage.TableGateway{DB: db, Table: table}

	if r.Method == "GET" {
		if id == 0 {
			data, err := gateway.GetAll()
			if err != nil {
				responseMalformed(w, err)
			} else {
				responseOK(w, data)
			}
		} else {
			data, err := gateway.GetByID(id)
			if err != nil {
				responseMalformed(w, err)
			} else {
				responseOK(w, data)
			}
		}
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

		err = gateway.Create(js)

		if err != nil {
			responseMalformed(w, err)
		} else {
			responseCreated(w)
		}
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

		err = gateway.Update(id, js)

		if err != nil {
			responseMalformed(w, err)
		}
	}

	if r.Method == "DELETE" {
		if id == 0 {
			http.Error(w, "Id must to be an integer", 422)
			return
		}

		err = gateway.Delete(id)

		if err != nil {
			responseMalformed(w, err)
		}
	}
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

	if r.Method == "GET" {
		data, err := gateway.GetAll()

		if err != nil {
			responseInternalError(w, err)
		} else {
			responseOK(w, data)
		}
	}

	if r.Method == "POST" {
		table := storage.Table{}
		err := json.NewDecoder(r.Body).Decode(&table)

		if err != nil {
			responseMalformed(w, err)
			return
		}

		err = gateway.Create(table)

		if err != nil {
			responseMalformed(w, err)
		} else {
			responseCreated(w)
		}
	}
}

func responseOK(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func responseCreated(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}

func responseMalformed(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 422)
	log.Println(err)
}

func responseInternalError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 500)
	log.Println(err)
}

func deleteEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
