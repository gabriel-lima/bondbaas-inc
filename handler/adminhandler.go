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
		payload, err := SerializeTableModel(service.AdminPresenter, r.Body)
		if err != nil {
			return
		}

		service.Create(payload)
	default:
		service.AdminPresenter.Fail(422, "Undefined HTTP Method")
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
