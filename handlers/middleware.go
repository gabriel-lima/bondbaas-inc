package handlers

import (
	"bondbaas/storage"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Middleware struct {
	Request      *http.Request
	Response     http.ResponseWriter
	AdminGateway storage.AdminGateway
}

func (m *Middleware) Route() (string, int, bool) {
	tableName, ID := split(m.Request.URL.Path)

	if m.tableNameIsEmpty(tableName) ||
		m.tableNameIsAReservedWord(tableName) ||
		m.tableNameNotFound(tableName) {
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

func (m *Middleware) tableNameIsEmpty(tableName string) bool {
	if tableName == "" {
		fail(m.Response, 404, "Please inform a table name")
		return true
	}
	return false
}

func (m *Middleware) tableNameIsAReservedWord(tableName string) bool {
	if tableName == "admin" {
		fail(m.Response, 422, "admin is a reserved name.")
		return true
	}
	return false
}

func (m *Middleware) tableNameNotFound(tableName string) bool {
	hasTable, err := m.AdminGateway.HasTable(tableName)
	if err != nil {
		fail(m.Response, 500, err.Error())
		return true
	}

	if !hasTable {
		fail(m.Response, 404, fmt.Sprintf("Table %s not found.", tableName))
		return true
	}

	return false
}
