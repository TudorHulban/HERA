package pghera

import (
	"log"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

// Column Concentrates table column data.
type Column struct {
	// If not required column could be null.
	Required   bool
	Unique     bool
	Index      bool
	ColumnName string
	DataType   string
	Default    string
}

type User struct {
	ID          int64  `hera:"pk"`
	name        string `hera:"default:xx"`
	age         int    `hera:"required"`
	isConnected bool   `hera:"default:true"`
	comment     string `hera:"-"`
}

// getTableName Gets table name from model. Use pointer like interface{}(&Model{}).
func (h Hera) getTableName(model interface{}) string {
	if t := reflect.TypeOf(model); t.Kind() == reflect.Ptr {
		return inflection.Plural(strcase.ToSnake(t.Elem().Name()))
	} else {
		return inflection.Plural(strcase.ToSnake(t.Name()))
	}
}

func allowedField(fieldType string) bool {
	allowedTypes := []string{"string", "*string", "int", "int64", "float64", "*float64", "bool", "*bool"}
	for _, v := range allowedTypes {
		if fieldType == v {
			return true
		}
	}
	return false
}

func (h Hera) getTableColumns(model interface{}) []Column {
	val := reflect.ValueOf(model).Elem()
	h.l.Debug("val:", val)

	var result []Column
	for i := 0; i < val.NumField(); i++ {
		fieldType := val.Type().Field(i).Type.String()
		if !allowedField(fieldType) {
			log.Println("skipping fieldType:", fieldType)
			continue
		}

		column := Column{
			ColumnName: strcase.ToSnake(val.Type().Field(i).Name),
		}

		tag := val.Type().Field(i).Tag.Get("hera")
		log.Println("Tag:", tag)

		if len(tag) > 0 {
			tags := strings.Split(tag, ",")

			for _, tagS := range tags {
				s := strings.ToLower(strings.TrimSpace(tagS))

				if s == "-" {
					continue
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

		if fieldType == "string" || fieldType == "*string" {
			column.DataType = "text"
		}
		if fieldType == "int64" || fieldType == "int" {
			column.DataType = "bigint"
		}
		if fieldType == "float64" || fieldType == "*float64" {
			column.DataType = "numeric"
		}
		if fieldType == "bool" || fieldType == "*bool" {
			column.DataType = "boolean"
		}
		result = append(result, column)
	}
	return result
}
