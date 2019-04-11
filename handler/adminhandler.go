package handler

import (
	"bondbaas/presenter"
	"bondbaas/service"
	"bondbaas/storage"
	"database/sql"
	"net/http"
)

type AdminHandler struct {
	DB *sql.DB
}

func (h *AdminHandler) Handle(w http.ResponseWriter, r *http.Request) {
	service := h.adminServiceFactory(w)

	switch r.Method {
	case http.MethodGet:
		service.Get()
	case http.MethodPost:
		payload, err := getPayloadSchema(r)
		if err != nil {
			// fail(h.Response, 422, err.Error())
			return
		}
		service.Create(payload)
	default:
		// TODO: define a message error
	}
}

func (h *AdminHandler) adminServiceFactory(w http.ResponseWriter) service.AdminService {
	storage := storage.AdminStorage{DB: h.DB}
	presenter := presenter.GenericPresenter{Response: w}
	return service.AdminService{
		AdminStorage:   storage,
		AdminPresenter: presenter,
	}
}
