package pghera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertModel(t *testing.T) {
	h, errCo := New(info, 3)

	// create table if not exists
	if assert.Nil(t, errCo) {
		defer h.DBConn.Close()

		// get table name to check if it exists
		tableName, errParse := h.CreateTable(interface{}(&User{}), true)
		assert.Nil(t, errParse)

		// check if table exists
		if h.TableExists(tableName) != nil {
			// does not exist. create it
			_, errCr := h.CreateTable(interface{}(&User{}), false)
			assert.Nil(t, errCr)
		}

		// define model data to insert
		mdata := User{
			name: "john",
			age:  34,
		}
		h.l.Debug(mdata)

		errIns := h.InsertModel(mdata)
		assert.Nil(t, errIns)

	}
}
