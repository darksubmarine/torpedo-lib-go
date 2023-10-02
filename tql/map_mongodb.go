package tql

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapToMongoDBFilter(q *Query, metadata map[string]string) (bson.D, error) {

	if q.Filter == nil {
		return nil, ErrEmptyFilter
	}

	toFilter := []interface{}{}
	for _, item := range q.Filter.Fields {
		if fType, ok := metadata[item.Field]; !ok {
			return nil, ErrInvalidFieldName
		} else {
			if docs, err := toMongoDoc(item, fType); err != nil {
				return nil, err
			} else {
				toFilter = append(toFilter, docs...)
			}
		}
	}

	if q.Filter.Type() == "any" {
		return bson.D{{"$or", toFilter}}, nil
	}
	return bson.D{{"$and", toFilter}}, nil
}

func toMongoDoc(item FilterItem, fieldType string) ([]interface{}, error) {
	if isSimpleOperator(item.Operator) {
		return toMongoDocSimpleOperator(item, fieldType)
	}

	if isBetweenOperator(item.Operator) {
		return toMongoDocBetweenOperator(item, fieldType)
	}

	if isInListOperator(item.Operator) {
		return toMongoDocInListOperator(item, fieldType)
	}

	if isStringOperator(item.Operator) {
		return toMongoDocStringOperator(item, fieldType)
	}

	return nil, ErrInvalidOperator
}

func toMongoOperator(operator string) (string, error) {
	switch operator {
	case OpNEQ:
		return "$ne", nil
	case OpEQ:
		return "$eq", nil
	case OpGT:
		return "$gt", nil
	case OpGTE:
		return "$gte", nil
	case OpLT:
		return "$lt", nil
	case OpLTE:
		return "$lte", nil
	}

	return "", ErrInvalidOperator
}

func toMongoDocSimpleOperator(item FilterItem, fieldType string) ([]interface{}, error) {
	if fieldType != "string" && fieldType != "int64" && fieldType != "float64" {
		return nil, ErrOperationNotSupported
	}

	var operator string
	if op, err := toMongoOperator(item.Operator); err != nil {
		return nil, err
	} else {
		operator = op
	}

	//if operator != OpEQ && operator != OpNEQ && fieldType == "string" {
	//	return nil, ErrOperationNotSupported
	//}

	return []interface{}{
		bson.D{{field(item.Field), bson.D{{operator, item.Value}}}},
	}, nil
}

func toMongoDocBetweenOperator(item FilterItem, fieldType string) ([]interface{}, error) {
	if fieldType != "int64" {
		return nil, ErrOperationNotSupported
	}

	var left float64
	var right float64
	if val, ok := item.Value.([]interface{}); !ok {
		return nil, ErrInvalidValue
	} else if len(val) != 2 {
		return nil, ErrInvalidValueBTLen
	} else {
		if v, ok := val[0].(float64); !ok {
			return nil, ErrInvalidValue
		} else {
			left = v
		}

		if v, ok := val[1].(float64); !ok {
			return nil, ErrInvalidValue
		} else {
			right = v
		}
	}

	switch item.Operator {
	case OpBTNoLimits:
		return []interface{}{
			bson.D{{"$and",
				[]bson.D{
					{{field(item.Field), bson.D{{"$gt", left}}}},
					{{field(item.Field), bson.D{{"$lt", right}}}},
				},
			}}}, nil
	case OpBTLimits:
		return []interface{}{
			bson.D{{"$and",
				[]bson.D{
					{{field(item.Field), bson.D{{"$gte", left}}}},
					{{field(item.Field), bson.D{{"$lte", right}}}},
				},
			}}}, nil
	case OpBTRightLimit:
		return []interface{}{
			bson.D{{"$and",
				[]bson.D{
					{{field(item.Field), bson.D{{"$gt", left}}}},
					{{field(item.Field), bson.D{{"$lte", right}}}},
				},
			}}}, nil
	case OpBTLeftLimit:
		return []interface{}{
			bson.D{{"$and",
				[]bson.D{
					{{field(item.Field), bson.D{{"$gte", left}}}},
					{{field(item.Field), bson.D{{"$lt", right}}}},
				},
			}}}, nil
	}

	return nil, ErrInvalidOperator
}

func toMongoDocInListOperator(item FilterItem, fieldType string) ([]interface{}, error) {
	if fieldType != "string" && fieldType != "int64" && fieldType != "float64" {
		return nil, ErrOperationNotSupported
	}

	if list, ok := item.Value.([]interface{}); ok {
		inList := make([]interface{}, len(list))

		for i, elem := range list {
			var val interface{}
			var okAssert bool
			switch fieldType {
			case "string":
				val, okAssert = elem.(string)
			case "int64":
				var v float64
				v, okAssert = elem.(float64) // because JSON numbers are float
				val = int64(v)
			case "float64":
				val, okAssert = elem.(float64)
			}

			if !okAssert {
				return nil, ErrInvalidValue
			}

			inList[i] = val
		}

		return []interface{}{
			bson.D{{field(item.Field), bson.D{{"$in", inList}}}},
		}, nil
	}
	return nil, ErrInvalidValue
}

func toMongoDocStringOperator(item FilterItem, fieldType string) ([]interface{}, error) {
	if fieldType != "string" {
		return nil, ErrOperationNotSupported
	}

	if _, ok := item.Value.(string); !ok {
		return nil, ErrInvalidValue
	}

	switch item.Operator {
	case OpPrefix:
		return []interface{}{
			bson.D{{field(item.Field), bson.D{{"$regex", primitive.Regex{Pattern: fmt.Sprintf("^%s", item.Value)}}}}},
		}, nil
	case OpSuffix:
		return []interface{}{
			bson.D{{field(item.Field), bson.D{{"$regex", primitive.Regex{Pattern: fmt.Sprintf("%s$", item.Value)}}}}},
		}, nil
	case OpContains:
		return []interface{}{
			bson.D{{field(item.Field), bson.D{{"$regex", primitive.Regex{Pattern: fmt.Sprintf("%s", item.Value)}}}}},
		}, nil
	}

	return nil, ErrInvalidOperator
}

func field(name string) string {
	if name == "id" {
		return "_id"
	}

	return name
}
