package pghera

/*
File concentrates model helpers specific to reflect.
*/

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jinzhu/inflection"
)

// ColumnShortData concentrates model field data.
type ColumnShortData struct {
	ColumnName string
	RDBMSType  string
	Value      reflect.Value
}

// Column Concentrates table column data.
type Column struct {
	// If primary key.
	PK bool
	// If not required column could be null.
	Required     bool
	Unique       bool
	Index        bool
	ColumnName   string
	RDBMSType    string
	DefaultValue string
}

type tableDefinition struct {
	TableName  string
	ColumnsDef []Column
}

// isItPointer Checks if param is a pointer.
func (h Hera) isItPointer(model interface{}) bool {
	h.L.Debug("Passed data of type: ", reflect.TypeOf(model))
	return strings.HasPrefix(reflect.TypeOf(model).String(), "*")
}

// produceTableColumnShortData Helper method to create table columns short info. helper takes a pointer type.
// should parse tag for column name override.
func (h Hera) produceTableColumnShortData(model interface{}) ([]ColumnShortData, error) {
	if !h.isItPointer(model) {
		return []ColumnShortData{}, ErrorNotAPointer
	}

	result := []ColumnShortData{}
	for i := 0; i < reflect.TypeOf(model).Elem().NumField(); i++ {
		// append only types of fields in table translation and that are not ignored (tag "-") or primary key ( auto incremented ).
		// fields that are not passed would be considered with default values.
		fieldRoot := reflect.TypeOf(model).Elem().FieldByIndex([]int{i})
		h.L.Debug("field type: ", fieldRoot.Type.String(), " - ", fieldRoot.Tag)

		if !strings.Contains(fmt.Sprintf("%v", fieldRoot.Tag), `"-"`) && !strings.Contains(fmt.Sprintf("%v", fieldRoot.Tag), `"pk"`) {
			// data for column. creating a pointer to fill it up.
			columnDef := new(Column)
			_, errPar := h.parseFieldTags(fmt.Sprintf("%v", fieldRoot.Tag), columnDef, false)
			if errPar != nil {
				return []ColumnShortData{}, errPar
			}

			// check if field type is in translation table
			rdbmsType := getRDBMSType(fieldRoot.Type.String(), columnDef.PK)
			if rdbmsType == "" {
				// type not supported. if it has tags it is an error. if no tags continue.
				if strings.Contains(fmt.Sprintf("%v", fieldRoot.Tag), `"hera"`) {
					return []ColumnShortData{}, errors.New("type " + fieldRoot.Type.String() + " cannot be translated to a RDBMS type")
				}
				continue
			}
			// check if any column name overridden in tag.
			var columnName string
			if columnDef.ColumnName != "" {
				columnName = columnDef.ColumnName
			} else {
				columnName = fieldRoot.Name
			}

			result = append(result, ColumnShortData{
				ColumnName: columnName,
				// to be taken from translation table
				RDBMSType: rdbmsType,
				Value:     reflect.ValueOf(model).Elem().FieldByIndex([]int{i}),
			})
		}
	}
	return result, nil
}

// reflectGetTableName Helper in case the table name value is needed using reflection types.
func reflectGetTableName(v reflect.Type) string {
	if v.Kind() == reflect.Ptr {
		return inflection.Plural(strings.ToLower(v.Elem().Name()))
	}
	return inflection.Plural(strings.ToLower(v.Name()))
}

// parseFieldTags Method takes "hera" field tags and a pointer to already populated column definition.
// It populates even more the column definition. It returns
// a. boolean to ignore the field with passed structure. Ignore if true.
// b. error
func (h Hera) parseFieldTags(fieldTags string, columnDef *Column, existsPK bool) (bool, error) {
	for _, tag := range strings.Split(fieldTags, ",") {
		s := strings.ToLower(strings.TrimSpace(tag))

		// tag to ignore field.
		if s == "-" {
			return true, nil
		}
		// tag for primary key.
		if s == "pk" {
			if existsPK {
				return false, errors.New("more than one primary key field detected. max is 1")
			}
			existsPK = true
			columnDef.PK = true
		}
		// tag for unicity on column.
		if s == "unique" {
			columnDef.Unique = true
		}
		// tag for being part in multi column index. for single column index use unique
		if s == "index" {
			columnDef.Index = true
		}
		// field cannot be null
		if s == "required" {
			columnDef.Required = true
		}
		// default value to be used.
		if strings.Contains(s, "default:") {
			columnDef.DefaultValue = s[8:]
		}
		// tag to everride column name created as per structure field.
		if strings.Contains(s, "column-name:") {
			columnDef.ColumnName = s[12:]
		}
	}
	return false, nil
}

// reflectGetTableDefinition Helper method get table definition directly from reflect param.
func (h Hera) reflectGetTableDefinition(v reflect.Value, getOnlyTableName bool) (tableDefinition, error) {
	// existsPK Signalizes if we already have a primary key field.
	existsPK := false
	result := tableDefinition{
		TableName:  reflectGetTableName(v.Type()),
		ColumnsDef: []Column{},
	}

	// loops through all structure fields! (not all passed fields)
	for i := 0; i < v.NumField(); i++ {
		// check if definition overrides table name
		if v.Type().Field(i).Name == "tableName" {
			result.TableName = strings.ToLower(v.Type().Field(i).Tag.Get("hera"))
			h.L.Debug("Overridden table name:", result.TableName)
		}

		if getOnlyTableName {
			return result, nil
		}

		// check if field definition exists in translation table. if not skip field.
		if getRDBMSType(v.Type().Field(i).Type.String(), false) == "" {
			h.L.Warnf("skipping field number: %v type: %s", i, v.Type().Field(i).Type.String())
			continue
		}

		column := Column{
			ColumnName: strings.ToLower(v.Type().Field(i).Name),
		}

		tag := v.Type().Field(i).Tag.Get("hera")
		h.L.Debug("Tag:", tag)

		// for the `hera:"-"` case
		ignoreField := false

		if len(tag) > 0 {
			var errParse error
			ignoreField, errParse = h.parseFieldTags(tag, &column, existsPK)
			if errParse != nil {
				return tableDefinition{}, errors.New("could not parse tag " + tag)
			}
		}
		if ignoreField {
			continue
		}

		// adding field type now that we defined the column data holder and parsed the tag.
		column.RDBMSType = getRDBMSType(v.Type().Field(i).Type.String(), column.PK)
		h.L.Debug("rdbms type: ", column.RDBMSType)

		// passing now complete column definition to final result slice.
		result.ColumnsDef = append(result.ColumnsDef, column)
	}
	return result, nil
}

// getTableDefinition Method gets table definition for model.
func (h Hera) getTableDefinition(model interface{}, getOnlyTableName bool) (tableDefinition, error) {
	if !h.isItPointer(model) {
		return tableDefinition{}, ErrorNotAPointer
	}
	return h.reflectGetTableDefinition(reflect.ValueOf(model).Elem(), getOnlyTableName)
}
