package hera_test

import (
	hera "ddl/newHERA"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtTable(t *testing.T) {
	table, errNew := hera.NewTable(
		&hera.Person{},
	)
	require.NoError(t, errNew)
	require.NotZero(t, table)

	fmt.Println(table.AsDDLPostgres())
}
