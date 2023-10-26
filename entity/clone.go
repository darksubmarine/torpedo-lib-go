package entity

import (
	"fmt"
	"reflect"
	"strings"
)

func Clone(from interface{}, to interface{}) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recoverd from %s", r)
		}
	}()

	toTypeOf := reflect.TypeOf(to)
	toValueOf := reflect.ValueOf(to)
	fromValueOf := reflect.ValueOf(from)

	iterateEntity(toTypeOf, &toValueOf, &fromValueOf)

	return nil
}

func iterateEntity(toTypeOf reflect.Type, toValueOf *reflect.Value, fromValueOf *reflect.Value) {
	toNumMethods := toTypeOf.NumMethod()
	for i := 0; i < toNumMethods; i++ {
		if strings.HasPrefix(toTypeOf.Method(i).Name, "Set") {
			readMethod, _ := strings.CutPrefix(toTypeOf.Method(i).Name, "Set")
			if fromValueOf.MethodByName(readMethod).Kind() != reflect.Invalid {
				if reflect.TypeOf(fromValueOf.MethodByName(readMethod).Interface()).NumOut() == 0 {
					continue
				}
				values := fromValueOf.MethodByName(readMethod).Call([]reflect.Value{})
				toValueOf.Method(i).Call(values)
			}
		}
	}
}
