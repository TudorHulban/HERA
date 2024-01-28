package hera

import (
	"os"
	"testing"
)

func TestSQLite1(t *testing.T) {
	var db DBSQLiteInfo
	db.DBFile = "lite01.dbf"
	testDB(db, "", t)

	err := os.Remove(db.DBFile)
	if err != nil {
		t.Error("Database file not removed")
	}
}

func TestSQLite2(t *testing.T) {
	var db DBSQLiteInfo
	db.DBFile = "lite02.dbf"
	testDB(db, "", t)

	err := os.Remove(db.DBFile)
	if err != nil {
		t.Error("Database file not removed")
	}
}
