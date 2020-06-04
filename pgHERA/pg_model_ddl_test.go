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
	h, errCo := New(info, 0)
	if assert.Nil(t, errCo) {
		defer h.DBConn.Close()

		// get table name to create first
		tableName, errParse := h.CreateTable(interface{}(&User{}), true)
		assert.Nil(t, errParse)

		// check if table exists already
		if h.TableExists(tableName) == nil {
			assert.Nil(t, h.DropTable(tableName, true))
		}

		// table was dropped or did not exist. create it.
		_, errCr := h.CreateTable(interface{}(&User{}), false)
		assert.Nil(t, errCr)

		// now we are sure we drop it.
		assert.Nil(t, h.DropTable(tableName, true))
	}
}
