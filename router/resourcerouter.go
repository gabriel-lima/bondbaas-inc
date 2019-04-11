package router

import (
	"bondbaas/presenter"
	"bondbaas/storage"
	"fmt"
)

type ResourceRouter struct {
	AdminStorage      storage.AdminStorage
	ResourcePresenter presenter.GenericPresenter
}

func (r *ResourceRouter) Route(path string) (string, int, bool) {
	tableName, ID := SplitPath(path)

	if r.resourceIsEmpty(tableName) ||
		r.resourceIsAReservedWord(tableName) ||
		r.resourceNotFound(tableName) {
		return "", 0, true
	}

	return tableName, ID, false
}

func (r *ResourceRouter) resourceIsEmpty(tableName string) bool {
	if tableName == "" {
		r.ResourcePresenter.Fail(404, "Please inform a table name")
		return true
	}
	return false
}

func (r *ResourceRouter) resourceIsAReservedWord(tableName string) bool {
	if tableName == "admin" {
		r.ResourcePresenter.Fail(422, "admin is a reserved name.")
		return true
	}
	return false
}

func (r *ResourceRouter) resourceNotFound(tableName string) bool {
	hasTable, err := r.AdminStorage.HasTable(tableName)
	if err != nil {
		r.ResourcePresenter.Fail(500, err.Error())
		return true
	}

	if !hasTable {
		r.ResourcePresenter.Fail(404, fmt.Sprintf("Table %s not found.", tableName))
		return true
	}

	return false
}
