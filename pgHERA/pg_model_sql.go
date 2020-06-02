package pghera

// Query Method returns data as slice of slice of interface{}.
func (h Hera) Query(sql string) (RowData, error) {
	rows, errQ := h.DBConn.Query(sql)
	if errQ != nil {
		return RowData{}, errQ
	}

	data, errS := rowsToSlice(rows)
	if errS != nil {
		return RowData{}, errS
	}
	return data, nil
}

// CreateTable Method returns data as slice of slice of interface{}.
func (h Hera) CreateTable() error {
	return nil
}
