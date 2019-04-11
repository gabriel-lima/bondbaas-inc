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

func (s *ResourceService) Get(ID int) {
	if ID > 0 {
		s.getByID(ID)
	} else {
		s.getAll()
	}
}

func (s *ResourceService) Create() {
	payload, err := getPayload(s.Request)
	if err != nil {
		fail(s.Response, 422, err.Error())
		return
	}

	err = s.ResourceStorage.Create(payload)

	if err != nil {
		fail(s.Response, 422, err.Error())
	} else {
		success(s.Response, 201, nil)
	}
}

func (s *ResourceService) Update(ID int) {
	payload, err := getPayload(s.Request)
	if err != nil {
		fail(s.Response, 422, err.Error())
		return
	}

	err = s.ResourceStorage.Update(ID, payload)

	if err != nil {
		fail(s.Response, 422, err.Error())
	} else {
		success(s.Response, 200, nil)
	}
}

func (s *ResourceService) Delete(ID int) {
	if ID == 0 {
		fail(s.Response, 422, "Id must to be an integer")
		return
	}

	err := s.ResourceStorage.Delete(ID)

	if err != nil {
		fail(s.Response, 422, err.Error())
	} else {
		success(s.Response, 200, nil)
	}
}

func (s *ResourceService) getByID(ID int) {
	data, err := s.ResourceStorage.GetByID(ID)

	if err != nil {
		fail(s.Response, 422, err.Error())
	} else {
		success(s.Response, 200, data)
	}
}

func (s *ResourceService) getAll() {
	data, err := s.ResourceStorage.GetAll()

	if err != nil {
		fail(s.Response, 422, err.Error())
	} else {
		success(s.Response, 200, data)
	}
}
