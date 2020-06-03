package pghera

import (
	"errors"
	"strings"
)

// CreateTable Method creates table based on model.
// It returns table name as table name could be overidden in struct. This way we are sure what was created.
func (h Hera) CreateTable(model interface{}) (string, error) {
	tbDef, errDef := h.getTableDefinition(model)
	if errDef != nil {
		return "", errDef
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
		return "", errCreate
	}

	// returning table name for eventual cleaning. table was created at this point.
	if indexDDL := getIndexDDL(tbDef); indexDDL != "" {
		if _, errIndex := h.DBConn.Exec(indexDDL); errIndex != nil {
			return tbDef.TableName, errIndex
		}
	}
	return tbDef.TableName, nil
}

// TableExists - returns nil if table exists
func (h Hera) TableExists(tableName string) error {
	theDML := "SELECT exists (select 1 from information_schema.tables WHERE table_schema='public' AND table_name=" + "'" + tableName + "'" + ")"

	var occurrences bool
	if errQ := h.DBConn.QueryRow(theDML).Scan(&occurrences); errQ != nil {
		return errQ
	}
	if occurrences {
		return nil
	}
	return errors.New("table " + tableName + " does not exist in " + h.DBName)
}
