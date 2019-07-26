package hera

import (
	"database/sql"
	"log"
	"testing"
)

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
