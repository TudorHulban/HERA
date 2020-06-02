package pghera

import (
	"strings"
)

// CreateTable Method creates table based on model.
func (h Hera) CreateTable(model interface{}) error {
	tbDef, errDef := h.getTableDefinition(model)
	if errDef != nil {
		return errDef
	}

	tbDDL := []string{"create table", tbDef.TableName, "("}

	for k, v := range tbDef.ColumnsDef {
		tbDDL = append(tbDDL, getColumnDDL(v))

		// adding comma between fields
		if k < len(tbDef.ColumnsDef)-1 {
			tbDDL = append(tbDDL, ",")
		}
	}
	tbDDL = append(tbDDL, ");")

	tableDDL := strings.Join(tbDDL, " ")

	// execute now the DDL
	if _, errCreate := h.DBConn.Exec(tableDDL); errCreate != nil {
		return errCreate
	}

	if indexDDL := getIndexDDL(tbDef); indexDDL != "" {
		if _, errIndex := h.DBConn.Exec(indexDDL); errIndex != nil {
			return errIndex
		}
	}
	return nil
}
