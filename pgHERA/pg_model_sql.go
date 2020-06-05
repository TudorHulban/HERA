package pghera

import (
	"strings"
)

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
	if !h.isItPointer(modelData) {
		return ErrorNotAPointer
	}
	modelData, errData := h.produceTableColumnShortData(modelData)
	if errData != nil {
		return errData
	}
	tbName, errName := h.getTableDefinition(modelData, true)
	if errName != nil {
		return errName
	}

	ddl := []string{"insert into", tbName.TableName, "("}
	for _, v := range modelData.([]ColumnShortData) {
		ddl = append(ddl, v.ColumnName)
	}
	ddl = append(ddl, ") VALUES (")
	for _, v := range modelData.([]ColumnShortData) {
		ddl = append(ddl, v.Value.String())
	}
	ddl = append(ddl, ");")
	ddlSQL := strings.Join(ddl, " ")
	h.l.Print("DDL: ", ddlSQL)
	return nil
}
