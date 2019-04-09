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

func (s *TableGateway) Update(ID int, fieldsAndValues map[string]interface{}) (err error) {
	removeIDField(fieldsAndValues)
	values := extractValuesToUpdateSQL(ID, fieldsAndValues)
	placeHolders := generatePlaceHoldersToUpdateSQL(fieldsAndValues)

	_, err = s.DB.Exec(
		fmt.Sprintf(`UPDATE %s SET %s WHERE id = $1`, s.Table, placeHolders),
		values...,
	)
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

func extractValuesToUpdateSQL(ID int, fieldsAndValues map[string]interface{}) []interface{} {
	values := make([]interface{}, 0, len(fieldsAndValues))
	values = append(values, ID)
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

func generatePlaceHoldersToUpdateSQL(fieldsAndValues map[string]interface{}) (placeHolders string) {
	placeHolderCounter := 2

	for field, _ := range fieldsAndValues {
		if placeHolderCounter > 1 {
			placeHolders += ", "
		}

		placeHolders += field + " = $" + strconv.Itoa(placeHolderCounter)
		placeHolderCounter++
	}
	return placeHolders
}
