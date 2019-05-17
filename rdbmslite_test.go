package main

import (
	"database/sql"
	"testing"
)

func TestNewTable(t *testing.T) {
	db := DBSQLiteInfo{"lite.db"}
	dbHandler, err := db.NewConnection()
	if err != nil {
		t.Error("dbHandler: ", err)
	}

	NewTbUsers(dbHandler)
	NewTbRoles(dbHandler)

	ex := db.TableExists(dbHandler, "users")
	if !ex {
		t.Error("Table not created")
	}

	ex = db.TableExists(dbHandler, "roles")
	if !ex {
		t.Error("Table not created")
	}

	_, err = dbHandler.Exec("drop table users")
	if err != nil {
		t.Error("Table not dropped")
	}

	ex = db.TableExists(dbHandler, "users")
	if ex {
		t.Error("Table still exists", ex)
	}

	_, err = dbHandler.Exec("drop table roles")
	if err != nil {
		t.Error("Table not dropped")
	}

	ex = db.TableExists(dbHandler, "roles")
	if ex {
		t.Error("Table still exists", ex)
	}
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
