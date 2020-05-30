package pghera

import (
	"database/sql"
)

// rowsToSlice - https://kylewbanks.com/blog/query-result-to-map-in-golang
func rowsToSlice(rows *sql.Rows) (TableData, error) {
	var data TableData
	var errColsClosed error

	data.ColumnNames, errColsClosed = rows.Columns()
	if errColsClosed != nil {
		return TableData{}, errColsClosed
	}

	for rows.Next() {
		tableColumns := make([]interface{}, len(data.ColumnNames))
		columnPointers := make([]interface{}, len(data.ColumnNames))

		for i := range tableColumns {
			columnPointers[i] = &tableColumns[i]
		}

		errScan := rows.Scan(columnPointers...)
		if errScan != nil {
			return TableData{}, errScan
		}

		data.Rows = append(data.Rows, columnPointers)
	}
	return data, nil
}
