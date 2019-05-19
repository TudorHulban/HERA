package main

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type DBSQLiteInfo struct {
	DBFile string //holds SQLIte DB File
}

func (r DBSQLiteInfo) NewConnection() (*sql.DB, error) {
	instance := new(sql.DB)
	var err error
	onceDB.Do(func() {
		instance, err = sql.Open("sqlite3", r.DBFile)
	})

	return instance, err
}

func (r DBSQLiteInfo) TableExists(pDB *sql.DB, pTable string) bool {
	var occurences int
	_ = pDB.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", pTable).Scan(&occurences)
	//log.Println(pTable, occurences)
	return (occurences == 1)
}

func (r DBSQLiteInfo) NewTable(pDB *sql.DB, pDDL TableDDL) error {
	var fieldDDL string
	var columnDDL = func(pDDL ColumnDef) string {
		var notnull, pk string

		if pDDL.NotNull {
			notnull = " " + "not null"
		}
		if pDDL.PrimaryKey {
			pk = " " + "PRIMARY KEY"
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

func (r DBSQLiteInfo) Insert(pDB *sql.DB, pTable string, pValues *RowData) error {

	theDDL := "insert into " + pTable + "(" + pValues.ColumnNames + ")" + " values(" + "\"" + strings.Join(pValues.Values, "\""+","+"\"") + "\"" + ")"
	_, err := pDB.Exec(theDDL)
	return err
}

func (r DBSQLiteInfo) BulkInsert(pDB *sql.DB, pBulk *BulkValues) error {

	theQuestionMarks := returnNoValues(pBulk.Values[0], "?")

	dbTransaction, err := pDB.Begin() // DB Transaction Start
	if err != nil {
		dbTransaction.Rollback()
		return err
	}

	statement := "insert into " + pBulk.TableName + "(" + pBulk.ColumnNames + ")" + " values " + theQuestionMarks
	dml, err := dbTransaction.Prepare(statement)
	defer dml.Close()

	if err != nil {
		dbTransaction.Rollback()
		return err
	}

	for _, columnValues := range pBulk.Values {
		_, err := dml.Exec(SliceToInterface(columnValues)...)

		if err != nil {
			dbTransaction.Rollback()
			return err
		}
	}
	dbTransaction.Commit() // DB Transaction End
	return nil
}

func (r DBSQLiteInfo) Query(pDB *sql.DB, pSQL string) (*sql.Rows, error) {
	return pDB.Query(pSQL)
}
