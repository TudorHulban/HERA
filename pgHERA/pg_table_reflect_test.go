package pghera

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTableName(t *testing.T) {
	h, errCo := New(info)
	if assert.Nil(t, errCo) {
		defer h.DBConn.Close()

		assert.Equal(t, h.getTableName(interface{}(&User{})), "users")
	}
}

func TestGetTableColumns(t *testing.T) {
	h, errCo := New(info)
	if assert.Nil(t, errCo) {
		defer h.DBConn.Close()

		c, errColumns := h.getTableColumns(interface{}(&User{}))
		assert.Nil(t, errColumns)
		assert.Equal(t, c[1].ColumnName, "name")
		assert.Equal(t, c[2].ColumnName, "age")
	}
}

func TestListTableColumns(t *testing.T) {
	h, errCo := New(info)
	if assert.Nil(t, errCo) {
		defer h.DBConn.Close()

		c, errColumns := h.getTableColumns(interface{}(&User{}))
		assert.Nil(t, errColumns)
		for k, v := range c {
			log.Println(k, v)
		}
	}
}
