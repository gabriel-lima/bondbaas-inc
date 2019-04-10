package handlers

import (
	"bondbaas/storage"
	"net/http"
)

type TableHandler struct {
	Response     http.ResponseWriter
	TableGateway storage.TableGateway
	Request      *http.Request
}

func (h *TableHandler) Get(ID int) {
	if ID > 0 {
		h.getByID(ID)
	} else {
		h.getAll()
	}
}

func (h *TableHandler) Create() {
	payload, err := getPayload(h.Request)
	if err != nil {
		fail(h.Response, 422, err.Error())
		return
	}

	err = h.TableGateway.Create(payload)

	if err != nil {
		fail(h.Response, 422, err.Error())
	} else {
		success(h.Response, 201, nil)
	}
}

func (h *TableHandler) Update(ID int) {
	payload, err := getPayload(h.Request)
	if err != nil {
		fail(h.Response, 422, err.Error())
		return
	}

	err = h.TableGateway.Update(ID, payload)

	if err != nil {
		fail(h.Response, 422, err.Error())
	} else {
		success(h.Response, 200, nil)
	}
}

func (h *TableHandler) Delete(ID int) {
	if ID == 0 {
		fail(h.Response, 422, "Id must to be an integer")
		return
	}

	err := h.TableGateway.Delete(ID)

	if err != nil {
		fail(h.Response, 422, err.Error())
	} else {
		success(h.Response, 200, nil)
	}
}

func (h *TableHandler) getByID(ID int) {
	data, err := h.TableGateway.GetByID(ID)

	if err != nil {
		fail(h.Response, 422, err.Error())
	} else {
		success(h.Response, 200, data)
	}
}

func (h *TableHandler) getAll() {
	data, err := h.TableGateway.GetAll()

	if err != nil {
		fail(h.Response, 422, err.Error())
	} else {
		success(h.Response, 200, data)
	}
}
