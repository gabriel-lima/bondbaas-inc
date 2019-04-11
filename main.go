package main

import (
	"bondbaas/handler"
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

	rh := handler.ResourceHandler{DB: db}
	ah := handler.AdminHandler{DB: db}

	http.HandleFunc("/", rh.Handle)
	http.HandleFunc("/admin/tables", ah.Handle)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), nil))
}
