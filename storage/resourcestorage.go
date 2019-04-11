package storage

import (
	"database/sql"
	"fmt"
)

type ResourceStorage struct {
	DB        *sql.DB
	TableName string
}

func (s *ResourceStorage) GetAll() ([]byte, error) {
	return queryToJson(
		s.DB,
		fmt.Sprintf(`SELECT * FROM %s`, s.TableName),
	)
}

func (s *ResourceStorage) GetByID(ID int) ([]byte, error) {
	return queryToJson(
		s.DB,
		fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, s.TableName),
		ID,
	)
}

func (s *ResourceStorage) Create(payload map[string]interface{}) (err error) {
	query := GenerateInsertQuery(s.TableName, payload)
	values := GenerateInsertValues(payload)

	_, err = s.DB.Exec(query, values...)

	return err
}

func (s *ResourceStorage) Update(ID int, payload map[string]interface{}) (err error) {
	query := GenerateUpdateQuery(s.TableName, payload)
	values := GenerateUpdateValues(ID, payload)

	_, err = s.DB.Exec(query, values...)
	return err
}

func (s *ResourceStorage) Delete(ID int) (err error) {
	_, err = s.DB.Exec(
		fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, s.TableName),
		ID,
	)
	return err
}
