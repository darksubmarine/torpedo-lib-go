package entity

import (
	"fmt"
	"reflect"
	"strings"
)

func OptionalFields(ety interface{}, skip ...string) []string {
	etyTypeOf := reflect.TypeOf(ety).Elem()
	etyValueOf := reflect.ValueOf(ety).Elem()

	skipMap := map[string]struct{}{}
	for _, v := range skip {
		skipMap[v] = struct{}{}
	}

	return metadataIterateFields(etyTypeOf, etyValueOf, skipMap)
}

func metadataIterateFields(etyTypeOf reflect.Type, etyValueOf reflect.Value, skip map[string]struct{}) []string {
	optionalFields := []string{}
	fromTypeOfNumField := etyTypeOf.NumField()
	for i := 0; i < fromTypeOfNumField; i++ {

		if etyTypeOf.Field(i).Type.Kind() == reflect.Pointer && reflect.Indirect(etyValueOf.Field(i)).Kind() == reflect.Struct {
			optionalFields = append(optionalFields,
				metadataIterateFields(reflect.Indirect(etyValueOf.Field(i)).Type(), reflect.Indirect(etyValueOf.Field(i)), skip)...)
		} else if etyTypeOf.Field(i).Type.Kind() == reflect.Struct {
			optionalFields = append(optionalFields, metadataIterateFields(etyTypeOf.Field(i).Type, etyValueOf.Field(i), skip)...)
		} else {

			fName := etyTypeOf.Field(i).Name
			if _, toSkip := skip[fName]; toSkip {
				continue
			}

			fm := readFieldMetadata(etyTypeOf.Field(i))
			if fm.optional {
				optionalFields = append(optionalFields, etyTypeOf.Field(i).Name)
				continue
			}
		}
	}

	return optionalFields
}

func FieldsMetadata(ety interface{}) map[string]*FieldMetadata {
	etyTypeOf := reflect.TypeOf(ety).Elem()
	etyValueOf := reflect.ValueOf(ety).Elem()

	return fetchFieldsMetadata(etyTypeOf, etyValueOf)
}

func fetchFieldsMetadata(etyTypeOf reflect.Type, etyValueOf reflect.Value) map[string]*FieldMetadata {
	metadataFields := map[string]*FieldMetadata{}
	fromTypeOfNumField := etyTypeOf.NumField()
	for i := 0; i < fromTypeOfNumField; i++ {

		var _typeOf reflect.Type
		var _valueOf reflect.Value

		if etyTypeOf.Field(i).Type.Kind() == reflect.Pointer && reflect.Indirect(etyValueOf.Field(i)).Kind() == reflect.Struct {
			_typeOf = reflect.Indirect(etyValueOf.Field(i)).Type()
			_valueOf = reflect.Indirect(etyValueOf.Field(i))
		} else if etyTypeOf.Field(i).Type.Kind() == reflect.Struct {
			_typeOf = etyTypeOf.Field(i).Type
			_valueOf = etyValueOf.Field(i)
		}

		if _typeOf != nil { // if we have type it is a Struct ... so, recursive call!
			auxMap := fetchFieldsMetadata(_typeOf, _valueOf)
			for k, v := range auxMap {
				metadataFields[k] = v
			}
		} else {
			fName := etyTypeOf.Field(i).Name
			metadataFields[fName] = readFieldMetadata(etyTypeOf.Field(i))
		}

		//if etyTypeOf.Field(i).Type.Kind() == reflect.Pointer && reflect.Indirect(etyValueOf.Field(i)).Kind() == reflect.Struct {
		//	auxMap := fetchFieldsMetadata(reflect.Indirect(etyValueOf.Field(i)).Type(), reflect.Indirect(etyValueOf.Field(i)))
		//	for k, v := range auxMap {
		//		metadataFields[k] = v
		//	}
		//} else if etyTypeOf.Field(i).Type.Kind() == reflect.Struct {
		//	auxMap := fetchFieldsMetadata(etyTypeOf.Field(i).Type, etyValueOf.Field(i))
		//	for k, v := range auxMap {
		//		metadataFields[k] = v
		//	}
		//} else {
		//
		//	fName := etyTypeOf.Field(i).Name
		//	metadataFields[fName] = readFieldMetadata(etyTypeOf.Field(i))
		//
		//}
	}

	return metadataFields
}

type FieldMetadata struct {
	optional  bool
	encrypted bool
	getter    string
	setter    string
	dmo       struct {
		memory  string
		redis   string
		mongodb string
		sql     string
	}
	dto struct {
		http string
	}
	qro string
}

func (f *FieldMetadata) IsOptional() bool       { return f.optional }
func (f *FieldMetadata) IsEncrypted() bool      { return f.encrypted }
func (f *FieldMetadata) Getter() string         { return f.getter }
func (f *FieldMetadata) Setter() string         { return f.setter }
func (f *FieldMetadata) DmoMemoryName() string  { return f.dmo.memory }
func (f *FieldMetadata) DmoRedisName() string   { return f.dmo.redis }
func (f *FieldMetadata) DmoMongodbName() string { return f.dmo.mongodb }
func (f *FieldMetadata) DmoSqlName() string     { return f.dmo.sql }
func (f *FieldMetadata) DtoHttpName() string    { return f.dto.http }
func (f *FieldMetadata) QroName() string        { return f.qro }

func readFieldMetadata(field reflect.StructField) *FieldMetadata {

	meta := &FieldMetadata{optional: false, encrypted: false}

	if tVal, ok := field.Tag.Lookup(tagField); ok {
		if strings.Contains(tVal, "optional") {
			meta.optional = true
		}

		if strings.Contains(tVal, "encrypted") {
			meta.encrypted = true
		}
	}

	if tVal, ok := field.Tag.Lookup(tagGetter); ok {
		meta.getter = tVal
	} else {
		methodName, _ := strings.CutSuffix(FieldNameToCode(field.Name), "_")
		meta.getter = methodName
	}

	if tVal, ok := field.Tag.Lookup(tagSetter); ok {
		meta.setter = tVal
	} else {
		methodName, _ := strings.CutSuffix(FieldNameToCode(field.Name), "_")
		meta.setter = fmt.Sprintf("Set%s", methodName)
	}

	if tVal, ok := field.Tag.Lookup(tagQRO); ok {
		meta.qro = tVal
	} else {
		meta.qro = FieldNameToCode(field.Name)
	}

	if tVal, ok := field.Tag.Lookup(tagDMO); ok {
		parts := strings.Split(tVal, ",")
		for _, output := range parts {
			dmo := strings.Split(output, "=")
			if len(dmo) == 2 {
				switch dmo[0] {
				case "memory":
					meta.dmo.memory = dmo[1]
				case "mongodb":
					meta.dmo.mongodb = dmo[1]
				case "redis":
					meta.dmo.redis = dmo[1]
				case "sql":
					meta.dmo.sql = dmo[1]
				}
			}
		}
	}

	if tVal, ok := field.Tag.Lookup(tagDTO); ok {
		parts := strings.Split(tVal, ",")
		for _, input := range parts {
			dto := strings.Split(input, "=")
			if len(dto) == 2 {
				switch dto[0] {
				case "http":
					meta.dto.http = dto[1]
				}
			}
		}
	}

	return meta
}
