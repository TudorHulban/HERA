package main

import (
	"log"

	"../hera"
)

func main() {
	dbPath := "lite.db"
	db := hera.DBSQLiteInfo{dbPath}
	dbHandler, errConn := db.NewConnection()
	if errConn != nil {
		log.Println("dbHandler: ", errConn)
	}
	defer dbHandler.Close()

	db.NewTable(dbHandler, DDLUsers())
	db.NewTable(dbHandler, DDLRoles())
}

func DDLUsers() hera.TableDDL {
	t := hera.TableDDL{Name: "users"}
	t.TableFields = append(t.TableFields, hera.ColumnDef{Name: "id", Type: "integer", PrimaryKey: true})
	t.TableFields = append(t.TableFields, hera.ColumnDef{Name: "first_name", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, hera.ColumnDef{Name: "last_name", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, hera.ColumnDef{Name: "password", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, hera.ColumnDef{Name: "role", Type: "integer", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, hera.ColumnDef{Name: "enabled", Type: "text", PrimaryKey: false, NotNull: true})

	return t
}

func DDLRoles() hera.TableDDL {
	t := hera.TableDDL{Name: "roles"}
	t.TableFields = append(t.TableFields, hera.ColumnDef{Name: "id", Type: "integer", PrimaryKey: true})
	t.TableFields = append(t.TableFields, hera.ColumnDef{Name: "code", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, hera.ColumnDef{Name: "description", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, hera.ColumnDef{Name: "enabled", Type: "text", PrimaryKey: false, NotNull: true})

	return t
}
