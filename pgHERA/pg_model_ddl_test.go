package pghera

import (
	"log"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestTableDDL(t *testing.T) {
	h, errCo := New(info, 0)
	if assert.Nil(t, errCo) {
		defer h.DBConn.Close()

		log.Println(h.CreateTable(interface{}(&User{})))
	}
}
