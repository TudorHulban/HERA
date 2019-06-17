package hera

import (
	"database/sql"
	"fmt"

	"strings"

	_ "github.com/lib/pq"
)

type DBPostgresInfo struct {
	ip       string
	port     uint
	user     string
	password string
	dbName   string
}

func (r DBPostgresInfo) NewConnection() (*sql.DB, error) {
	instance := new(sql.DB)
	var err error
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", r.ip, r.port, r.user, r.password, r.dbName)
	onceDB.Do(func() {
		instance, err = sql.Open("postgres", dbinfo)
		err = instance.Ping()
	})
	return instance, err
}

func (r DBPostgresInfo) TableExists(pDB *sql.DB, pDatabase, pTable string) bool {
	var occurences bool

	theDML := "SELECT exists (select 1 from information_schema.tables WHERE table_schema='public' AND table_name=" + "'" + pTable + "'" + ")"
	_ = pDB.QueryRow(theDML).Scan(&occurences)

	return occurences
}

func (r DBPostgresInfo) CreateTable(pDB *sql.DB, pDatabase, pTableName, pDDL string, pColumnPKAutoincrement int) (bool, error) {
	theDDL := pDDL

	if pColumnPKAutoincrement > 0 {
		theDDL = "\"id\" serial," + pDDL
	}
	theDDL = "CREATE TABLE " + pTableName + "(" + theDDL + ")"
	_, err := pDB.Exec(theDDL)
	if err != nil {
		return false, err
	}
	return r.TableExists(pDB, pDatabase, pTableName), nil
}

func (r DBPostgresInfo) SingleInsert(pDB *sql.DB, pTableName string, pValues []string) error {
	theDDL := "insert into " + pTableName + " values(" + "\"" + strings.Join(pValues, "\""+","+"\"") + "\"" + ")"
	_, err := pDB.Exec(theDDL)

	return err
}
