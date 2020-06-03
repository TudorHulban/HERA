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

		tableName, errCr := h.CreateTable(interface{}(&User{}))
		if assert.Nil(t, errCr) {
			errTbExists := h.TableExists(tableName)
			assert.Nil(t, errTbExists)
		}
	}
}
