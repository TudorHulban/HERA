package pghera

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTableName(t *testing.T) {
	assert.Equal(t, getTableName(interface{}(&User{})), "users")
}

func TestGetTableColumns(t *testing.T) {
	c := getTableColumns(interface{}(&User{}))

	assert.Equal(t, c[1].ColumnName, "name")
	assert.Equal(t, c[2].ColumnName, "age")
}

func TestListTableColumns(t *testing.T) {
	c := getTableColumns(interface{}(&User{}))

	for k, v := range c {
		log.Println(k, v)
	}
}
