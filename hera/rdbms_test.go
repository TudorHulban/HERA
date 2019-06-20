package hera

import (
	"database/sql"
	"errors"
	"log"

	"os"
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
	}
	log.Println("dropped in " + pDatabase + " table named " + pTableName)
	return nil
}

func dropTable(pDB *sql.DB, pTableName string) error {
	_, err := pDB.Exec("drop table " + pTableName)
	if err != nil {
		errors.New("Table not dropped")
	}
	return err
}

func TestSQLite1(t *testing.T) {
	var db DBSQLiteInfo
	db.DBFile = "lite01.dbf"
	//testDB(db, "", t)

	err := os.Remove(db.DBFile)
	if err != nil {
		t.Error("Database file not removed")
	}
}

func TestSQLite2(t *testing.T) {
	var db DBSQLiteInfo
	db.DBFile = "lite02.dbf"
	//testDB(db, "", t)

	err := os.Remove(db.DBFile)
	if err != nil {
		t.Error("Database file not removed")
	}
}

func TestMaria(t *testing.T) {
	var db DBMariaInfo
	db.ip = "192.168.1.13"
	db.user = "develop"
	db.password = "develop"
	db.dbName = "devops"
	db.port = 3306
	//testDB(db, db.dbName, t)
}

func TestPostgres(t *testing.T) {
	var db DBPostgresInfo
	db.ip = "192.168.1.15"
	db.user = "develop"
	db.password = "develop"
	db.dbName = "devops"
	db.port = 5432
	testDB(db, db.dbName, t)
}

func testDB(pRDBMS RDBMS, pDatabase string, t *testing.T) {
	dbHandler, err := pRDBMS.NewConnection()
	if err != nil {
		t.Error("could not connect: ", err)
		return
	}
	defer dbHandler.Close()

	cleanTable(dbHandler, pRDBMS, pDatabase, ddlUsers().Name)
	cleanTable(dbHandler, pRDBMS, pDatabase, ddlRoles().Name)

	err = pRDBMS.NewTable(dbHandler, *ddlUsers())
	if err != nil {
		t.Error("createTable: ", err)
		return
	}
	err = pRDBMS.TableExists(dbHandler, pDatabase, ddlUsers().Name)
	if err != nil {
		t.Error("table "+ddlUsers().Name+" not created: ", err)
		return
	}
	err = pRDBMS.NewTable(dbHandler, *ddlRoles())
	if err != nil {
		t.Error("createTable: ", err)
		return
	}
	err = pRDBMS.TableExists(dbHandler, pDatabase, ddlRoles().Name)
	if err != nil {
		t.Error("table "+ddlRoles().Name+" not created: ", err)
		return
	}
	// Testing Row Insert
	var i1 RowData
	i1.TableName = ddlUsers().Name
	for k, v := range ddlUsers().TableFields {
		switch k {
		case 0:
			continue
		case 1:
			i1.ColumnNames = v.Name
		default:
			i1.ColumnNames = i1.ColumnNames + "," + v.Name
		}
	}
	i1.Values = []string{"john", "smith", "123", "1", "Y"}

	err = pRDBMS.InsertRow(dbHandler, &i1)
	if err != nil {
		t.Error("insert row into "+ddlUsers().Name+" did not work: ", err)
		return
	}
	// Testing Bulk Insert
	var i2 BulkValues
	i2.TableName = ddlRoles().Name
	i2.ColumnNames = "code, description, enabled"
	i2.Values = append(i2.Values, []string{"ADMIN", "Full rights", "Y"})
	i2.Values = append(i2.Values, []string{"USER", "Some rights", "Y"})
	i2.Values = append(i2.Values, []string{"GUEST", "Few rights", "Y"})

	err = pRDBMS.InsertBulk(dbHandler, &i2)
	if err != nil {
		t.Error("insert bulk into "+ddlRoles().Name+" did not work: ", err)
		return
	}
	showQryResults(dbHandler, pRDBMS, "select * from users where id=1", t) // Testing Query - Single Insert
	showQryResults(dbHandler, pRDBMS, "select * from roles", t)            // Testing Query - Bulk Insert
}

func showQryResults(pDB *sql.DB, pRDBMS RDBMS, pQuery string, t *testing.T) error {
	rows, err := pRDBMS.Query(pDB, pQuery)
	if err != nil {
		t.Error("Query error")
		return err
	}
	if len(rows.Data) == 0 {
		t.Error("No rows returned by : ", pQuery)
		return nil
	}
	log.Println("Columns: ", rows.ColumnNames)
	rowData := rows.Data[0]
	log.Println("row data for: ", pQuery, "|", rowData)

	for k, v := range rowData {
		rowVal := *v.(*interface{})
		switch rowVal.(type) {
		case []byte:
			{
				log.Println("column ", k, string(rowVal.([]byte)))
			}
		default:
			log.Println("column ", k, rowVal)
		}
	}
	return nil
}

func ddlUsers() *TableDDL {
	t := TableDDL{Name: "users"}
	t.TableFields = append(t.TableFields, ColumnDef{Name: "id", Type: "integer", PrimaryKey: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "first_name", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "last_name", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "pass", Type: "text", PrimaryKey: false, NotNull: true})
	t.TableFields = append(t.TableFields, ColumnDef{Name: "rolebased", Type: "integer", PrimaryKey: false, NotNull: true})
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
