package pghera

import (
	"errors"
	"strings"
)

/*
File concentrates methods related to Data Definition Language (DDL) operations.
*/

// CreateTable Method creates table based on model.
// In simulation mode it only parses model to return table name that it would create and DDL.
// Returns:
// a. Table name - This way we are sure what was created as table name could have been overidden in structure tags.
// b. Table DDL - Helpfull in tests.
// c. Error - When error the other returns are default value for type, aka "".
func (h Hera) CreateTable(model interface{}, simulateOnly bool) (string, string, error) {
	tbDef, errDef := h.getTableDefinition(model, false)
	if errDef != nil {
		return "", "", errDef
	}

	tbDDL := []string{"create table", tbDef.TableName, "("}

	for k, v := range tbDef.ColumnsDef {
		h.L.Debug("v: ", v)
		tbDDL = append(tbDDL, getColumnDDL(v))

		// adding comma between fields
		if k < len(tbDef.ColumnsDef)-1 {
			tbDDL = append(tbDDL, ",")
		}
	}
	tbDDL = append(tbDDL, ");")
	tableDDL := strings.Join(tbDDL, " ")

	if simulateOnly {
		return tbDef.TableName, tableDDL, nil
	}
	// execute now the DDL
	if _, errCreate := h.DBConn.Exec(tableDDL); errCreate != nil {
		return "", "", errCreate
	}
	// checking if multi column index needs to be created.
	// if any error returning table name for eventual cleaning. table was created at this point.
	if indexDDL := getIndexDDL(tbDef); indexDDL != "" {
		if _, errIndex := h.DBConn.Exec(indexDDL); errIndex != nil {
			return tbDef.TableName, "", errIndex
		}
	}
	// now sure table created. index created.
	return tbDef.TableName, "", nil
}

// TableExists Method returns nil if table exists.
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

// DropTable Method would try to drop passed table.
func (h Hera) DropTable(tableName string, withCascade bool) error {
	ddl := []string{"drop table IF EXISTS", tableName}
	if withCascade {
		ddl = append(ddl, "CASCADE")
	}
	ddl = append(ddl, ";")
	dropDDL := strings.Join(ddl, " ")

	// execute now. w/o context as it is table dropped. might take some time.
	_, errDr := h.DBConn.Exec(dropDDL)
	return errDr
}
