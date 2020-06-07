package pghera

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// TestTableDDL Tests if table:
// a. created correctly
// b. check if it exists
// c. is dropped corectly
func TestTableDDL(t *testing.T) {
	h, errCo := New(LocalDBInfo, 3, true)
	if assert.Nil(t, errCo) {
		defer h.CloseDBConnection()

		// get table name to create first
		tableName, ddl, errParse := h.CreateTable(interface{}(&User{}), true)
		assert.Nil(t, errParse)
		h.L.Print("Table DDL: ", ddl)

		// check if table exists already
		if h.TableExists(tableName) == nil {
			assert.Nil(t, h.DropTable(tableName, true))
		}

		// table was dropped or did not exist. create it.
		_, _, errCr := h.CreateTable(interface{}(&User{}), false)
		assert.Nil(t, errCr)

		// now we are sure we drop it.
		assert.Nil(t, h.DropTable(tableName, true))
	}
}
