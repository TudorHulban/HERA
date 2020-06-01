package pghera

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCo(t *testing.T) {
	info := DBInfo{
		Port:     5432,
		IP:       "0.0.0.0",
		User:     "postgres",
		Password: "pp",
		DBName:   "postgres",
	}

	_, errCo := New(info)
	assert.Nil(t, errCo)
}
