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
	_, errCreate := h.DBConn.Exec(tableDDL)
	if errCreate != nil {
		return errCreate
	}

	indexDDL := getIndexDDL(tbDef)
	if indexDDL != "" {
		_, errIndex := h.DBConn.Exec(indexDDL)
		if errIndex != nil {
			return errIndex
		}
	}
	return nil
}
