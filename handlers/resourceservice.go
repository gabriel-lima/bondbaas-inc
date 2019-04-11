package handlers

import (
	"bondbaas/storage"
	"net/http"
)

type ResourceService struct {
	Response        http.ResponseWriter
	ResourceStorage storage.ResourceStorage
	Request         *http.Request
}

func (h *ResourceService) Get(ID int) {
	if ID > 0 {
		h.getByID(ID)
	} else {
		h.getAll()
	}
}

func (h *ResourceService) Create() {
	payload, err := getPayload(h.Request)
	if err != nil {
		fail(h.Response, 422, err.Error())
		return
	}

	err = h.ResourceStorage.Create(payload)

	if err != nil {
		fail(h.Response, 422, err.Error())
	} else {
		success(h.Response, 201, nil)
	}
}

func (h *ResourceService) Update(ID int) {
	payload, err := getPayload(h.Request)
	if err != nil {
		fail(h.Response, 422, err.Error())
		return
	}

	err = h.ResourceStorage.Update(ID, payload)

	if err != nil {
		fail(h.Response, 422, err.Error())
	} else {
		success(h.Response, 200, nil)
	}
}

func (h *ResourceService) Delete(ID int) {
	if ID == 0 {
		fail(h.Response, 422, "Id must to be an integer")
		return
	}

	err := h.ResourceStorage.Delete(ID)

	if err != nil {
		fail(h.Response, 422, err.Error())
	} else {
		success(h.Response, 200, nil)
	}
}

func (h *ResourceService) getByID(ID int) {
	data, err := h.ResourceStorage.GetByID(ID)

	if err != nil {
		fail(h.Response, 422, err.Error())
	} else {
		success(h.Response, 200, data)
	}
}

func (h *ResourceService) getAll() {
	data, err := h.ResourceStorage.GetAll()

	if err != nil {
		fail(h.Response, 422, err.Error())
	} else {
		success(h.Response, 200, data)
	}
}
