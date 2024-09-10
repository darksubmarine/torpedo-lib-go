package sql_utils

import (
	"bytes"
	"fmt"
	"github.com/darksubmarine/torpedo-lib-go/entity"
	"reflect"
	"strings"
)

func InsertStatementFromDMO(driverName string, tableName string, dmo interface{}, metadata map[string]*entity.FieldMetadata, skip ...string) string {
	dmoTypeOf := reflect.TypeOf(dmo).Elem()

	fields := iterateDMO(dmoTypeOf, metadata, skipMap(skip...))

	_colName := bytes.NewBufferString("")
	_colRef := bytes.NewBufferString("")
	_lastFieldPosition := len(fields) - 1
	for i, f := range fields {
		if driverName == "mysql" {
			_colName.WriteString(fmt.Sprintf("`%s`", f))
		} else {
			_colName.WriteString(fmt.Sprintf("%s", f))
		}

		_colRef.WriteString(fmt.Sprintf(":%s", f))

		if i < _lastFieldPosition {
			_colName.WriteString(",")
			_colRef.WriteString(",")
		}
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, _colName.String(), _colRef.String())
}

func UpdateStatementFromDMO(driverName string, tableName string, dmo interface{}, metadata map[string]*entity.FieldMetadata, skip ...string) string {
	dmoTypeOf := reflect.TypeOf(dmo).Elem()

	fields := iterateDMO(dmoTypeOf, metadata, skipMap(skip...))

	_set := bytes.NewBufferString("")
	_lastFieldPosition := len(fields) - 1
	for i, f := range fields {
		if f == "id" {
			continue
		}

		if driverName == "mysql" {
			_set.WriteString(fmt.Sprintf("`%s` = :%s", f, f))
		} else {
			_set.WriteString(fmt.Sprintf("%s = :%s", f, f))
		}

		if i < _lastFieldPosition {
			_set.WriteString(", ")
		}
	}

	return fmt.Sprintf("UPDATE %s SET %s WHERE id = :id", tableName, _set)
}

func skipMap(skip ...string) map[string]struct{} {
	sm := map[string]struct{}{}
	for _, v := range skip {
		sm[v] = struct{}{}
	}
	return sm
}

func iterateDMO(dmoTypeOf reflect.Type, metadata map[string]*entity.FieldMetadata, skip map[string]struct{}) []string {
	fields := make([]string, 0)
	fromTypeOfNumField := dmoTypeOf.NumField()
	for i := 0; i < fromTypeOfNumField; i++ {
		if dmoTypeOf.Field(i).Type.Kind() == reflect.Struct {
			if fls := iterateDMO(dmoTypeOf.Field(i).Type, metadata, skip); len(fls) > 0 {
				fields = append(fields, fls...)
			}
		} else {
			var fName = dmoTypeOf.Field(i).Name // redding field name from DMO
			var addToFields = false             // this field should be added to update doc?

			if !strings.HasSuffix(fName, "_") { // If field name not follows the naming convention "FieldName_" check metadata
				for _, mdata := range metadata {
					if mdata.DmoSqlName() == fName {
						addToFields = true
						break
					}
				}
			} else { // field name with naming should be added
				addToFields = true
			}

			if !addToFields {
				continue
			}

			if _, toSkip := skip[fName]; toSkip {
				continue
			}

			fields = append(fields, dmoTypeOf.Field(i).Tag.Get("db"))
		}
	}
	return fields
}
