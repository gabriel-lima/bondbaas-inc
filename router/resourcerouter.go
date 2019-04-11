package router

import (
	"bondbaas/presenter"
	"bondbaas/storage"
	"fmt"
	"strconv"
	"strings"
)

type ResourceRouter struct {
	AdminStorage      storage.AdminStorage
	ResourcePresenter presenter.GenericPresenter
}

func (r *ResourceRouter) Route(path string) (string, int, bool) {
	tableName, ID := split(path)

	if r.tableNameIsEmpty(tableName) ||
		r.tableNameIsAReservedWord(tableName) ||
		r.tableNameNotFound(tableName) {
		return "", 0, true
	}

	return tableName, ID, false
}

func split(path string) (tableName string, ID int) {
	paths := strings.Split(path, "/")

	paths = removeEmptySpace(paths)

	tableName = strings.ToLower(strings.Join(paths[0:1], ""))

	ID, _ = strconv.Atoi(strings.Join(paths[1:], ""))

	return tableName, ID
}

func removeEmptySpace(paths []string) (newPaths []string) {
	for _, path := range paths {
		if path != "" {
			newPaths = append(newPaths, path)
		}
	}
	return newPaths
}

func (r *ResourceRouter) tableNameIsEmpty(tableName string) bool {
	if tableName == "" {
		r.ResourcePresenter.Fail(404, "Please inform a table name")
		return true
	}
	return false
}

func (r *ResourceRouter) tableNameIsAReservedWord(tableName string) bool {
	if tableName == "admin" {
		r.ResourcePresenter.Fail(422, "admin is a reserved name.")
		return true
	}
	return false
}

func (r *ResourceRouter) tableNameNotFound(tableName string) bool {
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
