package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTable(t *testing.T) {
	columns, errParse := NewColumns(
		&Person{},
	)
	require.NoError(t, errParse)
	require.NotZero(t, columns)

	fmt.Println(columns)
}
