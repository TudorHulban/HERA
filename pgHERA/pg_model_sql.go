package pghera

// Query Method returns data as slice of slice of interface{}.
func (h Hera) Query(sql string) (TableData, error) {
	rows, errQ := h.DBConn.Query(sql)
	if errQ != nil {
		return TableData{}, errQ
	}

	data, errS := rowsToSlice(rows)
	if errS != nil {
		return TableData{}, errS
	}
	return data, nil
}
