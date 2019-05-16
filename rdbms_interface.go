package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DBSQLiteInfo struct {
	pathSQLiteFile string
}

type DBPostgresInfo struct {
	ip       string
	port     uint
	user     string
	password string
}

type DBMariaInfo struct {
	ip       string
	port     int
	user     string
	password string
}

var onceDB sync.Once

type RDBMS interface {
	NewConnection(pDBName string) (*sql.DB, error)
	TableExists(pDB *sql.DB, pDatabase, pTable string) bool
	CreateTable(pDB *sql.DB, pDatabase, pTableName, pDDL string, pColumnPKAutoincrement int) bool
	SingleInsert(pDB *sql.DB, pTableName string, pValues []string) error
	BulkInsert(pDB *sql.DB, pTableName string, pColumnNames []string, pValues [][]string) error
}

func (r DBSQLiteInfo) NewConnection(pDBName string) (*sql.DB, error) {

	instance := new(sql.DB)
	var err error
	onceDB.Do(func() {
		instance, err = sql.Open("sqlite3", r.pathSQLiteFile)
	})

	return instance, err
}

func (r DBPostgresInfo) NewConnection(pDBName string) (*sql.DB, error) {

	instance := new(sql.DB)
	var err error

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", r.ip, r.port, r.user, r.password, pDBName)
	onceDB.Do(func() {
		instance, err = sql.Open("postgres", dbinfo)
	})

	if err != nil {
		log.Println("NewConnection: ", err)
		return instance, err
	}
	err = instance.Ping()
	return instance, err
}

func (r DBMariaInfo) NewConnection(pDBName string) (*sql.DB, error) {

	instance := new(sql.DB)
	var err error

	dbinfo := r.user + ":" + r.password + "@tcp(" + r.ip + ":" + strconv.Itoa(r.port) + ")/" + pDBName
	onceDB.Do(func() {
		instance, err = sql.Open("mysql", dbinfo)
	})

	return instance, err
}

func (r DBSQLiteInfo) TableExists(pDB *sql.DB, pDatabase, pTable string) bool {
	var occurences int

	_ = pDB.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", pTable).Scan(&occurences)

	return (occurences == 1)
}

func (rDBPostgresInfo DBPostgresInfo) TableExists(pDB *sql.DB, pDatabase, pTable string) bool {
	var occurences bool

	theDML := "SELECT exists (select 1 from information_schema.tables WHERE table_schema='public' AND table_name=" + "'" + pTable + "'" + ")"
	_ = pDB.QueryRow(theDML).Scan(&occurences)

	return occurences
}

func (rDBMariaInfo DBMariaInfo) TableExists(pDB *sql.DB, pDatabase, pTable string) bool {
	var occurences bool

	theDML := "select count(1) from information_schema.tables WHERE table_schema=" + "'" + pDatabase + "'" + " AND table_name=" + "'" + pTable + "'" + " limit 1"
	_ = pDB.QueryRow(theDML).Scan(&occurences)

	return occurences
}

func (rDBSQLiteInfo DBSQLiteInfo) CreateTable(pDB *sql.DB, pDatabase, pTableName, pDDL string, pColumnPKAutoincrement int) bool {

	var theDDL string

	if pColumnPKAutoincrement > 0 {
		theDDL = "\"id\" INTEGER PRIMARY KEY AUTOINCREMENT," + pDDL
	} else {
		theDDL = pDDL
	}

	theDDL = "CREATE TABLE " + pTableName + "(" + theDDL + ")"

	_, err := pDB.Exec(theDDL)
	checkErr(err, "rDBSQLiteInfo CreateTable")

	return rDBSQLiteInfo.TableExists(pDB, pDatabase, pTableName)
}

func (rDBPostgresInfo DBPostgresInfo) CreateTable(pDB *sql.DB, pDatabase, pTableName, pDDL string, pColumnPKAutoincrement int) bool {

	var theDDL string

	if pColumnPKAutoincrement > 0 {
		theDDL = "\"id\" serial," + pDDL
	} else {
		theDDL = pDDL
	}

	theDDL = "CREATE TABLE " + pTableName + "(" + theDDL + ")"

	_, err := pDB.Exec(theDDL)
	checkErr(err, "rDBPostgresInfo CreateTable "+theDDL)

	return rDBPostgresInfo.TableExists(pDB, pDatabase, pTableName)
}

func (rDBMariaInfo DBMariaInfo) CreateTable(pDB *sql.DB, pDatabase, pTableName, pDDL string, pColumnPKAutoincrement int) bool {

	var theDDL string

	if pColumnPKAutoincrement > 0 {
		theDDL = "\"id\" serial," + pDDL
	} else {
		theDDL = pDDL
	}

	theDDL = "CREATE TABLE " + pTableName + " (" + strings.Replace(theDDL, "\"", "", -1) + ")"

	fmt.Println(theDDL, pColumnPKAutoincrement)

	_, err := pDB.Exec(theDDL)
	checkErr(err, "rDBMariaInfo CreateTable: "+theDDL)

	return rDBMariaInfo.TableExists(pDB, pDatabase, pTableName)
}

func (rDBSQLiteInfo DBSQLiteInfo) SingleInsert(pDB *sql.DB, pTableName string, pValues []string) error {

	theDDL := "insert into " + pTableName + " values(" + "\"" + strings.Join(pValues, "\""+","+"\"") + "\"" + ")"
	_, err := pDB.Exec(theDDL)

	return err
}

func (rDBPostgresInfo DBPostgresInfo) SingleInsert(pDB *sql.DB, pTableName string, pValues []string) error {

	theDDL := "insert into " + pTableName + " values(" + "\"" + strings.Join(pValues, "\""+","+"\"") + "\"" + ")"
	_, err := pDB.Exec(theDDL)

	return err
}

func (rDBMariaInfo DBMariaInfo) SingleInsert(pDB *sql.DB, pTableName string, pValues []string) error {

	theDDL := "insert into " + pTableName + " values(" + "\"" + strings.Join(pValues, "\""+","+"\"") + "\"" + ")"
	_, err := pDB.Exec(theDDL)

	return err
}

func (rDBSQLiteInfo DBSQLiteInfo) BulkInsert(pDB *sql.DB, pTableName string, pColumnNames []string, pValues [][]string) error {

	theQuestionMarks := returnNoValues(pValues[0], "?")

	// -------- DB Transaction Start -----------
	dbTransaction, err := pDB.Begin()
	checkErr(err, "pDB.Begin")

	statement := "insert into " + pTableName + "(" + strings.Join(pColumnNames, ",") + ")" + " values " + theQuestionMarks

	dml, err := dbTransaction.Prepare(statement)
	checkErr(err, "dbTransaction.Prepare")
	defer dml.Close()

	for _, columnValues := range pValues {
		_, err := dml.Exec(SliceToInterface(columnValues)...)

		if err != nil {
			fmt.Println("dml.Exec: ", err)
			panic(err)
		}
	}
	dbTransaction.Commit()
	// -------- DB Transaction End -----------

	return err
}

func (rDBMariaInfo DBMariaInfo) BulkInsert(pDB *sql.DB, pTableName string, pColumnNames []string, pValues [][]string) error {

	theQuestionMarks := returnNoValues(pValues[0], "?")

	// -------- DB Transaction Start -----------
	dbTransaction, err := pDB.Begin()
	checkErr(err, "pDB.Begin")

	statement := "insert into " + pTableName + "(" + strings.Join(pColumnNames, ",") + ")" + " values " + theQuestionMarks

	dml, err := dbTransaction.Prepare(statement)
	checkErr(err, "dbTransaction.Prepare")
	defer dml.Close()

	for _, columnValues := range pValues {
		_, err := dml.Exec(SliceToInterface(columnValues)...)

		if err != nil {
			fmt.Println("dml.Exec: ", err)
			panic(err)
		}
	}
	dbTransaction.Commit()
	// -------- DB Transaction End -----------

	return err
}

func prepareDBFields(pFields string) string {
	return strings.Replace(pFields, "'", "\"", -1)
}

func returnNoValues(pSlice []string, pCharToReturn string) string {
	var toReturn string
	for _ = range pSlice {
		toReturn = toReturn + pCharToReturn + ","
	}

	return "(" + toReturn[0:len(toReturn)-1] + ")"
}

func wrapSliceValuesx(pSlice []string, pCharToWrap string) string {
	return "(" + pCharToWrap + strings.Join(pSlice, pCharToWrap+","+pCharToWrap) + pCharToWrap + ")"
}

func SliceToInterface(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	result := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		result[i] = s.Index(i).Interface()
	}

	return result
}
