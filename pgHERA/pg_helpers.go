package pghera

import (
	"database/sql"
	"strings"
)

/*
File contains helpers that were not glued to model.
*/

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
	if ddl.DefaultValue != "" && ddl.RDBMSType != "text" {
		result = append(result, "DEFAULT")
		result = append(result, ddl.DefaultValue)
	}

	return strings.Join(result, " ")
}

// getIndexDDL Helper creates multi column index DDL based on tag.
func getIndexDDL(tbDef tableDefinition) string {
	var indexCols []string

	for _, v := range tbDef.ColumnsDef {
		if v.Index {
			indexCols = append(indexCols, v.ColumnName)
		}
	}
	// checking if any columns for multi column index
	if len(indexCols) == 0 {
		return ""
	}
	return "create index " + tbDef.TableName + "_idx" + " ON " + tbDef.TableName + " (" + strings.Join(indexCols, ",") + ");"
}
