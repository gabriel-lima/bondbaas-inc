package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"strings"
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	initDb()
    defer db.Close()

	http.HandleFunc("/", handler)
	http.HandleFunc("/schema", schemaHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func initDb() {
    var err error
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"))

    db, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        panic(err)
    }
    err = db.Ping()
    if err != nil {
        panic(err)
    }
    fmt.Println("Successfully connected!")
}

func schemaHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "schema")
}

func handler(w http.ResponseWriter, r *http.Request) {
	uri := strings.Split(r.URL.Path, "/")
	fmt.Fprintf(w, "domain: %s\n", uri[1])
	fmt.Fprintf(w, "id: %s\n", uri[2])
}
