package pghera

/*
File concentrates model helpers specific to reflect.
*/

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

// ColumnShortData concentrates model field data.
type ColumnShortData struct {
	ColumnName string
	RDBMSType  reflect.Type
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

type User struct {
	tableName   struct{}    `hera:"theTableName"`
	ID          int64       `hera:"pk"`                                  // nolint
	name        string      `hera:"default:xx, index"`                   // nolint
	age         int         `hera:"required, column-name:theage, index"` // nolint
	isConnected bool        `hera:"default:true"`                        // nolint
	comment     string      `hera:"-"`                                   // nolint
	toSkip      interface{} // nolint
}

// isItPointer Checks if param is a pointer.
func (h Hera) isItPointer(model interface{}) bool {
	h.l.Debug("Passed data of type: ", reflect.TypeOf(model))
	return strings.HasPrefix(reflect.TypeOf(model).String(), "*")
}

// need a helper to create table columns short info. helper takes a pointer type.
func (h Hera) produceTableColumnShortData(model interface{}) ([]ColumnShortData, error) {
	if !h.isItPointer(model) {
		return []ColumnShortData{}, ErrorNotAPointer
	}

	result := []ColumnShortData{}
	for i := 0; i < reflect.TypeOf(model).Elem().NumField(); i++ {
		// append only types of fields in table translation and that are not ignored (tag "-") or primary key ( auto incremented ).
		// fields that are not passed would be considered with default values.
		fieldRoot := reflect.TypeOf(model).Elem().FieldByIndex([]int{i})

		h.l.Print("field type: ", fieldRoot.Type.String(), " - ", fieldRoot.Tag)

		if _, exists := (*newTranslationTable())[fieldRoot.Type.String()]; exists {
			if !strings.Contains(fmt.Sprintf("%v", fieldRoot.Tag), `"-"`) && !strings.Contains(fmt.Sprintf("%v", fieldRoot.Tag), `"pk"`) {
				result = append(result, ColumnShortData{
					ColumnName: fieldRoot.Name,
					RDBMSType:  fieldRoot.Type,
					Value:      reflect.ValueOf(model).Elem().FieldByIndex([]int{i}),
				})
			}
		}
	}
	return result, nil
}

// getTableName Gets table name from model. Use pointer like interface{}(&Model{}).
func (h Hera) getTableName(model interface{}) string {
	return reflectGetTableName(reflect.TypeOf(model))
}

// reflectGetTableName Helper in case the value is needed using reflection types.
func reflectGetTableName(v reflect.Type) string {
	if v.Kind() == reflect.Ptr {
		return inflection.Plural(strcase.ToSnake(v.Elem().Name()))
	}
	return inflection.Plural(strcase.ToSnake(v.Name()))
}

// reflectGetTableDefinition Helper method get table definition directly from reflect param.
func (h Hera) reflectGetTableDefinition(v reflect.Value, getOnlyTableName bool) (tableDefinition, error) {
	h.l.Debug("reflected value:", v)

	// existsPK Signalizes if we already have a primary key field.
	existsPK := false
	result := tableDefinition{
		TableName:  reflectGetTableName(v.Type()),
		ColumnsDef: []Column{},
	}
	for i := 0; i < v.NumField(); i++ {
		// check if definition overrides table name
		if v.Type().Field(i).Name == "tableName" {
			result.TableName = strings.ToLower(v.Type().Field(i).Tag.Get("hera"))
			h.l.Debug("Overriden table name:", result.TableName)
		}

		if getOnlyTableName {
			return result, nil
		}

		// check if field definition exists in translation table. if not skip field.
		rdbmsFieldType, exists := (*h.transTable)[v.Type().Field(i).Type.String()]
		if !exists {
			h.l.Warnf("skipping field number: %v type: %s", i, rdbmsFieldType)
			continue
		}

		column := Column{
			ColumnName: strcase.ToSnake(v.Type().Field(i).Name),
		}
		// adding field type now that we defined the column data holder.
		column.RDBMSType = rdbmsFieldType

		tag := v.Type().Field(i).Tag.Get("hera")
		h.l.Debug("Tag:", tag)

		// for the `hera:"-"` case
		ignoreField := false

		if len(tag) > 0 {
			tags := strings.Split(tag, ",")

			for _, tagS := range tags {
				s := strings.ToLower(strings.TrimSpace(tagS))

				if s == "-" {
					ignoreField = true
					break
				}
				if s == "pk" {
					if existsPK {
						return tableDefinition{}, errors.New("more than one primary key field detected. max is 1")
					}
					existsPK = true
					column.PK = true
				}
				if s == "unique" {
					column.Unique = true
				}
				if s == "index" {
					column.Index = true
				}
				if s == "required" {
					column.Required = true
				}
				if strings.Contains(s, "default:") {
					column.DefaultValue = s[8:]
				}
				if strings.Contains(s, "column-name:") {
					column.ColumnName = s[12:]
				}
			}
		}
		if ignoreField {
			continue
		}
		result.ColumnsDef = append(result.ColumnsDef, column)
	}
	return result, nil
}

// getTableDefinition Method gets table definition for model.
func (h Hera) getTableDefinition(model interface{}, getOnlyTableName bool) (tableDefinition, error) {
	if !h.isItPointer(model) {
		return tableDefinition{}, ErrorNotAPointer
	}
	val := reflect.ValueOf(model).Elem()
	h.l.Debug("val:", val)

	return h.reflectGetTableDefinition(val, getOnlyTableName)
}
