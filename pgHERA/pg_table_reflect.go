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
func getTableName(model interface{}) string {
	if t := reflect.TypeOf(model); t.Kind() == reflect.Ptr {
		return inflection.Plural(strcase.ToSnake(t.Elem().Name()))
	} else {
		return inflection.Plural(strcase.ToSnake(t.Name()))
	}
}

func allowedField(fieldType string) bool {
	allowedTypes := []string{"string", "int", "int64", "float64", "bool"}
	for _, v := range allowedTypes {
		if fieldType == v {
			return true
		}
	}
	return false
}

func (h Hera) getTableColumns(model interface{}) []Column {
	var result []Column

	val := reflect.ValueOf(model).Elem()
	h.l.Debug("val:", val)

	for i := 0; i < val.NumField(); i++ {
		fieldType := val.Type().Field(i).Type.String()
		if !allowedField(fieldType) {
			log.Println("fieldType:", fieldType)
			continue
		}

		column := Column{
			ColumnName: strcase.ToSnake(val.Type().Field(i).Name),
		}

		tag := val.Type().Field(i).Tag.Get("hera")

		log.Println("Tag:", tag)
		ignoreField := false

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
				if strings.Contains(s, "default:") {
					column.Default = s[8:]
				}

			}
		}

		if fieldType == "string" {
			column.DataType = "text"
			if column.Default == "" {
				column.Default = "''"
			}
		}
		if fieldType == "*string" {
			column.DataType = "text"
		}
		if fieldType == "int64" || fieldType == "int" {
			column.DataType = "bigint"
		}
		if fieldType == "float64" || fieldType == "*float64" {
			column.DataType = "numeric"
			column.Default = "0.00"
		}
		if fieldType == "bool" || fieldType == "*bool" {
			column.DataType = "boolean"
			if column.Default == "" {
				column.Default = "false"
			}
		}
		res = append(res, column)
	}
	return res
}
