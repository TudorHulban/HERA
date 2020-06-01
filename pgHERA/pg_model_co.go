package pghera

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/TudorHulban/log"
)

// DBInfo Type concentrates information for connecting to a PostgreSQL db.
type DBInfo struct {
	Port     uint16
	IP       string
	User     string
	Password string
	DBName   string
}

// Hera Type concentrates RDBMS operations.
type Hera struct {
	DBInfo
	DBConn *sql.DB
	l      *log.LogInfo
}

// New Constructor for database connection. Preferable only one connection per DB.
func New(db DBInfo) (Hera, error) {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db.IP, db.Port, db.User, db.Password, db.DBName)

	dbconn, errOpen := sql.Open("postgres", dbinfo)
	if errOpen != nil {
		return Hera{}, errOpen
	}
	errAlive := dbconn.Ping()
	if errAlive != nil {
		return Hera{}, errAlive
	}
	return Hera{
		DBInfo: db,
		DBConn: dbconn,
		l:      log.New(3, os.Stderr),
	}, nil
}

// TableExists - returns nil if table exists
func (h Hera) TableExists(tableName string) error {
	theDML := "SELECT exists (select 1 from information_schema.tables WHERE table_schema='public' AND table_name=" + "'" + tableName + "'" + ")"

	var occurences bool
	if errQ := h.DBConn.QueryRow(theDML).Scan(&occurences); errQ != nil {
		return errQ
	}
	if occurences {
		return nil
	}
	return errors.New("table " + tableName + " does not exist in " + h.DBName)
}
