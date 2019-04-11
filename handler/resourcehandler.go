package handler

import (
	"bondbaas/presenter"
	"bondbaas/router"
	"bondbaas/service"
	"bondbaas/storage"
	"database/sql"
	"net/http"
)

type ResourceHandler struct {
	DB *sql.DB
}

func (h *ResourceHandler) ResourceHandle(w http.ResponseWriter, r *http.Request) {
	resourceRouter := h.resourceRouterFactory(w)

	tableName, ID, invalidRoute := resourceRouter.Route(r.URL.Path)
	if invalidRoute {
		return
	}

	service := h.resourceServiceFactory(w, tableName)

	switch r.Method {
	case http.MethodGet:
		service.Get(ID)
	case http.MethodPost:
		payload, err := getPayload(r)
		if err != nil {
			// s.ResourcePresenter.Fail(422, err.Error())
			return
		}
		service.Create(payload)
	case http.MethodPut:
		payload, err := getPayload(r)
		if err != nil {
			// s.ResourcePresenter.Fail(422, err.Error())
			return
		}
		service.Update(ID, payload)
	case http.MethodDelete:
		service.Delete(ID)
	default:
		// TODO: define a message error
	}
}

func (h *ResourceHandler) resourceRouterFactory(w http.ResponseWriter) router.ResourceRouter {
	adminStorage := storage.AdminStorage{DB: h.DB}
	resourcePresenter := presenter.GenericPresenter{Response: w}
	return router.ResourceRouter{
		AdminStorage:      adminStorage,
		ResourcePresenter: resourcePresenter,
	}
}

func (h *ResourceHandler) resourceServiceFactory(w http.ResponseWriter, tableName string) service.ResourceService {
	storage := storage.ResourceStorage{
		DB:        h.DB,
		TableName: tableName,
	}
	presenter := presenter.GenericPresenter{
		Response: w,
	}

	return service.ResourceService{
		ResourceStorage:   storage,
		ResourcePresenter: presenter,
	}
}
