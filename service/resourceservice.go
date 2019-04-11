package service

import (
	"bondbaas/storage"
	"bondbaas/presenter"
)

type ResourceService struct {
	ResourceStorage storage.ResourceStorage
	ResourcePresenter presenter.ResourcePresenter
}

func (s *ResourceService) Get(ID int) {
	if ID > 0 {
		s.getByID(ID)
	} else {
		s.getAll()
	}
}

func (s *ResourceService) Create(payload map[string]interface{}) {
	err := s.ResourceStorage.Create(payload)

	if err != nil {
		s.ResourcePresenter.Fail(422, err.Error())
	} else {
		s.ResourcePresenter.Success(201, nil)
	}
}

func (s *ResourceService) Update(ID int, payload map[string]interface{}) {
	err := s.ResourceStorage.Update(ID, payload)

	if err != nil {
		s.ResourcePresenter.Fail(422, err.Error())
	} else {
		s.ResourcePresenter.Success(200, nil)
	}
}

func (s *ResourceService) Delete(ID int) {
	if ID == 0 {
		s.ResourcePresenter.Fail(422, "Id must to be an integer")
		return
	}

	err := s.ResourceStorage.Delete(ID)

	if err != nil {
		s.ResourcePresenter.Fail(422, err.Error())
	} else {
		s.ResourcePresenter.Success(200, nil)
	}
}

func (s *ResourceService) getByID(ID int) {
	data, err := s.ResourceStorage.GetByID(ID)

	if err != nil {
		s.ResourcePresenter.Fail(422, err.Error())
	} else {
		s.ResourcePresenter.Success(200, data)
	}
}

func (s *ResourceService) getAll() {
	data, err := s.ResourceStorage.GetAll()

	if err != nil {
		s.ResourcePresenter.Fail(422, err.Error())
	} else {
		s.ResourcePresenter.Success(200, data)
	}
}
