package gin_utils

import (
	"github.com/darksubmarine/torpedo-lib-go/enum"
	"strings"
)

type MiddlewareTypeEnum enum.Type

const (
	Undefined MiddlewareTypeEnum = iota
	_
	Create
	Read
	Update
	Delete
	Query
)

func NewMiddlewareTypeEnumFromString(s string) MiddlewareTypeEnum {
	switch strings.ToLower(s) {
	case "create":
		return Create
	case "read":
		return Read
	case "update":
		return Update
	case "delete":
		return Delete
	case "query":
		return Query
	default:
		return Undefined
	}
}

func (c MiddlewareTypeEnum) ToInt() int {
	return int(c)
}

func (c MiddlewareTypeEnum) Value() enum.Type { return enum.Type(c) }

func (c MiddlewareTypeEnum) String() string {
	switch c {
	case Undefined:
		return "undefined"
	case Create:
		return "create"
	case Read:
		return "read"
	case Update:
		return "update"
	case Delete:
		return "delete"
	case Query:
		return "query"
	}

	return "undefined"
}
