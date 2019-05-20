package main

import (
	"log"
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
	defer dbHandler.Close()

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

	// Testing Row Insert
	var i1 RowData
	i1.TableName = "users"
	i1.ColumnNames = "first_name, last_name, password, role, enabled"
	i1.Values = []string{"john", "smith", "123", "1", "Y"}

	db.InsertRow(dbHandler, &i1)

	// Testing Bulk Insert
	var i2 BulkValues
	i2.TableName = "roles"
	i2.ColumnNames = "code, description, enabled"
	i2.Values = append(i2.Values, []string{"ADMIN", "Full rights", "Y"})
	i2.Values = append(i2.Values, []string{"USER", "Some rights", "Y"})
	i2.Values = append(i2.Values, []string{"GUEST", "Few rights", "Y"})

	db.InsertBulk(dbHandler, &i2)

	// Testing Query - https://kylewbanks.com/blog/query-result-to-map-in-golang
	rows, _ := db.Query(dbHandler, "select * from users where id=1")
	columns, _ := rows.Columns()

	for rows.Next() {
		cols := make([]interface{}, len(columns))
		colsPointers := make([]interface{}, len(columns))
		for i, _ := range cols {
			colsPointers[i] = &cols[i]
		}

		err := rows.Scan(colsPointers...)
		if err != nil {
			log.Println("scan: ", err)
		}

		type cellValue struct {
			columnName string
			cellData   interface{}
		}
		m := []cellValue{}

		for k, colName := range columns {
			m = append(m, cellValue{columnName: colName, cellData: *colsPointers[k].(*interface{})})
		}
		log.Print("m:", m)
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
