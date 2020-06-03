package pghera

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestTableDDL(t *testing.T) {
	h, errCo := New(info, 0)
	if assert.Nil(t, errCo) {
		defer h.DBConn.Close()

		// get table name to create first
		tableName, errParse := h.CreateTable(interface{}(&User{}), true)
		assert.Nil(t, errParse)

		// check if table exists already
		if h.TableExists(tableName) == nil {
			errDrop := h.DropTable(tableName, true)
			assert.Nil(t, errDrop)
		}

		// table was dropped or did not exist. create it.
		_, errCr := h.CreateTable(interface{}(&User{}), false)
		assert.Nil(t, errCr)
	}
}
