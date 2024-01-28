package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTable(t *testing.T) {
	table, errNew := NewTable(
		&Person{},
	)
	require.NoError(t, errNew)
	require.NotZero(t, table)

	fmt.Println(table.AsDDLPostgres())
}
