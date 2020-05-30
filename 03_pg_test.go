package hera

import (
	"testing"
)

func TestPostgres(t *testing.T) {
	var db DBPostgresInfo
	db.ip = "192.168.1.15"
	db.user = "develop"
	db.password = "develop"
	db.dbName = "devops"
	db.port = 5432
	//testDB(db, db.dbName, t)
}
