package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	v, err := NewSQLiteDB("x1.sdb")
	log.Println(v, err)
}

func NewSQLiteDB(pDBFile string) (string, error) {

	_, err := os.Stat(pDBFile)

	if err == nil {
		return "", errors.New("File already exists!")
	}

	v, err := GetSQLiteVersion(pDBFile)
	return v, err
}

func GetSQLiteVersion(pDBFile string) (string, error) {
	db, err := sql.Open("sqlite3", pDBFile)
	defer db.Close()

	if err != nil {
		log.Println("sqlite version: ", err)
		return "", err
	}

	var version string

	err = db.QueryRow("select sqlite_version()").Scan(&version)
	return version, err
}
