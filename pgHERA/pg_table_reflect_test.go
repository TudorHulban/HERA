package pghera

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetTableName(t *testing.T) {
	h, errCo := New(info, 0, false)
	if assert.Nil(t, errCo) {
		defer h.DBConn.Close()

		assert.Equal(t, h.getTableName(interface{}(&User{})), "users")
	}
}

func TestGetTableColumns(t *testing.T) {
	h, errCo := New(info, 0, false)
	if assert.Nil(t, errCo) {
		defer h.DBConn.Close()

		c, errColumns := h.getTableDefinition(interface{}(&User{}))
		assert.Nil(t, errColumns)
		assert.Equal(t, c.ColumnsDef[1].ColumnName, "name")
		assert.Equal(t, c.ColumnsDef[2].ColumnName, "theage")
	}
}

func TestListTableColumns(t *testing.T) {
	h, errCo := New(info, 1, false)
	if assert.Nil(t, errCo) {
		defer h.DBConn.Close()

		c, errColumns := h.getTableDefinition(interface{}(&User{}))
		assert.Nil(t, errColumns)
		for k, v := range c.ColumnsDef {
			h.l.Print(k, v)
		}
	}
}

func TestProduceModelData(t *testing.T) {
	// define model data to insert
	mdata := User{
		name: "john",
		age:  34,
	}

	// not considering error as we do not need DB
	h, _ := New(info, 3, false)
	_, errPro := h.produceTableColumnShortData(mdata)
	assert.Error(t, errPro)

	data, errData := h.produceTableColumnShortData(&mdata)
	assert.Nil(t, errData)

	h.l.Print("data: ", data)
}
