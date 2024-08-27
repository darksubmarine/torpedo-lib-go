package entity

import (
	"encoding/json"
	"fmt"
	"github.com/darksubmarine/torpedo-lib-go/ptr"
	"github.com/darksubmarine/torpedo-lib-go/storage/sql_utils/data_type"
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
	rootEntityValueOf := reflect.ValueOf(entity)
	entityValueOf := reflect.ValueOf(entity).Elem()
	entityTypeOf := reflect.TypeOf(entity).Elem()

	copyFields := map[string]struct{}{}
	for _, v := range field {
		copyFields[v] = struct{}{}
	}

	iterateToEntity(entityTypeOf, &entityValueOf, &rootEntityValueOf, toTypeOf, &toValueOf, copyFields)

	return nil
}

func iterateToEntity(etyTypeOf reflect.Type, etyValueOf *reflect.Value, rootEtyValueOf *reflect.Value, toTypeOf reflect.Type, toValueOf *reflect.Value, onlyCopyThisFields map[string]struct{}) {

	toTypeOfString := strings.Split(toTypeOf.String(), ".")
	toPkg := toTypeOfString[0]
	toObject := toTypeOfString[1]

	etyTypeOfNumField := etyTypeOf.NumField()
	for i := 0; i < etyTypeOfNumField; i++ {
		if etyTypeOf.Field(i).Type.Kind() == reflect.Pointer && reflect.Indirect(etyValueOf.Field(i)).Kind() == reflect.Struct {
			_etyValueOf := etyValueOf.Field(i)
			iterateToEntity(reflect.Indirect(etyValueOf.Field(i)).Type(), &_etyValueOf, rootEtyValueOf, toTypeOf, toValueOf, onlyCopyThisFields)
		} else if etyTypeOf.Field(i).Type.Kind() == reflect.Struct {
			_etyValueOf := etyValueOf.Field(i)
			iterateToEntity(etyTypeOf.Field(i).Type, &_etyValueOf, rootEtyValueOf, toTypeOf, toValueOf, onlyCopyThisFields)
		} else {

			fName := etyTypeOf.Field(i).Name

			var toInput = false
			fMeta := readFieldMetadata(etyTypeOf.Field(i))
			toFieldName := FieldNameToCode(fName)

			if toObject == "EntityQRO" {
				if fMeta.qro != "" {
					toFieldName = fMeta.qro
				}
			} else {
				switch toPkg {
				// inputs
				case "http", "gin", "dto":
					toInput = true
					if fMeta.dto.http != "" {
						toFieldName = fMeta.dto.http
					}

				// outputs
				case "memory":
					if fMeta.dmo.memory != "" {
						toFieldName = fMeta.dmo.memory
					}
				case "redis":
					if fMeta.dmo.redis != "" {
						toFieldName = fMeta.dmo.redis
					}
				case "mongodb":
					if fMeta.dmo.mongodb != "" {
						toFieldName = fMeta.dmo.mongodb
					}
				case "sql":
					if fMeta.dmo.sql != "" {
						toFieldName = fMeta.dmo.sql
					}

				// testing
				case "entity_test":
					if fMeta.dmo.memory != "" {
						toFieldName = fMeta.dmo.memory
					}
				}
			}

			// TODO Warning with this logic...
			// this is used at QRO  when the query has a projection to copy only (from Entity TO QRO) projected fields
			if len(onlyCopyThisFields) > 0 {
				if _, exists := onlyCopyThisFields[toFieldName]; !exists {
					continue
				}
			}

			methodName := fMeta.getter
			if rootEtyValueOf.MethodByName(methodName).Kind() != reflect.Invalid {
				if reflect.TypeOf(rootEtyValueOf.MethodByName(methodName).Interface()).NumOut() == 0 {
					continue
				}

				values := rootEtyValueOf.MethodByName(methodName).Call([]reflect.Value{})
				valueToSet := values[0]

				if fMeta.encrypted {
					if !toInput {
						if toValueOf.MethodByName("EncryptString").Kind() != reflect.Invalid {
							vals := toValueOf.MethodByName("EncryptString").Call([]reflect.Value{valueToSet})
							valueToSet = vals[0]
						}
					}
				}

				if toFieldV := toValueOf.FieldByName(toFieldName); toFieldV.IsValid() {
					setValue(&toFieldV, &valueToSet)
				}
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

	case data_type.JsonArrayFloat, data_type.JsonArrayInteger, data_type.JsonArrayString, data_type.JsonArrayDate, data_type.JsonArrayBoolean:
		if jsonStr, err := json.Marshal(value.Interface()); err == nil {
			str := string(jsonStr)
			switch dv.(type) {
			case data_type.JsonArrayFloat:
				dest.Set(reflect.ValueOf(data_type.JsonArrayFloat(&str)))
			case data_type.JsonArrayInteger:
				dest.Set(reflect.ValueOf(data_type.JsonArrayInteger(&str)))
			case data_type.JsonArrayString:
				dest.Set(reflect.ValueOf(data_type.JsonArrayString(&str)))
			case data_type.JsonArrayDate:
				dest.Set(reflect.ValueOf(data_type.JsonArrayDate(&str)))
			case data_type.JsonArrayBoolean:
				dest.Set(reflect.ValueOf(data_type.JsonArrayBoolean(&str)))

			}
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
	case []int32:
		if value.Len() > 0 {
			switch value.Index(0).Kind() {
			case reflect.Int32:
				dest.Set(reflect.ValueOf(value.Interface().([]int32)))
			}
		}
	case []int64:
		if value.Len() > 0 {
			switch value.Index(0).Kind() {
			case reflect.Int64:
				dest.Set(reflect.ValueOf(value.Interface().([]int64)))
			}
		}
	case []float64:
		if value.Len() > 0 {
			switch value.Index(0).Kind() {
			case reflect.Float64:
				dest.Set(reflect.ValueOf(value.Interface().([]float64)))
			}
		}
	case []float32:
		if value.Len() > 0 {
			switch value.Index(0).Kind() {
			case reflect.Float32:
				dest.Set(reflect.ValueOf(value.Interface().([]float32)))
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
