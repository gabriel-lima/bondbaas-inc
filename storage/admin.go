package storage

import (
	"database/sql"
	"fmt"
	"strings"
)

type AdminGateway struct {
	DB *sql.DB
}

func (g *AdminGateway) GetAll() ([]byte, error) {
	return queryToJson(
		g.DB,
		`SELECT t.table_name, c.column_name, c.data_type, c.is_nullable
		FROM information_schema.tables as t
		JOIN information_schema.columns as c ON c.table_name = t.table_name
		WHERE t.table_schema = 'public'`,
	)
}

func (g *AdminGateway) Create(table Table) (err error) {
	columnsSchema := generateColumnsSchema(table.Columns)

	_, err = g.DB.Exec(
		fmt.Sprintf(`CREATE TABLE %s (id SERIAL PRIMARY KEY %s)`, table.Name, columnsSchema),
	)

	return err
}

func generateColumnsSchema(columns []Column) (columnsSchema string) {
	for _, c := range columns {
		// ID field is always create in a table, so avoiding duplicated field
		if strings.ToLower(c.Name) == "id" {
			continue
		}
		columnsSchema += fmt.Sprintf(", %s %s %s", c.Name, c.Type, c.Constraint)
	}
	return columnsSchema
}
