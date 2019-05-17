package main

import (
	"database/sql"
	"log"
)

func main() {
	i := DBSQLiteInfo{"x1.sdb"}

	log.Println("1")
	v, err := i.NewConnection(i.DBFile)
	log.Println(v, err)

	NewTbUsers(v)
}

func NewTbUsers(pDB *sql.DB) error {

	f1 := ColumnDef{Name: "id", Type: "integer", PrimaryKey: true}
	f2 := ColumnDef{Name: "first_name", Type: "text", PrimaryKey: false, NotNull: true}
	f3 := ColumnDef{Name: "last_name", Type: "text", PrimaryKey: false, NotNull: true}
	f4 := ColumnDef{Name: "password", Type: "text", PrimaryKey: false, NotNull: true}
	f5 := ColumnDef{Name: "role", Type: "integer", PrimaryKey: false, NotNull: true}
	f6 := ColumnDef{Name: "enabled", Type: "text", PrimaryKey: false, NotNull: true}

	table := TableDDL{Name: "users", TableFields: []ColumnDef{f1, f2, f3, f4, f5, f6}}
	return NewTable(pDB, table)
}

func NewTable(pDB *sql.DB, pDDL TableDDL) error {

	return nil
}

func columnDDL(pDDL ColumnDef) string {
	var notnull string
	if pDDL.NotNull {
		notnull = "not null"
	}

	ddl := pDDL.Name + " " + pDDL.Type + " " + notnull + " ,"
	return ddl
}
