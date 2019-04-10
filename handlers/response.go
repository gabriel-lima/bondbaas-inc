package handlers

import (
	"log"
	"net/http"
)

func success(response http.ResponseWriter, status int, data []byte) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	response.Write(data)
}

func fail(response http.ResponseWriter, status int, err string) {
	http.Error(response, err, status)
	log.Println(err)
}
