package pghera

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestInsertModel(t *testing.T) {
	h, errCo := New(LocalDBInfo, 1, true)

	// create table if not exists
	if assert.Nil(t, errCo) {
		defer h.CloseDBConnection() // nolint

		// get table name to check if it exists
		tableName, _, errParse := h.CreateTable(interface{}(&User{}), true)
		assert.Nil(t, errParse)

		// check if table exists
		if h.TableExists(tableName) != nil {
			// does not exist. create it
			_, _, errCr := h.CreateTable(interface{}(&User{}), false)
			assert.Nil(t, errCr)
		}

		// define model data to insert
		mdata := &User{
			name: "john",
			age:  34,
		}

		errIns := h.InsertModel(mdata)
		assert.Nil(t, errIns)
	}
}
