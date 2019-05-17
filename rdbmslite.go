package main

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type DBSQLiteInfo struct {
	DBFile string //holds SQLIte DB File
}

func (r DBSQLiteInfo) NewConnection(pDBName string) (*sql.DB, error) {

	instance := new(sql.DB)
	var err error
	onceDB.Do(func() {
		instance, err = sql.Open("sqlite3", r.DBFile)
	})

	return instance, err
}

func (r DBSQLiteInfo) TableExists(pDB *sql.DB, pDatabase, pTable string) bool {
	var occurences int

	_ = pDB.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", pTable).Scan(&occurences)

	return (occurences == 1)
}

func (r DBSQLiteInfo) CreateTable(pDB *sql.DB, pDatabase, pTableName, pDDL string, pColumnPKAutoincrement int) (bool, error) {
	theDDL := pDDL

	if pColumnPKAutoincrement > 0 {
		theDDL = "\"id\" INTEGER PRIMARY KEY AUTOINCREMENT," + pDDL
	}
	theDDL = "CREATE TABLE " + pTableName + "(" + theDDL + ")"
	_, err := pDB.Exec(theDDL)
	if err != nil {
		return false, err
	}
	return r.TableExists(pDB, pDatabase, pTableName), nil
}

func (r DBSQLiteInfo) SingleInsert(pDB *sql.DB, pTableName string, pValues []string) error {
	theDDL := "insert into " + pTableName + " values(" + "\"" + strings.Join(pValues, "\""+","+"\"") + "\"" + ")"
	_, err := pDB.Exec(theDDL)

	return err
}

func (r DBSQLiteInfo) BulkInsert(pDB *sql.DB, pTableName string, pColumnNames []string, pValues [][]string) error {

	theQuestionMarks := returnNoValues(pValues[0], "?")

	dbTransaction, err := pDB.Begin() // DB Transaction Start
	if err != nil {
		dbTransaction.Rollback()
		return err
	}

	statement := "insert into " + pTableName + "(" + strings.Join(pColumnNames, ",") + ")" + " values " + theQuestionMarks
	dml, err := dbTransaction.Prepare(statement)
	defer dml.Close()

	if err != nil {
		dbTransaction.Rollback()
		return err
	}

	for _, columnValues := range pValues {
		_, err := dml.Exec(SliceToInterface(columnValues)...)

		if err != nil {
			dbTransaction.Rollback()
			return err
		}
	}
	dbTransaction.Commit() // DB Transaction End
	return nil
}
