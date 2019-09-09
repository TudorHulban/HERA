package hera

import (
	"database/sql"
	"log"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3" // importing sqlite driver
	"github.com/pkg/errors"
)

var onceDBSQLite sync.Once

// DBSQLiteInfo contains RDBMS connection info.
type DBSQLiteInfo struct {
	DBFile string //holds SQLIte DB File
}

// NewConnection - lazy open db, does not truly open until first request
func (r DBSQLiteInfo) NewConnection() (*sql.DB, error) {
	instance := new(sql.DB)
	var err error
	onceDBSQLite.Do(func() {
		instance, err = sql.Open("sqlite3", r.DBFile)
		err = instance.Ping()
	})
	return instance, err
}

// TableExists - returns nil if table exists
func (r DBSQLiteInfo) TableExists(pDB *sql.DB, pDatabase, pTableName string) error {
	var occurences int
	err := pDB.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", pTableName).Scan(&occurences)
	if err != nil {
		return err
	}
	//log.Println(pTable, occurences)
	if occurences != 1 {
		return errors.New("Table does not exist")
	}
	return nil
}

// NewTable is helper for new DB table.
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
		return pDDL.Name + " " + pDDL.Type + pk + notnull
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

// InsertRow is helper for row insert.
func (r DBSQLiteInfo) InsertRow(pDB *sql.DB, pValues *RowData) error {
	theDDL := "insert into " + pValues.TableName + "(" + pValues.ColumnNames + ")" + " values(" + "\"" + strings.Join(pValues.Values, "\""+","+"\"") + "\"" + ")"
	_, err := pDB.Exec(theDDL)
	return err
}

// InsertBulk is helper for insert of multiple rows.
func (r DBSQLiteInfo) InsertBulk(pDB *sql.DB, pBulk *BulkValues) error {
	theQuestionMarks := returnNoValues(pBulk.Values[0], "?", false)

	dbTransaction, errBeginTx := pDB.Begin() // DB Transaction Start
	if errBeginTx != nil {
		errRollBack := dbTransaction.Rollback()
		return errors.Wrap(errBeginTx, errRollBack.Error())
	}

	statement := "insert into " + pBulk.TableName + "(" + pBulk.ColumnNames + ")" + " values " + theQuestionMarks
	dml, errPrepare := dbTransaction.Prepare(statement)
	if errPrepare != nil {
		errRollBack := dbTransaction.Rollback()
		return errors.Wrap(errPrepare, errRollBack.Error())
	}
	defer dml.Close()

	for _, columnValues := range pBulk.Values {
		_, errDML := dml.Exec(sliceToInterface(columnValues)...)

		if errDML != nil {
			errRollBack := dbTransaction.Rollback()
			return errors.Wrap(errDML, errRollBack.Error())
		}
	}
	return dbTransaction.Commit() // DB Transaction End
}

// Query - returns data as slice of slice of interface{}
func (r DBSQLiteInfo) Query(pDB *sql.DB, pSQL string) (*TableData, error) {
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
