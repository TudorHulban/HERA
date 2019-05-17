package main

import (
	"database/sql"
	"reflect"

	"strings"
	"sync"
)

var onceDB sync.Once

type RDBMS interface {
	NewConnection(pDBName string) (*sql.DB, error)
	TableExists(pDB *sql.DB, pDatabase, pTable string) bool
	CreateTable(pDB *sql.DB, pDatabase, pTableName, pDDL string, pColumnPKAutoincrement int) bool
	SingleInsert(pDB *sql.DB, pTableName string, pValues []string) error
	BulkInsert(pDB *sql.DB, pTableName string, pColumnNames []string, pValues [][]string) error
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
