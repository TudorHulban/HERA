package main

import (
	"errors"
	"reflect"
	"strings"
)

type Column struct {
	Name   string
	PGType string

	ValueDefault string

	IsPK       bool
	IsNullable bool
	IsUnique   bool
	IsIndexed  bool
}

func NewColumns(object any) (Columns, error) {
	result := make([]*Column, 0)

	var alreadyHavePK bool

	for i := 0; i < reflect.TypeOf(object).Elem().NumField(); i++ {
		fieldRoot := reflect.TypeOf(object).
			Elem().
			FieldByIndex([]int{i})

		column := Column{
			Name: fieldRoot.Name,
		}

		column.PGType = reflectToPG(fieldRoot.Type.String(), column.IsPK)
		column.UpdateWith(fieldRoot.Tag.Get(_TagName), alreadyHavePK)

		if column.IsPK {
			alreadyHavePK = true
		}

		result = append(result, &column)
	}

	return result,
		nil
}

func (col *Column) UpdateWith(tags string, alreadyHavePK bool) error {
	for _, tag := range strings.Split(
		tags, ",",
	) {
		tagClean := strings.ToLower(strings.TrimSpace(tag))

		if len(tagClean) == 0 {
			continue
		}

		if tagClean == "-" {
			continue
		}

		if tagClean == _TagPK {
			if alreadyHavePK {
				return errors.New("more than one primary key field detected. max is 1")
			}

			col.IsPK = true
		}

		if tagClean == _TagUnique {
			col.IsUnique = true
		}

		if tagClean == _TagIndexed {
			col.IsIndexed = true
		}

		if tagClean == _TagRequired {
			col.IsIndexed = true
		}
	}

	return nil
}
