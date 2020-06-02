package pghera

import (
	"errors"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

// Column Concentrates table column data.
type Column struct {
	// If primary key.
	PK bool
	// If not required column could be null.
	Required   bool
	Unique     bool
	Index      bool
	ColumnName string
	RDBMSType  string
	Default    string
}

type User struct {
	ID          int64       `hera:"pk"`           // nolint
	name        string      `hera:"default:xx"`   // nolint
	age         int         `hera:"required"`     // nolint
	isConnected bool        `hera:"default:true"` // nolint
	comment     string      `hera:"-"`            // nolint
	toSkip      interface{} // nolint
}

// getTableName Gets table name from model. Use pointer like interface{}(&Model{}).
func (h Hera) getTableName(model interface{}) string {
	if t := reflect.TypeOf(model); t.Kind() == reflect.Ptr {
		return inflection.Plural(strcase.ToSnake(t.Elem().Name()))
	} else { // nolint
		return inflection.Plural(strcase.ToSnake(t.Name()))
	}
}

func (h Hera) getTableColumns(model interface{}) ([]Column, error) {
	val := reflect.ValueOf(model).Elem()
	h.l.Debug("val:", val)

	// existsPK Signalizes if we already have a primary key field.
	existsPK := false
	var result []Column

	for i := 0; i < val.NumField(); i++ {
		rdbmsFieldType, exists := (*h.transTable)[val.Type().Field(i).Type.String()]
		if !exists {
			h.l.Warnf("skipping field number: %v", i)
			continue
		}

		column := Column{
			ColumnName: strcase.ToSnake(val.Type().Field(i).Name),
		}
		// adding field type now that we defined the column data holder.
		column.RDBMSType = rdbmsFieldType

		tag := val.Type().Field(i).Tag.Get("hera")
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
						return []Column{}, errors.New("more than one primary key field detected. max is 1")
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
					column.Default = s[8:]
				}
			}
		}
		if ignoreField {
			continue
		}
		result = append(result, column)
	}
	return result, nil
}
