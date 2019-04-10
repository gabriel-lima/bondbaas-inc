package storage

import (
	"fmt"
)

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
	for key, _ := range payload {
		columns += fmt.Sprintf("%v = $%v, ", key, counter)
		counter++
	}

	// Remove last comma
	columns = columns[0 : len(columns)-2]

	base := `UPDATE %s SET %s WHERE id = $1`
	return fmt.Sprintf(base, tableName, columns)
}

func GenerateUpdateValues(ID int, payload map[string]interface{}) (values []interface{}) {
	values = append(values, ID)
	for _, v := range payload {
		values = append(values, v)
	}
	return values
}
