package main

import (
	"os"
	"testing"
)

func TestSQLite(t *testing.T) {
	dbPath := "lite.db"
	db := DBSQLiteInfo{dbPath}
	dbHandler, err := db.NewConnection()
	if err != nil {
		t.Error("dbHandler: ", err)
	}

	db.NewTable(dbHandler, *ddlUsers())
	db.NewTable(dbHandler, *ddlRoles())

	ex := db.TableExists(dbHandler, "users")
	if !ex {
		t.Error("Table not created")
	}

	ex = db.TableExists(dbHandler, "roles")
	if !ex {
		t.Error("Table not created")
	}

	// Testing Insert
	var i1 RowValues
	i1.TableName = "users"
	i1.ColumnNames = "first_name, last_name, password, role, enabled"
	i1.Values = []string{"john", "smith", "123", "1", "Y"}

	db.Insert(dbHandler, &i1)

	// Testing Query

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

	err = os.Remove(dbPath)
	if err != nil {
		t.Error("Database file not removed")
	}
}

func ddlUsers() *TableDDL {
	t := TableDDL{Name: "users"}
	t.TableFields = append(t.TableFields, ColumnDef{Name: "id", Type: "integer", PrimaryKey: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "first_name", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "last_name", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "password", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "role", Type: "integer", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "enabled", Type: "text", PrimaryKey: false, NotNull: true})

	return &t
}

func ddlRoles() *TableDDL {
	t := TableDDL{Name: "roles"}
	t.TableFields = append(t.TableFields, ColumnDef{Name: "id", Type: "integer", PrimaryKey: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "code", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "description", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "enabled", Type: "text", PrimaryKey: false, NotNull: true})

	return &t
}
