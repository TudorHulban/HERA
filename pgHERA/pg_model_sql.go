package pghera

import "reflect"

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
	reType := reflect.TypeOf(modelData)

	tbName := reflectGetTableName(reType)
	reValue := reflect.ValueOf(reType).Elem()
	h.l.Debug("Type:", reType, reValue, tbName)

	return nil
}
