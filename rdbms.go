package main

import (
	"reflect"

	"strings"
	"sync"
)

var onceDB sync.Once

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

// RowValues - to be used for insert, select for single row
type RowValues struct {
	TableName   string
	ColumnNames string
	Values      []string
}

type BulkValues struct {
	TableName   string
	ColumnNames string
	Values      [][]string
}

type SchemaDDL struct {
	Tables []TableDDL
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
