package service

import (
	"bondbaas/model"
	"bondbaas/presenter"
	"bondbaas/storage"
)

type AdminService struct {
	AdminStorage   storage.AdminStorage
	AdminPresenter presenter.GenericPresenter
}

func (s *AdminService) Get() {
	data, err := s.AdminStorage.GetAll()

	if err != nil {
		s.AdminPresenter.Fail(500, err.Error())
	} else {
		s.AdminPresenter.Success(200, data)
	}
}

func (s *AdminService) Create(payload model.TableModel) {
	err := s.AdminStorage.Create(payload)

	if err != nil {
		s.AdminPresenter.Fail(422, err.Error())
	} else {
		s.AdminPresenter.Success(201, nil)
	}
}
