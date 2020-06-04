package pghera

import "fmt"

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
func (h Hera) InsertModel(modelData interface{}) error {
	// get model type
	h.l.Debug("Type:", fmt.Sprintf("%T", modelData))

	return nil
}
