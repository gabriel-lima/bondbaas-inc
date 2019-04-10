package handlers

import (
	"bondbaas/storage"
	"net/http"
)

type AdminHandler struct {
	Response     http.ResponseWriter
	AdminGateway storage.AdminGateway
}

func (h *AdminHandler) Get() {
	data, err := h.AdminGateway.GetAll()

	if err != nil {
		fail(h.Response, 500, err.Error())
	} else {
		success(h.Response, 200, data)
	}
}

func (h *AdminHandler) Create(payload storage.Table) {
	err := h.AdminGateway.Create(payload)

	if err != nil {
		fail(h.Response, 422, err.Error())
	} else {
		success(h.Response, 201, nil)
	}
}
