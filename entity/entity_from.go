package entity

import (
	"encoding/json"
	"fmt"
	"github.com/darksubmarine/torpedo-lib-go/ptr"
	"github.com/darksubmarine/torpedo-lib-go/storage/sql_utils/data_type"
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
	rootEntityValueOf := reflect.ValueOf(entity)
	entityValueOf := reflect.ValueOf(entity).Elem()
	entityTypeOf := reflect.TypeOf(entity).Elem()

	return iterateFromEntity(fromTypeOf, &fromValueOf, entityTypeOf, &entityValueOf, &rootEntityValueOf)
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

		case data_type.JsonArrayFloat:
			str := v.Interface().(data_type.JsonArrayFloat)
			var obj []float64
			json.Unmarshal([]byte(*str), &obj)
			valueOf = reflect.ValueOf(obj)
		case data_type.JsonArrayInteger:
			str := v.Interface().(data_type.JsonArrayInteger)
			var obj []int64
			json.Unmarshal([]byte(*str), &obj)
			valueOf = reflect.ValueOf(obj)
		case data_type.JsonArrayString:
			str := v.Interface().(data_type.JsonArrayString)
			var obj []string
			json.Unmarshal([]byte(*str), &obj)
			valueOf = reflect.ValueOf(obj)
		case data_type.JsonArrayDate:
			str := v.Interface().(data_type.JsonArrayDate)
			var obj []int64
			json.Unmarshal([]byte(*str), &obj)
			valueOf = reflect.ValueOf(obj)
		case data_type.JsonArrayBoolean:
			str := v.Interface().(data_type.JsonArrayBoolean)
			var obj []bool
			json.Unmarshal([]byte(*str), &obj)
			valueOf = reflect.ValueOf(obj)
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

func iterateFromEntity(fromTypeOf reflect.Type, fromValueOf *reflect.Value, etyTypeOf reflect.Type, etyValueOf *reflect.Value, rootEtyValueOf *reflect.Value) error {
	etyTypeOfNumField := etyTypeOf.NumField()
	for i := 0; i < etyTypeOfNumField; i++ {

		fMeta := readFieldMetadata(etyTypeOf.Field(i))
		if fMeta.IsRelationship() {
			continue
		}

		if etyTypeOf.Field(i).Type.Kind() == reflect.Pointer && reflect.Indirect(etyValueOf.Field(i)).Kind() == reflect.Struct {
			_etyValueOf := etyValueOf.Field(i)
			if err := iterateFromEntity(fromTypeOf, fromValueOf, reflect.Indirect(etyValueOf.Field(i)).Type(), &_etyValueOf, rootEtyValueOf); err != nil {
				return err
			}
		} else if etyTypeOf.Field(i).Type.Kind() == reflect.Struct {
			_etyValueOf := etyValueOf.Field(i)
			if err := iterateFromEntity(fromTypeOf, fromValueOf, etyTypeOf.Field(i).Type, &_etyValueOf, rootEtyValueOf); err != nil {
				return err
			}
		} else {
			fName := etyTypeOf.Field(i).Name
			var fromOutput = false
			fromFieldName := FieldNameToCode(fName)
			fromPkg := strings.Split(fromTypeOf.String(), ".")[0]

			switch fromPkg {
			// inputs
			case "http", "gin", "dto":
				if fMeta.dto.http != "" {
					fromFieldName = fMeta.dto.http
				}

			// outputs
			case "memory":
				fromOutput = true
				if fMeta.dmo.memory != "" {
					fromFieldName = fMeta.dmo.memory
				}
			case "redis":
				fromOutput = true
				if fMeta.dmo.redis != "" {
					fromFieldName = fMeta.dmo.redis
				}
			case "mongodb":
				fromOutput = true
				if fMeta.dmo.mongodb != "" {
					fromFieldName = fMeta.dmo.mongodb
				}
			case "sql":
				fromOutput = true
				if fMeta.dmo.sql != "" {
					fromFieldName = fMeta.dmo.sql
				}

			// testing
			case "entity_test":
				fromOutput = true
				if fMeta.dto.http != "" {
					fromFieldName = fMeta.dto.http
				}
			}

			methodName := fMeta.setter
			if rootEtyValueOf.MethodByName(methodName).Kind() != reflect.Invalid {
				if val := getValue(fromValueOf.FieldByName(fromFieldName)); val.Kind() != reflect.Invalid {

					var valueToSet = val
					// Checking for encrypted fields
					if fMeta.encrypted {
						if fromOutput {
							if fromValueOf.MethodByName("DecryptString").Kind() != reflect.Invalid {
								vals := fromValueOf.MethodByName("DecryptString").Call([]reflect.Value{val})
								valueToSet = vals[0] // TODO handle error
							}
						}
					}

					// setting value
					res := rootEtyValueOf.MethodByName(methodName).Call([]reflect.Value{valueToSet})

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
