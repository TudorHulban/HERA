package pghera

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCo(t *testing.T) {
	h, errCo := New(info)
	assert.Nil(t, errCo)
	h.DBConn.Close()
}
