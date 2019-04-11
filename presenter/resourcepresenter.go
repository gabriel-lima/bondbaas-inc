package presenter

import (
	"log"
	"net/http"
)

type ResourcePresenter struct {
	Response http.ResponseWriter
}

func (p *ResourcePresenter) Success(status int, data []byte) {
	p.Response.Header().Set("Content-Type", "application/json")
	p.Response.WriteHeader(status)
	p.Response.Write(data)
}

func (p *ResourcePresenter) Fail(status int, err string) {
	http.Error(p.Response, err, status)
	log.Println(err)
}
