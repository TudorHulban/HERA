package pghera

import (
	"strings"
)

// CreateTable Method creates table based on model.
func (h Hera) CreateTable(model interface{}) (string, error) {
	tbDef, errDef := h.getTableDefinition(model)
	if errDef != nil {
		return "", errDef
	}

	colDDL := func(ddl Column) string {
		result := []string{ddl.ColumnName, ddl.RDBMSType}

		if ddl.PK {
			result = append(result, "PRIMARY KEY")
		}
		if ddl.Required {
			result = append(result, "NOT NULL")
		}

		return strings.Join(result, " ")
	}

	tbDDL := []string{"create table", tbDef.TableName, "("}

	for k, v := range tbDef.ColumnsDef {
		tbDDL = append(tbDDL, colDDL(v))
		// adding comma between fields
		if k < len(tbDef.ColumnsDef)-1 {
			tbDDL = append(tbDDL, ",")
		}
	}
	tbDDL = append(tbDDL, ");")
	return strings.Join(tbDDL, " "), nil
}
