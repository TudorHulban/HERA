package pghera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTableName(t *testing.T) {
	assert.Equal(t, getTableName(interface{}(&User{})), "users")
}

func TestGetTableColumns(t *testing.T) {
	c := getFields(interface{}(&User{}))

	assert.Equal(t, c[0].ColumnName, "name")
	assert.Equal(t, c[1].ColumnName, "age")
}
