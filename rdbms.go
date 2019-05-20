package main

import (
	"database/sql"
	"log"
	"reflect"

	"strings"
	"sync"
)

var onceDB sync.Once

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

func RowsToSlice(pRows *sql.Rows) (*TableData, error) {
	d := new(TableData)

	columns, _ := pRows.Columns()
	cols := make([]interface{}, len(columns))

	for _, v := range columns {
		d.ColumnNames = append(d.ColumnNames, v)
	}

	for pRows.Next() {
		colsPointers := make([]interface{}, len(columns)) // initialize pointers for each row

		for i, _ := range cols {
			colsPointers[i] = &cols[i]
		}

		err := pRows.Scan(colsPointers...)
		if err != nil {
			log.Println("scan: ", err)
			return nil, err
		}

		d.Data = append(d.Data, colsPointers)
	}
	return d, nil
}
