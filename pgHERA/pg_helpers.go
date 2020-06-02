package pghera

import (
	"database/sql"
	"strings"
)

// rowsToSlice - https://kylewbanks.com/blog/query-result-to-map-in-golang
// transformation independent of table name.
func rowsToSlice(rows *sql.Rows) (RowData, error) {
	var data RowData
	var errColsClosed error

	data.ColumnNames, errColsClosed = rows.Columns()
	if errColsClosed != nil {
		return RowData{}, errColsClosed
	}

	for rows.Next() {
		tableColumns := make([]interface{}, len(data.ColumnNames))
		columnPointers := make([]interface{}, len(data.ColumnNames))

		for i := range tableColumns {
			columnPointers[i] = &tableColumns[i]
		}

		errScan := rows.Scan(columnPointers...)
		if errScan != nil {
			return RowData{}, errScan
		}
		var rowValues RowValues
		rowValues.Values = columnPointers
		data.Data = append(data.Data, rowValues)
	}
	return data, nil
}

// getColumnDDL Produces column DDL based on column definition.
func getColumnDDL(ddl Column) string {
	result := []string{ddl.ColumnName, ddl.RDBMSType}

	// checking each field in column definition
	if ddl.PK {
		result = append(result, "PRIMARY KEY")
	}
	if ddl.Unique {
		result = append(result, "UNIQUE")
	}
	if ddl.Required {
		result = append(result, "NOT NULL")
	}
	if ddl.DefaultValue != "" {
		var sep string
		if ddl.RDBMSType == "text" {
			sep = `"`
		}
		result = append(result, "DEFAULT")
		result = append(result, sep+ddl.DefaultValue+sep)
	}

	return strings.Join(result, " ")
}
