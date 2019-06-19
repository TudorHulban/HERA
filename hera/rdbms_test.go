package hera

import (
	"database/sql"
	"errors"
	"log"

	//"os"
	"testing"
)

func cleanTable(pDB *sql.DB, pRDBMS RDBMS, pDatabase, pTableName string) error {
	err := pRDBMS.TableExists(pDB, pDatabase, pTableName)
	if err != nil {
		log.Println(pDatabase + " does NOT contains " + pTableName)
		return err
	}
	log.Println(pDatabase + " contains " + pTableName)

	err = dropTable(pDB, pTableName)
	if err != nil {
		log.Println("cannot drop table in " + pDatabase + " named " + pTableName)
		return errors.New("drop table: " + err.Error())
	} else {
		log.Println("dropped in " + pDatabase + " table named " + pTableName)
	}
	return nil
}

func dropTable(pDB *sql.DB, pTableName string) error {
	_, err := pDB.Exec("drop table " + pTableName)
	if err != nil {
		errors.New("Table not dropped")
	}
	return err
}

/*
func TestMaria(t *testing.T) {
	var db DBMariaInfo
	db.ip = "192.168.1.13"
	db.user = "develop"
	db.password = "develop"
	db.dbName = "devops"
	db.port = 3306

	// Connect DB
	dbHandler, err := db.NewConnection()
	if err != nil {
		t.Error("could not connect: ", err)
	}
	defer dbHandler.Close()

	cleanTable(dbHandler, db, db.dbName, ddlUsers().Name)
	cleanTable(dbHandler, db, db.dbName, ddlRoles().Name)

	err = db.NewTable(dbHandler, *ddlUsers())
	if err != nil {
		t.Error("createTable: ", err)
	}

	err = db.TableExists(dbHandler, db.dbName, ddlUsers().Name)
	if err != nil {
		t.Error("table "+ddlUsers().Name+" not created: ", err)
	}

	err = db.NewTable(dbHandler, *ddlRoles())
	if err != nil {
		t.Error("createTable: ", err)
	}

	err = db.TableExists(dbHandler, db.dbName, ddlRoles().Name)
	if err != nil {
		t.Error("table "+ddlRoles().Name+" not created: ", err)
	}

	// Testing Row Insert
	var i1 RowData
	i1.TableName = ddlUsers().Name
	i1.ColumnNames = "first_name, last_name, password, role, enabled"
	i1.Values = []string{"john", "smith", "123", "1", "Y"}

	err = db.InsertRow(dbHandler, &i1)
	if err != nil {
		t.Error("insert row into "+ddlUsers().Name+" dit not work: ", err)
	}

	// Testing Bulk Insert
	var i2 BulkValues
	i2.TableName = ddlRoles().Name
	i2.ColumnNames = "code, description, enabled"
	i2.Values = append(i2.Values, []string{"ADMIN", "Full rights", "Y"})
	i2.Values = append(i2.Values, []string{"USER", "Some rights", "Y"})
	i2.Values = append(i2.Values, []string{"GUEST", "Few rights", "Y"})

	err = db.InsertBulk(dbHandler, &i2)
	if err != nil {
		t.Error("insert bulk into "+ddlRoles().Name+" dit not work: ", err)
	}

	// Testing Query
	rows, err := db.Query(dbHandler, "select * from users where id=1")
	if err != nil {
		t.Error("Query error")
	}
	log.Println("Columns: ", rows.ColumnNames)
	rowData := rows.Data[0]

	for k, v := range rowData {
		log.Println("column ", k, *v.(*interface{}))
	}

}
*/

func TestSQLite(t *testing.T) {
	db := DBSQLiteInfo{"lite.dbf"}
	dbHandler, err := db.NewConnection()
	if err != nil {
		t.Error("dbHandler: ", err)
	}
	defer dbHandler.Close()

	// Connect DB
	dbHandler, err = db.NewConnection()
	if err != nil {
		t.Error("could not connect: ", err)
	}
	defer dbHandler.Close()

	cleanTable(dbHandler, db, "", ddlUsers().Name)
	cleanTable(dbHandler, db, "", ddlRoles().Name)

	err = db.NewTable(dbHandler, *ddlUsers())
	if err != nil {
		t.Error("createTable: ", err)
	}

	err = db.TableExists(dbHandler, "", ddlUsers().Name)
	if err != nil {
		t.Error("table "+ddlUsers().Name+" not created: ", err)
	}

	err = db.NewTable(dbHandler, *ddlRoles())
	if err != nil {
		t.Error("createTable: ", err)
	}

	err = db.TableExists(dbHandler, "", ddlRoles().Name)
	if err != nil {
		t.Error("table "+ddlRoles().Name+" not created: ", err)
	}

	// Testing Row Insert
	var i1 RowData
	i1.TableName = ddlUsers().Name
	i1.ColumnNames = "first_name, last_name, password, role, enabled"
	i1.Values = []string{"john", "smith", "123", "1", "Y"}

	err = db.InsertRow(dbHandler, &i1)
	if err != nil {
		t.Error("insert row into "+ddlUsers().Name+" dit not work: ", err)
	}

	// Testing Bulk Insert
	var i2 BulkValues
	i2.TableName = ddlRoles().Name
	i2.ColumnNames = "code, description, enabled"
	i2.Values = append(i2.Values, []string{"ADMIN", "Full rights", "Y"})
	i2.Values = append(i2.Values, []string{"USER", "Some rights", "Y"})
	i2.Values = append(i2.Values, []string{"GUEST", "Few rights", "Y"})

	err = db.InsertBulk(dbHandler, &i2)
	if err != nil {
		t.Error("insert bulk into "+ddlRoles().Name+" dit not work: ", err)
	}

	// Testing Query
	rows, err := db.Query(dbHandler, "select * from users where id=1")
	if err != nil {
		t.Error("Query error")
	}
	log.Println("Columns: ", rows.ColumnNames)
	rowData := rows.Data[0]

	for k, v := range rowData {
		log.Println("column ", k, *v.(*interface{}))
	}

}

/*


	// Testing Query - multiple rows returned
	bulk, err := db.Query(dbHandler, "select * from roles")
	if err != nil {
		t.Error("Query error")
	}
	log.Println(bulk.ColumnNames)
	log.Println("rows returned: ", len(bulk.Data))

	for k1, v1 := range bulk.Data {
		for k2, v2 := range v1 {
			log.Println("row: ", k1, "field: ", bulk.ColumnNames[k2], "value: ", *v2.(*interface{}))
		}
	}

	err = os.Remove(dbPath)
	if err != nil {
		t.Error("Database file not removed")
	}
*/

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
