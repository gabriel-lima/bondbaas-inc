package storage

import (
	"fmt"
	"sort"
)

func GenerateInsertQuery(tableName string, payload map[string]interface{}) string {
	// Remove ID fields from payload, it is not allowed insert a value into a PK
	delete(payload, "id")

	// Store a string like: ,name, price ...
	var columns string
	// Store a string like: ,$1, $2 ...
	var placeHolders string
	// Initialize at 1 the 1st placeholder
	counter := 1

	keys := orderByColumn(payload)

	for _, key := range keys {
		columns += fmt.Sprintf(", %v", key)
		placeHolders += fmt.Sprintf(", $%v", counter)
		counter++
	}

	base := `INSERT INTO %s (id%s) VALUES (DEFAULT%s)`
	return fmt.Sprintf(base, tableName, columns, placeHolders)
}

func GenerateInsertValues(payload map[string]interface{}) (values []interface{}) {
	keys := orderByColumn(payload)
	for _, k := range keys {
		values = append(values, payload[k])
	}
	return values
}

func GenerateUpdateQuery(tableName string, payload map[string]interface{}) string {
	// Remove ID fields from payload, it is not allowed update a PK
	delete(payload, "id")

	if len(payload) == 0 {
		return ""
	}

	// Store a string like: name = $1, price = $2 ...
	var columns string
	// Initialize at 2 because 1st placeholder is reserver for ID field
	counter := 2

	keys := orderByColumn(payload)

	for _, key := range keys {
		columns += fmt.Sprintf("%v = $%v, ", key, counter)
		counter++
	}

	// Remove last comma
	columns = columns[0 : len(columns)-2]

	base := `UPDATE %s SET %s WHERE id = $1`
	return fmt.Sprintf(base, tableName, columns)
}

func GenerateUpdateValues(ID int, payload map[string]interface{}) (values []interface{}) {
	keys := orderByColumn(payload)
	values = append(values, ID)
	for _, k := range keys {
		values = append(values, payload[k])
	}
	return values
}

func orderByColumn(payload map[string]interface{}) (keys []string) {
	for k := range payload {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
