package router

import (
	"strconv"
	"strings"
)

func SplitPath(path string) (tableName string, ID int) {
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
