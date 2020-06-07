package pghera

import (
	"fmt"
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

// InsertModel Method tries to insert data passed as model data.
// Fields in model that are not passed are inserted with Go default values.
func (h Hera) InsertModel(modelData interface{}) error {
	if !h.isItPointer(modelData) {
		return ErrorNotAPointer
	}
	modelColumnData, errData := h.produceTableColumnShortData(modelData)
	if errData != nil {
		return errData
	}
	tbName, errName := h.getTableDefinition(modelData, true)
	if errName != nil {
		return errName
	}
	ddl := []string{"insert into", tbName.TableName, "("}
	for k, v := range modelColumnData {
		ddl = append(ddl, strings.ToLower(v.ColumnName))
		if k < len(modelColumnData)-1 {
			ddl = append(ddl, ",")
		}
	}
	ddl = append(ddl, ") VALUES (")
	for k, v := range modelColumnData {
		var delim string
		if v.RDBMSType == "string" {
			delim = `'`
		}
		ddl = append(ddl, delim+fmt.Sprintf("%v", v.Value)+delim)
		if k < len(modelColumnData)-1 {
			ddl = append(ddl, ",")
		}
	}
	ddl = append(ddl, ");")
	ddlSQL := strings.Join(ddl, " ")
	h.L.Print("SQL: ", ddlSQL)

	_, errIns := h.DBConn.Exec(ddlSQL)
	return errIns
}
