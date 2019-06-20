package hera

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	pq "github.com/lib/pq" // imported for Postgres DB
)

var onceDBPg sync.Once

// DBPostgresInfo - connection details for PostgreSQL type of connections
type DBPostgresInfo struct {
	ip       string
	port     uint
	user     string
	password string
	dbName   string
}

// NewConnection - singleton, only one connection per DB
func (r DBPostgresInfo) NewConnection() (*sql.DB, error) {
	instance := new(sql.DB)
	var err error
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", r.ip, r.port, r.user, r.password, r.dbName)
	onceDBPg.Do(func() {
		instance, err = sql.Open("postgres", dbinfo)
		err = instance.Ping()
	})
	return instance, err
}

// NewTable - Primary Key field is auto incremented
func (r DBPostgresInfo) NewTable(pDB *sql.DB, pDDL TableDDL) error {
	var fieldDDL string
	var columnDDL = func(pDDL ColumnDef) string {
		var notnull, ddl string

		if pDDL.NotNull {
			notnull = " " + "not null"
		}
		if pDDL.PrimaryKey {
			if pDDL.Type == "integer" {
				ddl = pDDL.Name + " " + "SERIAL PRIMARY KEY"
			} else {
				ddl = pDDL.Name + " " + "PRIMARY KEY"
			}
		} else {
			ddl = pDDL.Name + " " + pDDL.Type + notnull
		}
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
func (r DBPostgresInfo) TableExists(pDB *sql.DB, pDatabase, pTableName string) error {
	var occurences bool
	theDML := "SELECT exists (select 1 from information_schema.tables WHERE table_schema='public' AND table_name=" + "'" + pTableName + "'" + ")"
	log.Println("PostgreSQL, checking if table exists: ", theDML)

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
func (r DBPostgresInfo) InsertRow(pDB *sql.DB, pValues *RowData) error {
	theDDL := "insert into " + pValues.TableName + "(" + pValues.ColumnNames + ")" + " values(" + "'" + strings.Join(pValues.Values, "'"+","+"'") + "'" + ")"
	log.Println(theDDL)
	_, err := pDB.Exec(theDDL)
	return err
}

// InsertBulk - insert for multiple rows
func (r DBPostgresInfo) InsertBulk(pDB *sql.DB, pBulk *BulkValues) error {
	//theQuestionMarks := returnNoValues(pBulk.Values[0], "'?'")

	dbTransaction, err := pDB.Begin() // DB Transaction Start
	if err != nil {
		dbTransaction.Rollback()
		return err
	}
	preparedStatem, err := dbTransaction.Prepare(pq.CopyIn(pBulk.TableName, "code", "description", "enabled"))
	if err != nil {
		dbTransaction.Rollback()
		return err
	}
	//statement := "insert into " + pBulk.TableName + "(" + pBulk.ColumnNames + ")" + " values " + theQuestionMarks
	log.Println("insert bulk statement: ", preparedStatem)
	/*
		dml, err := dbTransaction.Prepare(statement)
		if err != nil {
			log.Println("------------------------ Rollback Prepare")
			dbTransaction.Rollback()
			return err
		}
	*/

	for k, columnValues := range pBulk.Values {
		log.Println(k)
		rowValues := sliceToInterface(columnValues)
		stringValues := []string{}
		for _, v := range rowValues {
			stringValues = append(stringValues, "'"+v.(string)+"'")
		}
		log.Println(strings.Join(stringValues, ","))
		result, err := preparedStatem.Exec(strings.Join(stringValues, ","))
		if err != nil {
			log.Println(result)
			log.Println("------------------------ Rollback Transaction")
			dbTransaction.Rollback()
			return err
		}
	}
	log.Println("------------------------ Exit")
	preparedStatem.Exec()
	preparedStatem.Close()
	dbTransaction.Commit() // DB Transaction End
	return nil
}

// Query - returns data as slice of slice of interface{}
func (r DBPostgresInfo) Query(pDB *sql.DB, pSQL string) (*TableData, error) {
	tableData := new(TableData)

	rows, err := pDB.Query(pSQL)
	if err != nil {
		return nil, err
	}
	tableData, err = RowsToSlice(rows)
	if err != nil {
		return nil, err
	}
	return tableData, nil
}
