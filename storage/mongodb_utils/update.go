package mongodb_utils

import (
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strings"
)

func ToBSONDocument(dmo interface{}, skip ...string) bson.D {
	doc := bson.D{}

	dmoTypeOf := reflect.TypeOf(dmo).Elem()
	dmoValueOf := reflect.ValueOf(dmo).Elem()

	skipMap := map[string]struct{}{}
	for _, v := range skip {
		skipMap[v] = struct{}{}
	}

	iterateDMO(dmoTypeOf, dmoValueOf, &doc, skipMap)

	return doc
}

func iterateDMO(dmoTypeOf reflect.Type, dmoValueOf reflect.Value, doc *bson.D, skip map[string]struct{}) {
	fromTypeOfNumField := dmoTypeOf.NumField()
	for i := 0; i < fromTypeOfNumField; i++ {
		if dmoTypeOf.Field(i).Type.Kind() == reflect.Struct {
			iterateDMO(dmoTypeOf.Field(i).Type, dmoValueOf.Field(i), doc, skip)
		} else if name := dmoTypeOf.Field(i).Name; strings.HasSuffix(name, "_") {

			fName := dmoTypeOf.Field(i).Name
			if _, toSkip := skip[fName]; toSkip {
				continue
			}

			//fmt.Println(dmoTypeOf.Field(i).Name, dmoTypeOf.Field(i).Tag.Get("bson"), dmoValueOf.Field(i).Type().String(), dmoValueOf.Field(i).Interface())
			*doc = append(*doc, bson.E{Key: dmoTypeOf.Field(i).Tag.Get("bson"), Value: dmoValueOf.Field(i).Interface()})
		}
	}
}
