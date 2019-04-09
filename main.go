package main

import (
	"net/http"
	"log"
	"strconv"
	"strings"
	"fmt"
	"os"
	"database/sql"
	"reflect"
	"encoding/json"
	"io/ioutil"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	initDB()
	defer db.Close()
	
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/admin/tables", adminTablesHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), nil))
}

func initDB() {
	var err error
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), 
		os.Getenv("DB_PORT"), 
		os.Getenv("DB_USER"), 
		os.Getenv("DB_PASSWORD"), 
		os.Getenv("DB_DATABASE"))
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var sqlStatement string
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
		http.Error(w, err.Error(), 500)
		log.Println(err)
		return
	}
	if !exists {
		http.Error(w, "Table not found.", 404)
		return
	}

	if len(uri) == 2 {
		id, err = strconv.Atoi(uri[1])
		if err != nil {
			http.Error(w, "Id must to be an integer", 422)
			log.Println(err)
			return
		}
	}

	if r.Method == "GET" {
		if id == 0 {
			var js []byte
			sqlStatement = fmt.Sprintf(`
				SELECT * 
				FROM %s`, table)
			js, err = queryToJson(db, sqlStatement)
			if err != nil {
				http.Error(w, err.Error(), 500)
				log.Println(err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		} else {
			var js []byte
			sqlStatement = fmt.Sprintf(`
				SELECT * 
				FROM %s 
				WHERE id = $1`, table)
			js, err = queryToJson(db, sqlStatement, id)
			if err != nil {
				http.Error(w, err.Error(), 500)
				log.Println(err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
	}
	if r.Method == "POST" {
		var js map[string]interface{}
		var body []byte
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 422)
			log.Println(err)
			return
		}
		err = json.Unmarshal([]byte(body), &js)
		if err != nil {
			http.Error(w, err.Error(), 422)
			log.Println(err)
			return
		}
		delete(js, "id")

		var placeHolderNumber int
		var placeHolders string
		valuesFromJs := make([]interface{}, 0, len(js))
		var fieldsFromJs string
		for field, value := range js {
			fieldsFromJs += ", " + field

			valuesFromJs = append(valuesFromJs, value)
			placeHolderNumber++
			placeHolders += ", $" + strconv.Itoa(placeHolderNumber)
		}

		sqlStatement = fmt.Sprintf(`
			INSERT INTO %s 
			(id%s) 
			VALUES (DEFAULT%s)`, table, fieldsFromJs, placeHolders)
		_, err = db.Exec(sqlStatement, valuesFromJs...)
		if err != nil {
			http.Error(w, err.Error(), 422)
			log.Println(err)
			return
		} else {
			w.WriteHeader(http.StatusCreated)
			return
		}
	}
	if r.Method == "PUT" {
		var js map[string]interface{}
		var body []byte
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 422)
			log.Println(err)
			return
		}
		err = json.Unmarshal([]byte(body), &js)
		if err != nil {
			http.Error(w, err.Error(), 422)
			log.Println(err)
			return
		}
		
		params := make([]interface{}, 0, len(js))
		params = append(params, js["id"])
		delete(js, "id")

		placeHolderCounter := 1
		var setValues string

		for field, value := range js {
			params = append(params, value)
			if placeHolderCounter > 1 {
				setValues += ", "
			}
			placeHolderCounter++
			setValues += field + " = $" + strconv.Itoa(placeHolderCounter)
		}

		sqlStatement = fmt.Sprintf(`
			UPDATE %s 
			SET %s
			WHERE id = $1`, table, setValues)
		
		_, err = db.Exec(sqlStatement, params...)
		if err != nil {
			http.Error(w, err.Error(), 422)
			log.Println(err)
			return
		}
	}
	if r.Method == "DELETE" {
		if id == 0 {
			http.Error(w, "Id must to be an integer", 422)
			return
		}

		sqlStatement = fmt.Sprintf(`
			DELETE FROM %s 
			WHERE id = $1`, table)
		_, err = db.Exec(sqlStatement, id)
		if err != nil {
			http.Error(w, err.Error(), 422)
			log.Println(err)
			return
		}
	}
}

/// JSON schema
type Table struct {
	Name string `json:name`
	Columns []Columns `json:columns`
}
type Columns struct {
	Name string `json:name`
	Type string `json:type`
	Constraint string `json:constraint`
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
	var err error
	var sqlStatement string

	if r.Method == "GET" {
		var js []byte
		sqlStatement = `
			SELECT t.table_name, c.column_name, c.data_type, c.is_nullable
			FROM information_schema.tables as t
			JOIN information_schema.columns as c ON c.table_name = t.table_name
			WHERE t.table_schema = 'public'`
		js, err = queryToJson(db, sqlStatement)
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}
	if r.Method == "POST" {
		table := Table{}
		err = json.NewDecoder(r.Body).Decode(&table)
		if err != nil {
			http.Error(w, err.Error(), 422)
			log.Println(err)
			return
		}

		var params string
		for _, column := range table.Columns {
			// ID field is always create in a table, so avoiding duplicated field
			if strings.ToLower(column.Name) == "id" {
				continue
			}
			params += fmt.Sprintf(", %s %s %s", column.Name, column.Type, column.Constraint)
		}

		sqlStatement = fmt.Sprintf(`
			CREATE TABLE %s (
				id SERIAL PRIMARY KEY
				%s
			)`, table.Name, params)
		_, err := db.Exec(sqlStatement)
		if err != nil {
			http.Error(w, err.Error(), 422)
			log.Println(err)
			return
		} else {
			w.WriteHeader(http.StatusCreated)
			return
		}
	}
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

func queryToJson(db *sql.DB, query string, args ...interface{}) ([]byte, error) {
	// an array of JSON objects
	// the map key is the field name
	var objects []map[string]interface{}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		// figure out what columns were returned
		// the column names will be the JSON object field keys
		columns, err := rows.ColumnTypes()
		if err != nil {
			return nil, err
		}

		// Scan needs an array of pointers to the values it is setting
		// This creates the object and sets the values correctly
		values := make([]interface{}, len(columns))
		object := map[string]interface{}{}
		for i, column := range columns {
			object[column.Name()] = reflect.New(column.ScanType()).Interface()
			values[i] = object[column.Name()]
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		objects = append(objects, object)
	}

	// indent because I want to read the output
	return json.MarshalIndent(objects, "", "\t")
}
