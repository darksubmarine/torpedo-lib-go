package entity

import (
	"fmt"
	"github.com/darksubmarine/torpedo-lib-go/ptr"
	"reflect"
	"strings"
)

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

func From(from interface{}, entity interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recoverd from %s", r)
		}
	}()

	fromTypeOf := reflect.TypeOf(from).Elem()
	fromValueOf := reflect.ValueOf(from).Elem()
	entityValueOf := reflect.ValueOf(entity)
	return iterateFrom(fromTypeOf, fromValueOf, &entityValueOf)
}

func getValue(v reflect.Value) reflect.Value {
	var kind reflect.Kind
	if v.Kind() == reflect.Pointer {
		kind = v.Elem().Kind()
	} else {
		kind = v.Kind()
	}

	var valueOf reflect.Value
	switch kind {
	case reflect.String:
		switch v.Interface().(type) {
		case string:
			valueOf = reflect.ValueOf(v.Interface().(string))
		case *string:
			valueOf = reflect.ValueOf(ptr.ToString(v.Interface().(*string)))
		}

	case reflect.Int:
		switch v.Interface().(type) {
		case int:
			valueOf = reflect.ValueOf(v.Interface().(int))
		case *int:
			valueOf = reflect.ValueOf(ptr.ToInt(v.Interface().(*int)))
		}

	case reflect.Int64:
		switch v.Interface().(type) {
		case int64:
			valueOf = reflect.ValueOf(v.Interface().(int64))
		case *int64:
			valueOf = reflect.ValueOf(ptr.ToInt64(v.Interface().(*int64)))
		}

	case reflect.Bool:
		switch v.Interface().(type) {
		case bool:
			valueOf = reflect.ValueOf(v.Interface().(bool))
		case *bool:
			valueOf = reflect.ValueOf(ptr.ToBool(v.Interface().(*bool)))
		}

	case reflect.Float32:
		switch v.Interface().(type) {
		case float32:
			valueOf = reflect.ValueOf(v.Interface().(float32))
		case *float32:
			valueOf = reflect.ValueOf(ptr.ToFloat32(v.Interface().(*float32)))
		}

	case reflect.Float64:
		switch v.Interface().(type) {
		case float64:
			valueOf = reflect.ValueOf(v.Interface().(float64))
		case *float64:
			valueOf = reflect.ValueOf(ptr.ToFloat64(v.Interface().(*float64)))
		}

	case reflect.Slice:
		if v.Len() > 0 {
			switch v.Index(0).Kind() {
			case reflect.Int:
				valueOf = reflect.ValueOf(v.Interface().([]int))
			case reflect.Int64:
				valueOf = reflect.ValueOf(v.Interface().([]int64))
			case reflect.Float32:
				valueOf = reflect.ValueOf(v.Interface().([]float32))
			case reflect.Float64:
				valueOf = reflect.ValueOf(v.Interface().([]float64))
			case reflect.String:
				valueOf = reflect.ValueOf(v.Interface().([]string))
			case reflect.Interface:
				valueOf = reflect.ValueOf(v.Interface().([]interface{}))
			}
		}
	}

	return valueOf
}

func iterateFrom(fromTypeOf reflect.Type, fromValueOf reflect.Value, entityValueOf *reflect.Value) error {
	fromTypeOfNumField := fromTypeOf.NumField()
	for i := 0; i < fromTypeOfNumField; i++ {
		if fromTypeOf.Field(i).Type.Kind() == reflect.Struct {
			if err := iterateFrom(fromTypeOf.Field(i).Type, fromValueOf.Field(i), entityValueOf); err != nil {
				return err
			}
		} else if name := fromTypeOf.Field(i).Name; strings.HasSuffix(name, "_") {
			varName, _ := strings.CutSuffix(name, "_")
			methodName := fmt.Sprintf("Set%s", varName)
			if entityValueOf.MethodByName(methodName).Kind() != reflect.Invalid {
				if val := getValue(fromValueOf.Field(i)); val.Kind() != reflect.Invalid {

					var valueToSet = val
					// Checking for encrypted fields
					if tagVal, ok := fromTypeOf.Field(i).Tag.Lookup("tpdo"); ok {
						if tagVal == "encrypted" {
							if fromValueOf.MethodByName("DecryptString").Kind() != reflect.Invalid {
								vals := fromValueOf.MethodByName("DecryptString").Call([]reflect.Value{val})
								valueToSet = vals[0] // TODO handle error
							}
						}
					}

					// setting value
					res := entityValueOf.MethodByName(methodName).Call([]reflect.Value{valueToSet})

					// checking if result has an error
					for _, r := range res {

						if r.Interface() != nil && r.Type().Implements(errorInterface) {
							return r.Interface().(error)
						}

					}
				}
			}
		}
	}
	return nil
}
