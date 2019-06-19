package hera

import (
	"database/sql"
	"log"
	"reflect"

	"strings"
	"sync"
)

var onceDB sync.Once

// --------- Decouple
type RDBMS interface {
	NewConnection() (*sql.DB, error)
	NewTable(pDB *sql.DB, pDDL TableDDL) error
// TableExists - returns nil if table exists
	TableExists(pDB *sql.DB, pDatabase, pTableName string) error
	//InsertRow(pDB *sql.DB, pValues *RowData) error
}

// --------- DDL
type ColumnDef struct {
	Name       string
	Type       string
	PrimaryKey bool
	NotNull    bool
}

type TableDDL struct {
	Name        string
	TableFields []ColumnDef
}

type SchemaDDL struct {
	Tables []TableDDL
}

// ----------- INSERT

// RowValues - to be used for insert single row
type RowData struct {
	TableName   string
	ColumnNames string
	Values      []string
}

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

type RowValues struct {
	ColumnNames []string
	Values      []interface{}
}

type TableData struct {
	ColumnNames []string
	Data        [][]interface{}
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

// RowsToSlice - https://kylewbanks.com/blog/query-result-to-map-in-golang
func RowsToSlice(pRows *sql.Rows) (*TableData, error) {
	d := new(TableData)
	d.ColumnNames, _ = pRows.Columns()

	for pRows.Next() {
		columns := make([]interface{}, len(d.ColumnNames))
		columnPointers := make([]interface{}, len(d.ColumnNames))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		err := pRows.Scan(columnPointers...)
		if err != nil {
			log.Println("scan: ", err)
			return nil, err
		}
		d.Data = append(d.Data, columnPointers)
	}
	return d, nil
}
