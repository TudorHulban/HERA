package main

import (
	"database/sql"
	"log"
)

func main() {
	dbPath := "lite.db"
	db := DBSQLiteInfo{dbPath}
	dbHandler, err := db.NewConnection()
	if err != nil {
		log.Println("dbHandler: ", err)
	}

	NewTbUsers(dbHandler)

	var i1 TableValues
	i1.TableName = "users"
	i1.ColumnNames = "first_name, last_name, password, role, enabled"
	i1.Values = []string{"john", "smith", "123", "1", "Y"}

	Insert(dbHandler, &i1)
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
