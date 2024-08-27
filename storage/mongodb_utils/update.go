package mongodb_utils

import (
	"github.com/darksubmarine/torpedo-lib-go/entity"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strings"
)

func ToBSONDocument(dmo interface{}, metadata map[string]*entity.FieldMetadata, skip ...string) bson.D {
	doc := bson.D{}

	dmoTypeOf := reflect.TypeOf(dmo).Elem()
	dmoValueOf := reflect.ValueOf(dmo).Elem()

	skipMap := map[string]struct{}{}
	for _, v := range skip {
		skipMap[v] = struct{}{}
	}

	iterateDMO(dmoTypeOf, dmoValueOf, &doc, skipMap, metadata)

	return doc
}

func iterateDMO(dmoTypeOf reflect.Type, dmoValueOf reflect.Value, doc *bson.D, skip map[string]struct{}, metadata map[string]*entity.FieldMetadata) {
	fromTypeOfNumField := dmoTypeOf.NumField()
	for i := 0; i < fromTypeOfNumField; i++ {
		if dmoTypeOf.Field(i).Type.Kind() == reflect.Struct {
			iterateDMO(dmoTypeOf.Field(i).Type, dmoValueOf.Field(i), doc, skip, metadata)
		} else {
			var fName = dmoTypeOf.Field(i).Name // redding field name from DMO
			var addToUpdateDoc = false          // this field should be added to update doc?
			if !strings.HasSuffix(fName, "_") { // If field name not follows the naming convention "FieldName_" check metadata
				for _, mdata := range metadata {
					if mdata.DmoMongodbName() == fName {
						addToUpdateDoc = true
						break
					}
				}
			} else { // field name with naming should be added
				addToUpdateDoc = true
			}

			if !addToUpdateDoc {
				continue
			}

			if _, toSkip := skip[fName]; toSkip {
				continue
			}

			*doc = append(*doc, bson.E{Key: dmoTypeOf.Field(i).Tag.Get("bson"), Value: dmoValueOf.Field(i).Interface()})
		}
	}
}
