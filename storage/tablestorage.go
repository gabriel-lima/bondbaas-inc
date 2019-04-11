package storage

import (
	"database/sql"
	"fmt"
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

func (s *TableGateway) Create(payload map[string]interface{}) (err error) {
	query := GenerateInsertQuery(s.Table, payload)
	values := GenerateInsertValues(payload)

	_, err = s.DB.Exec(query, values...)

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
