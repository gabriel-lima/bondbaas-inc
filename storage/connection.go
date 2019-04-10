package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func InitDB() (DB *sql.DB) {
	var err error
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"))

	DB, err = sql.Open("postgres", dataSourceName)

	if err != nil {
		log.Panic(err)
	}

	if err = DB.Ping(); err != nil {
		log.Panic(err)
	}

	return DB
}

