package hera

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql" // imported for Maria DB
	"github.com/pkg/errors"
)

var onceDBMaria sync.Once

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
	onceDBMaria.Do(func() {
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
	theQuestionMarks := returnNoValues(pBulk.Values[0], "?", false)

	dbTransaction, errBegin := pDB.Begin() // DB Transaction Start
	if errBegin != nil {
		errRollBack := dbTransaction.Rollback()
		return errors.Wrap(errBegin, errRollBack.Error())
	}
	statement := "insert into " + pBulk.TableName + "(" + pBulk.ColumnNames + ")" + " values " + theQuestionMarks
	log.Println("insert bulk statement: ", statement)
	dml, errPrepare := dbTransaction.Prepare(statement)
	if errPrepare != nil {
		errRollBack := dbTransaction.Rollback()
		return errors.Wrap(errPrepare, errRollBack.Error())
	}
	defer dml.Close()

	for _, columnValues := range pBulk.Values {
		log.Println(sliceToInterface(columnValues)...)
		_, errExec := dml.Exec(sliceToInterface(columnValues)...)
		if errExec != nil {
			errRollBack := dbTransaction.Rollback()
			return errors.Wrap(errExec, errRollBack.Error())
		}
	}
	return dbTransaction.Commit() // DB Transaction End
}

// Query - returns data as slice of slice of interface{}
func (r DBMariaInfo) Query(pDB *sql.DB, pSQL string) (*TableData, error) {
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
