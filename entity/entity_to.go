package entity

import (
	"fmt"
	"github.com/darksubmarine/torpedo-lib-go/ptr"
	"reflect"
	"strings"
)

func To(entity interface{}, to interface{}, field ...string) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recoverd from %s", r)
		}
	}()

	toTypeOf := reflect.TypeOf(to).Elem()
	toValueOf := reflect.ValueOf(to).Elem()
	entityValueOf := reflect.ValueOf(entity)

	copyFields := map[string]struct{}{}
	for _, v := range field {
		copyFields[v] = struct{}{}
	}

	iterateTo(toTypeOf, &toValueOf, &entityValueOf, copyFields)

	return nil
}

func iterateTo(toTypeOf reflect.Type, toValueOf *reflect.Value, entityValueOf *reflect.Value, fields map[string]struct{}) {
	toTypeOfNumField := toTypeOf.NumField()
	for i := 0; i < toTypeOfNumField; i++ {
		if toTypeOf.Field(i).Type.Kind() == reflect.Struct {
			v := toValueOf.Field(i)
			iterateTo(toTypeOf.Field(i).Type, &v, entityValueOf, fields)
		} else if name := toTypeOf.Field(i).Name; strings.HasSuffix(name, "_") {

			// TODO Warning with this logic
			if len(fields) > 0 {
				if _, exists := fields[name]; !exists {
					continue
				}
			}

			var methodName string
			if val, ok := toTypeOf.Field(i).Tag.Lookup("read_method"); ok {
				methodName = val
			} else {
				varName, _ := strings.CutSuffix(name, "_")
				methodName = fmt.Sprintf("%s", varName)
			}

			if entityValueOf.MethodByName(methodName).Kind() != reflect.Invalid {

				if reflect.TypeOf(entityValueOf.MethodByName(methodName).Interface()).NumOut() == 0 {
					continue
				}

				values := entityValueOf.MethodByName(methodName).Call([]reflect.Value{})
				valueToSet := values[0]
				toValueOfField := toValueOf.Field(i)

				// Checking for encrypted fields
				if tagVal, ok := toTypeOf.Field(i).Tag.Lookup("tpdo"); ok {
					if tagVal == "encrypted" {
						if toValueOf.MethodByName("EncryptString").Kind() != reflect.Invalid {
							vals := toValueOf.MethodByName("EncryptString").Call([]reflect.Value{valueToSet})
							valueToSet = vals[0]
						}
					}
				}

				setValue(&toValueOfField, &valueToSet)
			}
		}
	}
}

func setValue(dest *reflect.Value, value *reflect.Value) {
	var dv = dest.Interface()

	switch dv.(type) {
	case *string:
		switch v := value.Interface().(type) {
		case string:
			dest.Set(reflect.ValueOf(ptr.String(v)))
		case *string:
			dest.Set(reflect.ValueOf(v))
		}

	case string:
		switch v := value.Interface().(type) {
		case string:
			dest.Set(reflect.ValueOf(v))
		case *string:
			dest.Set(reflect.ValueOf(ptr.ToString(v)))
		}

	case *int:
		switch v := value.Interface().(type) {
		case int:
			dest.Set(reflect.ValueOf(ptr.Int(v)))
		case *int:
			dest.Set(reflect.ValueOf(v))
		}

	case int:
		switch v := value.Interface().(type) {
		case int:
			dest.Set(reflect.ValueOf(v))
		case *int:
			dest.Set(reflect.ValueOf(ptr.ToInt(v)))
		}

	case *int64:
		switch v := value.Interface().(type) {
		case int64:
			dest.Set(reflect.ValueOf(ptr.Int64(v)))
		case *int64:
			dest.Set(reflect.ValueOf(v))
		}

	case int64:
		switch v := value.Interface().(type) {
		case int64:
			dest.Set(reflect.ValueOf(v))
		case *int64:
			dest.Set(reflect.ValueOf(ptr.ToInt64(v)))
		}

	case *bool:
		switch v := value.Interface().(type) {
		case bool:
			dest.Set(reflect.ValueOf(ptr.Bool(v)))
		case *bool:
			dest.Set(reflect.ValueOf(v))
		}

	case bool:
		switch v := value.Interface().(type) {
		case bool:
			dest.Set(reflect.ValueOf(v))
		case *bool:
			dest.Set(reflect.ValueOf(ptr.ToBool(v)))
		}

	case *float32:
		switch v := value.Interface().(type) {
		case float32:
			dest.Set(reflect.ValueOf(ptr.Float32(v)))
		case *float32:
			dest.Set(reflect.ValueOf(v))
		}

	case float32:
		switch v := value.Interface().(type) {
		case float32:
			dest.Set(reflect.ValueOf(v))
		case *float32:
			dest.Set(reflect.ValueOf(ptr.ToFloat32(v)))
		}

	case *float64:
		switch v := value.Interface().(type) {
		case float64:
			dest.Set(reflect.ValueOf(ptr.Float64(v)))
		case *float64:
			dest.Set(reflect.ValueOf(v))
		}

	case float64:
		switch v := value.Interface().(type) {
		case float64:
			dest.Set(reflect.ValueOf(v))
		case *float64:
			dest.Set(reflect.ValueOf(ptr.ToFloat64(v)))
		}

	case []int:
		if value.Len() > 0 {
			switch value.Index(0).Kind() {
			case reflect.Int:
				dest.Set(reflect.ValueOf(value.Interface().([]int)))
			}
		}
	case []string:
		if value.Len() > 0 {
			switch value.Index(0).Kind() {
			case reflect.String:
				dest.Set(reflect.ValueOf(value.Interface().([]string)))
			}
		}
	case []interface{}:
		if value.Len() > 0 {
			switch value.Index(0).Kind() {
			case reflect.Interface:
				dest.Set(reflect.ValueOf(value.Interface().([]interface{})))
			}
		}
	}
}
