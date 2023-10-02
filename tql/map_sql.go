package tql

import (
	"fmt"
	"strings"
)

func MapToSQLFilter(q *Query, metadata map[string]string) (string, map[string]interface{}, error) {

	if q.Filter == nil {
		return "", nil, ErrEmptyFilter
	}

	toValues := map[string]interface{}{}
	toFilter := []string{}
	for itid, item := range q.Filter.Fields {
		if fType, ok := metadata[item.Field]; !ok {
			return "", nil, ErrInvalidFieldName
		} else {
			if qry, values, err := toSQLQuery(item, fType, itid); err != nil {
				return "", nil, err
			} else {
				toFilter = append(toFilter, qry)
				for k, v := range values {
					toValues[k] = v
				}
			}
		}
	}

	if q.Filter.Type() == "any" {
		return strings.Join(toFilter, " OR "), toValues, nil
	}

	return strings.Join(toFilter, " AND "), toValues, nil
}

func toSQLQuery(item FilterItem, fieldType string, itid int) (string, map[string]interface{}, error) {
	if isSimpleOperator(item.Operator) {
		return toSQLSimpleOperator(item, fieldType, itid)
	}

	if isBetweenOperator(item.Operator) {
		return toSQLBetweenOperator(item, fieldType, itid)
	}

	if isInListOperator(item.Operator) {
		str, err := toSQLInListOperator(item, fieldType)
		return str, nil, err
	}

	if isStringOperator(item.Operator) {
		str, err := toSQLStringOperator(item, fieldType)
		return str, nil, err
	}

	return "", nil, ErrInvalidOperator
}

func toSQLOperator(operator string) (string, error) {
	switch operator {
	case OpNEQ:
		return "!=", nil
	case OpEQ:
		return "==", nil
	case OpGT:
		return ">", nil
	case OpGTE:
		return ">=", nil
	case OpLT:
		return "<", nil
	case OpLTE:
		return "<=", nil
	}

	return "", ErrInvalidOperator
}

func toSQLSimpleOperator(item FilterItem, fieldType string, itid int) (string, map[string]interface{}, error) {
	if fieldType != "string" && fieldType != "int64" && fieldType != "float64" {
		return "", nil, ErrOperationNotSupported
	}

	var operator string
	if op, err := toSQLOperator(item.Operator); err != nil {
		return "", nil, err
	} else {
		operator = op
	}

	//if operator != OpEQ && operator != OpNEQ && fieldType == "string" {
	//	return nil, ErrOperationNotSupported
	//}

	namedField := fmt.Sprintf("%s_%d", item.Field, itid)
	namedQuery := fmt.Sprintf("%s %s :%s", item.Field, operator, namedField)
	values := map[string]interface{}{}

	switch fieldType {
	case "string":
		values[namedField] = fmt.Sprintf("%s", item.Value)
		return namedQuery, values, nil
	case "int64":
		v, okAssert := item.Value.(float64) // because JSON numbers are float
		if okAssert {
			values[namedField] = int64(v)
			return namedQuery, values, nil
		}
	case "float64":
		v, okAssert := item.Value.(float64)
		if okAssert {
			values[namedField] = v
			return namedQuery, values, nil
		}
	}

	return "", nil, ErrOperationNotSupported
}

func toSQLBetweenOperator(item FilterItem, fieldType string, itid int) (string, map[string]interface{}, error) {
	if fieldType != "int64" && fieldType != "float64" {
		return "", nil, ErrOperationNotSupported
	}

	var left float64
	var right float64
	if val, ok := item.Value.([]interface{}); !ok {
		return "", nil, ErrInvalidValue
	} else if len(val) != 2 {
		return "", nil, ErrInvalidValueBTLen
	} else {
		if v, ok := val[0].(float64); !ok {
			return "", nil, ErrInvalidValue
		} else {
			left = v
		}

		if v, ok := val[1].(float64); !ok {
			return "", nil, ErrInvalidValue
		} else {
			right = v
		}
	}

	namedFieldL := fmt.Sprintf("%s_%d_L", item.Field, itid)
	namedFieldR := fmt.Sprintf("%s_%d_R", item.Field, itid)
	var namedQuery string
	values := map[string]interface{}{}

	if fieldType == "int64" {
		values[namedFieldL] = int64(left)
		values[namedFieldR] = int64(right)
	} else {
		values[namedFieldL] = left
		values[namedFieldR] = right
	}

	switch item.Operator {
	case OpBTNoLimits:
		namedQuery = fmt.Sprintf("(%s > :%s AND %s < :%s)", item.Field, namedFieldL, item.Field, namedFieldR)
		return namedQuery, values, nil
	case OpBTLimits:
		namedQuery = fmt.Sprintf("(%s >= :%s AND %s <= :%s)", item.Field, namedFieldL, item.Field, namedFieldR)
		return namedQuery, values, nil
	case OpBTRightLimit:
		namedQuery = fmt.Sprintf("(%s > :%s AND %s <= :%s)", item.Field, namedFieldL, item.Field, namedFieldR)
		return namedQuery, values, nil
	case OpBTLeftLimit:
		namedQuery = fmt.Sprintf("(%s >= :%s AND %s < :%s)", item.Field, namedFieldL, item.Field, namedFieldR)
		return namedQuery, values, nil
	}

	return "", nil, ErrInvalidOperator
}

func toSQLInListOperator(item FilterItem, fieldType string) (string, error) {
	if fieldType != "string" && fieldType != "int64" && fieldType != "float64" {
		return "", ErrOperationNotSupported
	}

	if list, ok := item.Value.([]interface{}); ok {
		inList := make([]string, len(list))

		for i, elem := range list {
			var val string
			var okAssert bool
			switch fieldType {
			case "string":
				val, okAssert = elem.(string)
				val = sanitizeQryStr(val)
			case "int64":
				var v float64
				v, okAssert = elem.(float64) // because JSON numbers are float
				if okAssert {
					val = fmt.Sprintf("%d", int64(v))
				}
			case "float64":
				var v float64
				v, okAssert = elem.(float64)
				if okAssert {
					val = fmt.Sprintf("%f", v)
				}
			}

			if !okAssert {
				return "", ErrInvalidValue
			}

			inList[i] = val
		}

		return fmt.Sprintf("%s IN (%s)", item.Field, strings.Join(inList, ",")), nil
	}
	return "", ErrInvalidValue
}

func toSQLStringOperator(item FilterItem, fieldType string) (string, error) {
	if fieldType != "string" {
		return "", ErrOperationNotSupported
	}

	if val, ok := item.Value.(string); !ok {
		return "", ErrInvalidValue
	} else {
		switch item.Operator {
		case OpPrefix:
			return fmt.Sprintf("%s LIKE '%s%%'", item.Field, sanitizeQryLikeStr(val)), nil
		case OpSuffix:
			return fmt.Sprintf("%s LIKE '%%%s'", item.Field, sanitizeQryLikeStr(val)), nil
		case OpContains:
			return fmt.Sprintf("%s LIKE '%%%s%%'", item.Field, sanitizeQryLikeStr(val)), nil
		}
	}
	return "", ErrInvalidOperator
}
