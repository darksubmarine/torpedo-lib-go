package sql_utils

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

func InsertStatementFromDMO(tableName string, dmo interface{}, skip ...string) string {
	dmoTypeOf := reflect.TypeOf(dmo).Elem()

	fields := iterateDMO(dmoTypeOf, skipMap(skip...))

	_colName := bytes.NewBufferString("")
	_colRef := bytes.NewBufferString("")
	_lastFieldPosition := len(fields) - 1
	for i, f := range fields {
		_colName.WriteString(f)
		_colRef.WriteString(fmt.Sprintf(":%s", f))

		if i < _lastFieldPosition {
			_colName.WriteString(",")
			_colRef.WriteString(",")
		}
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, _colName.String(), _colRef.String())
}

func UpdateStatementFromDMO(tableName string, dmo interface{}, skip ...string) string {
	dmoTypeOf := reflect.TypeOf(dmo).Elem()

	fields := iterateDMO(dmoTypeOf, skipMap(skip...))

	_set := bytes.NewBufferString("")
	_lastFieldPosition := len(fields) - 1
	for i, f := range fields {
		if f == "id" {
			continue
		}
		_set.WriteString(fmt.Sprintf("%s = :%s", f, f))

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

func iterateDMO(dmoTypeOf reflect.Type, skip map[string]struct{}) []string {
	fields := make([]string, 0)
	fromTypeOfNumField := dmoTypeOf.NumField()
	for i := 0; i < fromTypeOfNumField; i++ {
		if dmoTypeOf.Field(i).Type.Kind() == reflect.Struct {
			if fls := iterateDMO(dmoTypeOf.Field(i).Type, skip); len(fls) > 0 {
				fields = append(fields, fls...)
			}
		} else if name := dmoTypeOf.Field(i).Name; strings.HasSuffix(name, "_") {
			if _, skipIt := skip[name]; skipIt {
				continue
			}

			fields = append(fields, dmoTypeOf.Field(i).Tag.Get("db"))
		}
	}
	return fields
}
