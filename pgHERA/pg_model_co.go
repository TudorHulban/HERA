package pghera

import (
	"database/sql"
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
	// used for translating structure fields to RDBMS field types
	transTable *translationTable
	L          *log.LogInfo
}

// New Constructor for database connection. Preferable only one connection per DB.
func New(db DBInfo, logLevel int, wCaller bool) (Hera, error) {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db.IP, db.Port, db.User, db.Password, db.DBName)

	dbconn, errOpen := sql.Open("postgres", dbinfo)
	if errOpen != nil {
		return Hera{}, errOpen
	}
	if errAlive := dbconn.Ping(); errAlive != nil {
		return Hera{}, errAlive
	}
	return Hera{
		DBInfo:     db,
		DBConn:     dbconn,
		transTable: newTranslationTable(),
		L:          log.New(logLevel, os.Stderr, wCaller),
	}, nil
}

// CloseConnection Method closes connection to database.
func (h Hera) CloseDBConnection() error {
	return h.DBConn.Close()
}
