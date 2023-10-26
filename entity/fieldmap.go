package entity

import (
	"reflect"
)

type FieldMap map[string]string

func (fm FieldMap) Add(name string, kind string) {
	fm[name] = kind
}

func (fm FieldMap) HasField(name string) bool {
	_, ok := fm[name]
	return ok
}

func (fm FieldMap) FieldType(name string) (string, bool) {
	if val, ok := fm[name]; ok {
		return val, true
	}

	return "", false
}

func ToFieldMap(ety interface{}, skip ...string) FieldMap {
	fm := FieldMap{}

	etyTypeOf := reflect.TypeOf(ety).Elem()
	etyValueOf := reflect.ValueOf(ety).Elem()

	skipMap := map[string]struct{}{}
	for _, v := range skip {
		skipMap[v] = struct{}{}
	}

	iterateEntityFields(etyTypeOf, etyValueOf, &fm, skipMap)

	return fm
}

func iterateEntityFields(etyTypeOf reflect.Type, etyValueOf reflect.Value, fieldMap *FieldMap, skip map[string]struct{}) {
	fromTypeOfNumField := etyTypeOf.NumField()
	for i := 0; i < fromTypeOfNumField; i++ {

		if etyTypeOf.Field(i).Type.Kind() == reflect.Pointer && reflect.Indirect(etyValueOf.Field(i)).Kind() == reflect.Struct {
			iterateEntityFields(reflect.Indirect(etyValueOf.Field(i)).Type(), reflect.Indirect(etyValueOf.Field(i)), fieldMap, skip)
		} else if etyTypeOf.Field(i).Type.Kind() == reflect.Struct {
			iterateEntityFields(etyTypeOf.Field(i).Type, etyValueOf.Field(i), fieldMap, skip)
		} else {

			fName := etyTypeOf.Field(i).Name
			if _, toSkip := skip[fName]; toSkip {
				continue
			}

			if etyValueOf.Field(i).Type().String() == "map[string]validator.IValidator" {
				continue
			}

			fieldMap.Add(etyTypeOf.Field(i).Name, etyValueOf.Field(i).Type().String())
		}
	}
}
