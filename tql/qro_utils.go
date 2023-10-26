package tql

import (
	"reflect"
)

func FieldValue(qro interface{}, field string) interface{} {
	qroTypeOf := reflect.TypeOf(qro).Elem()
	qroValueOf := reflect.ValueOf(qro).Elem()

	return iterateQROFields(qroTypeOf, qroValueOf, field)
}

func iterateQROFields(qroTypeOf reflect.Type, qroValueOf reflect.Value, field string) interface{} {
	qroTypeOfNumField := qroTypeOf.NumField()
	for i := 0; i < qroTypeOfNumField; i++ {
		if qroTypeOf.Field(i).Type.Kind() == reflect.Struct {
			if ret := iterateQROFields(qroTypeOf.Field(i).Type, qroValueOf.Field(i), field); ret != nil {
				return ret
			}

		} else if name := qroTypeOf.Field(i).Name; name == field {
			if qroTypeOf.Field(i).Type.Kind() == reflect.Pointer {
				return reflect.Indirect(qroValueOf.Field(i)).Interface()
			} else {
				return qroValueOf.Field(i).Interface()
			}
		}
	}

	return nil
}
