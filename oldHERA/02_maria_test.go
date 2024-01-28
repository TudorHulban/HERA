package hera

import (
	"testing"
)

func TestMaria(t *testing.T) {
	var db DBMariaInfo
	db.ip = "192.168.1.13"
	db.user = "develop"
	db.password = "develop"
	db.dbName = "devops"
	db.port = 3306
	//testDB(db, db.dbName, t)
}
