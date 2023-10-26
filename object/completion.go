package object

import (
	"fmt"
	"github.com/darksubmarine/torpedo-lib-go/data_struct"
	"reflect"
)

// IsComplete iterates all struct fields by reflection and check if values are nil.
func IsComplete(obj interface{}, skip ...string) (bool, error) {

	var objTypeOf reflect.Type
	var objValueOf reflect.Value

	if reflect.TypeOf(obj).Kind() == reflect.Pointer {
		objTypeOf = reflect.TypeOf(obj).Elem()
		objValueOf = reflect.ValueOf(obj).Elem()
	} else {
		objTypeOf = reflect.TypeOf(obj)
		objValueOf = reflect.ValueOf(obj)
	}

	return iterateFields(objTypeOf, objValueOf, data_struct.SkipMap(skip...))
}

func iterateFields(objTypeOf reflect.Type, objValueOf reflect.Value, skip map[string]struct{}) (bool, error) {
	objTypeOfNumField := objTypeOf.NumField()
	for i := 0; i < objTypeOfNumField; i++ {
		if objTypeOf.Field(i).Type.Kind() == reflect.Struct {
			if ok, err := iterateFields(objTypeOf.Field(i).Type, objValueOf.Field(i), skip); err != nil {
				return false, err
			} else if !ok {
				return false, nil
			}
		} else {

			if _, ok := skip[objTypeOf.Field(i).Name]; ok {
				continue
			}

			if objValueOf.Field(i).Kind() == reflect.Pointer || objValueOf.Field(i).Kind() == reflect.Slice {
				if objValueOf.Field(i).IsNil() {
					return false, nil // value is not set
				}
			} else {
				return false, fmt.Errorf("the object fields must be pointers and got %s(%s)", objTypeOf.Field(i).Name, objValueOf.Field(i).Kind().String())
			}
		}
	}

	return true, nil
}
