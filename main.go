package main

import (
	"database/sql"
	"log"
)

func main() {
	v, err := DBSQLiteInfo{"x3.db"}.NewConnection()
	log.Println(v, err)

	NewTbUsers(v)
	NewTbRoles(v)
}

func NewTbUsers(pDB *sql.DB) error {

	t := TableDDL{Name: "users"}
	t.TableFields = append(t.TableFields, ColumnDef{Name: "id", Type: "integer", PrimaryKey: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "first_name", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "last_name", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "password", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "role", Type: "integer", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "enabled", Type: "text", PrimaryKey: false, NotNull: true})

	return NewTable(pDB, t)
}

func NewTbRoles(pDB *sql.DB) error {

	t := TableDDL{Name: "roles"}
	t.TableFields = append(t.TableFields, ColumnDef{Name: "id", Type: "integer", PrimaryKey: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "code", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "description", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "enabled", Type: "text", PrimaryKey: false, NotNull: true})

	return NewTable(pDB, t)
}

func NewTable(pDB *sql.DB, pDDL TableDDL) error {
	var fieldDDL string

	var columnDDL = func(pDDL ColumnDef) string {
		var notnull, pk string

		if pDDL.NotNull {
			notnull = " " + "not null"
		}

		if pDDL.PrimaryKey {
			pk = " " + "PRIMARY KEY"
		}

		ddl := pDDL.Name + " " + pDDL.Type + pk + notnull
		return ddl
	}

	for k, v := range pDDL.TableFields {
		if k == 0 {
			fieldDDL = columnDDL(v)
		} else {
			fieldDDL = fieldDDL + "," + columnDDL(v)
		}
	}

	ddl := "create table " + pDDL.Name + "(" + fieldDDL + ")"
	log.Println("DDL: ", ddl)

	_, err := pDB.Exec(ddl)
	return err
}
