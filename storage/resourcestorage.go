package storage

import (
	"database/sql"
	"fmt"
)

type ResourceStorage struct {
	DB    *sql.DB
	Table string
}

func (s *ResourceStorage) GetAll() ([]byte, error) {
	return queryToJson(
		s.DB,
		fmt.Sprintf(`SELECT * FROM %s`, s.Table),
	)
}

func (s *ResourceStorage) GetByID(ID int) ([]byte, error) {
	return queryToJson(
		s.DB,
		fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, s.Table),
		ID,
	)
}

func (s *ResourceStorage) Create(payload map[string]interface{}) (err error) {
	query := GenerateInsertQuery(s.Table, payload)
	values := GenerateInsertValues(payload)

	_, err = s.DB.Exec(query, values...)

	return err
}

func (s *ResourceStorage) Update(ID int, payload map[string]interface{}) (err error) {
	query := GenerateUpdateQuery(s.Table, payload)
	values := GenerateUpdateValues(ID, payload)

	_, err = s.DB.Exec(query, values...)
	return err
}

func (s *ResourceStorage) Delete(ID int) (err error) {
	_, err = s.DB.Exec(
		fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, s.Table),
		ID,
	)
	return err
}
