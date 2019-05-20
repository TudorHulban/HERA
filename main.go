package main

import (
	"log"
)

func main() {
	dbPath := "lite.db"
	db := DBSQLiteInfo{dbPath}
	dbHandler, err := db.NewConnection()
	defer dbHandler.Close()
	if err != nil {
		log.Println("dbHandler: ", err)
	}

	db.NewTable(dbHandler, *DDLUsers())
	db.NewTable(dbHandler, *DDLRoles())

}

func DDLUsers() *TableDDL {
	t := TableDDL{Name: "users"}
	t.TableFields = append(t.TableFields, ColumnDef{Name: "id", Type: "integer", PrimaryKey: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "first_name", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "last_name", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "password", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "role", Type: "integer", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "enabled", Type: "text", PrimaryKey: false, NotNull: true})

	return &t
}

func DDLRoles() *TableDDL {
	t := TableDDL{Name: "roles"}
	t.TableFields = append(t.TableFields, ColumnDef{Name: "id", Type: "integer", PrimaryKey: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "code", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "description", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "enabled", Type: "text", PrimaryKey: false, NotNull: true})

	return &t
}
