package pghera

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetTableColumns(t *testing.T) {
	h, errCo := New(LocalDBInfo, 0, false)
	if assert.Nil(t, errCo) {
		defer h.CloseDBConnection()

		c, errColumns := h.getTableDefinition(interface{}(&User{}), false)
		assert.Nil(t, errColumns)
		assert.Equal(t, c.ColumnsDef[1].ColumnName, "name")
		assert.Equal(t, c.ColumnsDef[2].ColumnName, "theage")
	}
}

func TestListTableColumns(t *testing.T) {
	h, errCo := New(LocalDBInfo, 1, false)
	if assert.Nil(t, errCo) {
		defer h.CloseDBConnection()

		c, errColumns := h.getTableDefinition(interface{}(&User{}), false)
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
	h, _ := New(LocalDBInfo, 3, false)

	_, errPro := h.produceTableColumnShortData(mdata)
	assert.Error(t, errPro)

	data, errData := h.produceTableColumnShortData(&mdata)
	assert.Nil(t, errData)

	h.l.Print("data: ", data)
}
