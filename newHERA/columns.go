package main

import (
	"fmt"
	"strings"
)

type Columns []*Column

func (cols Columns) String() string {
	result := []string{
		fmt.Sprintf(
			"Columns: %d",
			len(cols),
		),
	}

	for ix, column := range cols {
		result = append(result,
			fmt.Sprintf(
				"%d. Name: %s, Type: %s, IsPK: %t, IsNullable: %t, IsIndexed: %t",
				ix+1,
				column.Name,
				column.PGType,
				column.IsPK,
				column.IsNullable,
				column.IsIndexed,
			),
		)
	}

	return strings.Join(result, "\n")
}
