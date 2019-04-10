package storage

import (
	"database/sql"
	"fmt"
	"strconv"
)

type TableGateway struct {
	DB    *sql.DB
	Table string
}

func (s *TableGateway) GetAll() ([]byte, error) {
	return queryToJson(
		s.DB,
		fmt.Sprintf(`SELECT * FROM %s`, s.Table),
	)
}

func (s *TableGateway) GetByID(ID int) ([]byte, error) {
	return queryToJson(
		s.DB,
		fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, s.Table),
		ID,
	)
}

func (s *TableGateway) Create(fieldsAndValues map[string]interface{}) (err error) {
	removeIDField(fieldsAndValues)
	fields := extractFieldsToInsertSQL(fieldsAndValues)
	values := extractValuesToInsertSQL(fieldsAndValues)
	placeHolders := generatePlaceHoldersToInsertSQL(fieldsAndValues)

	_, err = s.DB.Exec(
		fmt.Sprintf(`INSERT INTO %s (id%s) VALUES (DEFAULT%s)`, s.Table, fields, placeHolders),
		values...,
	)

	return err
}

func (s *TableGateway) Update(ID int, payload map[string]interface{}) (err error) {
	query := GenerateUpdateQuery(s.Table, payload)
	values := GenerateUpdateValues(ID, payload)

	_, err = s.DB.Exec(query, values...)
	return err
}

func (s *TableGateway) Delete(ID int) (err error) {
	_, err = s.DB.Exec(
		fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, s.Table),
		ID,
	)
	return err
}

func removeIDField(fieldsAndValues map[string]interface{}) {
	delete(fieldsAndValues, "id")
}

func extractFieldsToInsertSQL(fieldsAndValues map[string]interface{}) (fields string) {
	for f, _ := range fieldsAndValues {
		fields += ", " + f
	}
	return fields
}

func extractValuesToInsertSQL(fieldsAndValues map[string]interface{}) []interface{} {
	values := make([]interface{}, 0, len(fieldsAndValues))
	for _, v := range fieldsAndValues {
		values = append(values, v)
	}
	return values
}

func generatePlaceHoldersToInsertSQL(fieldsAndValues map[string]interface{}) (placeHolders string) {
	placeHolderCounter := 1

	for range fieldsAndValues {
		placeHolders += ", $" + strconv.Itoa(placeHolderCounter)
		placeHolderCounter++
	}
	return placeHolders
}
