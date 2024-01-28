package pghera

import (
	"database/sql"
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
