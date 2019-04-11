package handlers

import (
	"bondbaas/storage"
	"net/http"
)

type AdminHandler struct {
	Response     http.ResponseWriter
	AdminStorage storage.AdminStorage
	Request      *http.Request
}

func (h *AdminHandler) Get() {
	data, err := h.AdminStorage.GetAll()

	if err != nil {
		fail(h.Response, 500, err.Error())
	} else {
		success(h.Response, 200, data)
	}
}

func (h *AdminHandler) Create() {
	payload, err := getPayloadSchema(h.Request)
	if err != nil {
		fail(h.Response, 422, err.Error())
		return
	}

	err = h.AdminStorage.Create(payload)

	if err != nil {
		fail(h.Response, 422, err.Error())
	} else {
		success(h.Response, 201, nil)
	}
}
