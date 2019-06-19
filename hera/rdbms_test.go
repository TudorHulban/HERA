package hera

import (
	"database/sql"
	"errors"
	"log"

	//"log"
	//"os"
	"testing"
)

func createTable(pDB *sql.DB, pRDBMS RDBMS, pTableInfo *TableDDL) error {
	return pRDBMS.NewTable(pDB, *pTableInfo)
}

func hasTable(pDB *sql.DB, pRDBMS RDBMS, pDatabase, pTableName string) error {
	return pRDBMS.TableExists(pDB, pDatabase, pTableName)
}

func createRow(pRDBMS RDBMS, pTableName string, pValues *RowData) error {
	return nil
}

func dropTable(pDB *sql.DB, pTableName string) error {
	_, err := pDB.Exec("drop table " + pTableName)
	if err != nil {
		errors.New("Table not dropped")
	}
	return err
}

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

	err = hasTable(dbHandler, db, db.dbName, ddlUsers().Name)
	if err != nil {
		log.Println(db.dbName + " contains " + ddlUsers().Name)

		err = dropTable(dbHandler, ddlUsers().Name)
		if err != nil {
			log.Println("cannot drop table in " + db.dbName + " named " + ddlUsers().Name)
			t.Error("drop table: ", err)
		}
	} else {
		log.Println(db.dbName + " does NOT contains " + ddlUsers().Name)
	}

	err = createTable(dbHandler, db, ddlUsers())
	if err != nil {
		t.Error("NewTable: ", err)
	}

	/*


		err = createTable(db, ddlUsers())

		err := testTableExists(db, "", ddlUsers().Name)
		if err != nil {
			t.Error("TableExists: ", err01)
		}

		err03 := testDropTable(db, ddlUsers().Name)
		if err03 != nil {
			t.Error("DropTable: ", err01)
		}

	*/
}

/*
func TestSQLite(t *testing.T) {
	var db DBSQLiteInfo
	db.DBFile = "lite.db"

	err01 := testNewTable(db, ddlUsers())
	if err01 != nil {
		t.Error("NewTable: ", err01)
	}

	err02 := testTableExists(db, "", ddlUsers().Name)
	if err02 != nil {
		t.Error("TableExists: ", err01)
	}

	err03 := testDropTable(db, ddlUsers().Name)
	if err03 != nil {
		t.Error("DropTable: ", err01)
	}
}

*/

/*
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

	// Testing Query
	rows, err := db.Query(dbHandler, "select * from users where id=1")
	if err != nil {
		t.Error("Query error")
	}
	log.Println(rows.ColumnNames)
	rowData := rows.Data[0]

	for k, v := range rowData {
		log.Println(k, *v.(*interface{}))
	}

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
