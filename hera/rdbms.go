package hera

import (
	"database/sql"
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"
)

// RDBMS contracts methods for easy working with various RDBMS.
type RDBMS interface {
	NewConnection() (*sql.DB, error)
	NewTable(pDB *sql.DB, pDDL TableDDL) error
	// TableExists - returns nil if table exists
	TableExists(pDB *sql.DB, pDatabase, pTableName string) error
	InsertRow(pDB *sql.DB, pValues *RowData) error
	InsertBulk(pDB *sql.DB, pBulk *BulkValues) error
	Query(pDB *sql.DB, pSQL string) (*TableData, error)
}

// --------- DDL

// ColumnDef provides DDL for table column.
type ColumnDef struct {
	Name       string
	Type       string
	PrimaryKey bool
	NotNull    bool
}

// TableDDL is nutshell for database table DDL.
type TableDDL struct {
	Name        string
	TableFields []ColumnDef
}

// SchemaDDL is nutshell for database schema.
type SchemaDDL struct {
	Tables []TableDDL
}

// ----------- INSERT

// RowData is to be used for insert single row
type RowData struct {
	TableName   string
	ColumnNames string
	Values      []string
}

// BulkValues to be used when working with data from several rows from a specific table.
type BulkValues struct {
	TableName   string
	ColumnNames string
	Values      [][]string
}

// CellValue - for select
type CellValue struct {
	ColumnName string
	CellData   interface{}
}

// RowValues contains data from given row.
type RowValues struct {
	ColumnNames []string
	Values      []interface{}
}

// TableData contains data from given table.
type TableData struct {
	ColumnNames []string
	Data        [][]interface{}
}

// prepareDBFields - switches from single apostrophe used in SQL to double quotes as used in Go
func prepareDBFields(pFields string) string {
	return strings.Replace(pFields, "'", "\"", -1)
}

// cleanTable - drop table if exists
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
		return errors.New("table not dropped: " + err.Error())
	}
	return err
}

// returnNoValues - returns parametrized symbols for inserts
func returnNoValues(pSlice []string, pCharToReturn string, pWithNumber bool) string {
	var toReturn string

	if pWithNumber {
		for k := range pSlice {
			toReturn = toReturn + pCharToReturn + strconv.Itoa(k+1) + ","
		}
	} else {
		for range pSlice {
			toReturn = toReturn + pCharToReturn + ","
		}
	}
	return "(" + toReturn[0:len(toReturn)-1] + ")"
}

// wrapSliceValuesx -
func wrapSliceValuesx(pSlice []string, pCharToWrap string) string {
	return "(" + pCharToWrap + strings.Join(pSlice, pCharToWrap+","+pCharToWrap) + pCharToWrap + ")"
}

// sliceToInterface - transforms unknown into slice of unknown
func sliceToInterface(pSlice interface{}) []interface{} {
	s := reflect.ValueOf(pSlice)

	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}
	result := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		result[i] = s.Index(i).Interface()
	}
	return result
}

// RowsToSlice - https://kylewbanks.com/blog/query-result-to-map-in-golang
func RowsToSlice(pRows *sql.Rows) (*TableData, error) {
	d := new(TableData)
	d.ColumnNames, _ = pRows.Columns()

	for pRows.Next() {
		columns := make([]interface{}, len(d.ColumnNames))
		columnPointers := make([]interface{}, len(d.ColumnNames))

		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		errScan := pRows.Scan(columnPointers...)
		if errScan != nil {
			log.Println("scan: ", errScan)
			return nil, errScan
		}
		d.Data = append(d.Data, columnPointers)
	}
	return d, nil
}
