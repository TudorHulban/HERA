package hera

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type DBMariaInfo struct {
	ip       string
	port     int
	user     string
	password string
	dbName   string
}

func (r DBMariaInfo) NewConnection() (*sql.DB, error) {
	instance := new(sql.DB)
	var err error
	onceDB.Do(func() {
		instance, err = sql.Open("mysql", r.user+":"+r.password+"@tcp("+r.ip+":"+strconv.Itoa(r.port)+")/"+r.dbName)
		err = instance.Ping()
	})
	return instance, err
}

func (r DBMariaInfo) TableExists(pDB *sql.DB, pDatabase, pTable string) bool {
	var occurences bool

	theDML := "select count(1) from information_schema.tables WHERE table_schema=" + "'" + pDatabase + "'" + " AND table_name=" + "'" + pTable + "'" + " limit 1"
	_ = pDB.QueryRow(theDML).Scan(&occurences)

	return occurences
}

func (r DBMariaInfo) NewTable(pDB *sql.DB, pDDL TableDDL) error {
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

func (r DBMariaInfo) CreateTable(pDB *sql.DB, pDatabase, pTableName, pDDL string, pColumnPKAutoincrement int) (bool, error) {
	theDDL := pDDL

	if pColumnPKAutoincrement > 0 {
		theDDL = "\"id\" serial," + pDDL
	}
	theDDL = "CREATE TABLE " + pTableName + " (" + strings.Replace(theDDL, "\"", "", -1) + ")"
	_, err := pDB.Exec(theDDL)
	if err != nil {
		return false, err
	}
	return r.TableExists(pDB, pDatabase, pTableName), nil
}

func (r DBMariaInfo) SingleInsert(pDB *sql.DB, pTableName string, pValues []string) error {
	theDDL := "insert into " + pTableName + " values(" + "\"" + strings.Join(pValues, "\""+","+"\"") + "\"" + ")"
	_, err := pDB.Exec(theDDL)

	return err
}

func (r DBMariaInfo) BulkInsert(pDB *sql.DB, pTableName string, pColumnNames []string, pValues [][]string) error {

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
