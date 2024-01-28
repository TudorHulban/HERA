package main

import "strings"

type Table struct {
	Name    string
	Columns Columns
}

func NewTable(object any) (*Table, error) {
	columns, errGetColumns := NewColumns(object)
	if errGetColumns != nil {
		return nil,
			errGetColumns
	}

	tableName := strings.ToLower(getObjectName(object))

	if strings.HasPrefix(tableName, _TagPointer) {
		tableName = tableName[1:]
	}

	return &Table{
			Name:    pluralize(tableName),
			Columns: columns,
		},
		nil
}

func (t *Table) AsDDLPostgres() string {
	result := []string{
		"create table if does not exist ",
		t.Name,
		"(",
	}

	for ix, column := range t.Columns {
		result = append(result,
			column.AsDDLPostgres(),
		)

		if ix < len(t.Columns)-1 {
			result = append(result,
				",",
			)
		}
	}

	result = append(result, ");")

	return strings.Join(result, "")
}
