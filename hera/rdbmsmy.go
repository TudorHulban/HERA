package hera

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql" // imported for Maria DB
)

// DBMariaInfo - connection details for Maria DB type of connections
type DBMariaInfo struct {
	ip       string
	port     int
	user     string
	password string
	dbName   string
}

// NewConnection - singleton, only one connection per DB
func (r DBMariaInfo) NewConnection() (*sql.DB, error) {
	instance := new(sql.DB)
	var err error
	onceDB.Do(func() {
		instance, err = sql.Open("mysql", r.user+":"+r.password+"@tcp("+r.ip+":"+strconv.Itoa(r.port)+")/"+r.dbName)
		err = instance.Ping()
	})
	return instance, err
}

// NewTable - Primary Key field is auto incremented
func (r DBMariaInfo) NewTable(pDB *sql.DB, pDDL TableDDL) error {
	var fieldDDL string
	var columnDDL = func(pDDL ColumnDef) string {
		var notnull, pk string

		if pDDL.NotNull {
			notnull = " " + "not null"
		}
		if pDDL.PrimaryKey {
			pk = " " + "AUTO_INCREMENT PRIMARY KEY"
		}
		ddl := pDDL.Name + " " + pDDL.Type + pk + notnull
		return ddl
	}

	for k, v := range pDDL.TableFields {
		if k == 0 {
			fieldDDL = columnDDL(v)
		} else {
			fieldDDL = fieldDDL + "," + columnDDL(v)
		}
	}
	ddl := "create table " + pDDL.Name + "(" + fieldDDL + ")"
	log.Println("DDL: ", ddl)

	_, err := pDB.Exec(ddl)
	return err
}

// TableExists - returns nil if table exists
func (r DBMariaInfo) TableExists(pDB *sql.DB, pDatabase, pTableName string) error {
	var occurences bool
	theDML := "select count(1) from information_schema.tables WHERE table_schema=" + "'" + pDatabase + "'" + " AND table_name=" + "'" + pTableName + "'" + " limit 1"
	log.Println("Maria, checking if table exists: ", theDML)

	err := pDB.QueryRow(theDML).Scan(&occurences)
	if err != nil {
		return err
	}
	if occurences {
		return nil
	}
	return errors.New("Table " + pTableName + " does not exist in " + pDatabase)
}

// InsertRow - single row only
func (r DBMariaInfo) InsertRow(pDB *sql.DB, pValues *RowData) error {
	theDDL := "insert into " + pValues.TableName + "(" + pValues.ColumnNames + ")" + " values(" + "\"" + strings.Join(pValues.Values, "\""+","+"\"") + "\"" + ")"
	_, err := pDB.Exec(theDDL)
	return err
}

// InsertBulk - insert for multiple rows
func (r DBMariaInfo) InsertBulk(pDB *sql.DB, pBulk *BulkValues) error {
	theQuestionMarks := returnNoValues(pBulk.Values[0], "?")

	dbTransaction, err := pDB.Begin() // DB Transaction Start
	if err != nil {
		dbTransaction.Rollback()
		return err
	}

	statement := "insert into " + pBulk.TableName + "(" + pBulk.ColumnNames + ")" + " values " + theQuestionMarks
	log.Println("insert bulk statement: ", statement)
	dml, err := dbTransaction.Prepare(statement)
	if err != nil {
		dbTransaction.Rollback()
		return err
	}
	defer dml.Close()

	for _, columnValues := range pBulk.Values {
		log.Println(SliceToInterface(columnValues)...)
		_, err := dml.Exec(SliceToInterface(columnValues)...)
		if err != nil {
			dbTransaction.Rollback()
			return err
		}
	}
	dbTransaction.Commit() // DB Transaction End
	return nil
}
