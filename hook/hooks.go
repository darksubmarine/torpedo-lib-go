package hook

import (
	"github.com/darksubmarine/torpedo-lib-go/enum"
	"strings"
)

type HookEnum enum.Type

const (
	Undefined HookEnum = iota
	_
	BeforeCreate
	AfterCreate
	BeforeRead
	AfterRead
	BeforeUpdate
	AfterUpdate
	BeforeDelete
	AfterDelete
)

func NewHookEnumFromString(s string) HookEnum {
	switch strings.ToLower(s) {
	case "beforecreate":
		return BeforeCreate
	case "aftercreate":
		return AfterCreate
	case "beforeread":
		return BeforeRead
	case "afterread":
		return AfterRead
	case "beforeupdate":
		return BeforeUpdate
	case "afterupdate":
		return AfterUpdate
	case "beforedelete":
		return BeforeDelete
	case "afterdelete":
		return AfterDelete
	default:
		return Undefined
	}
}

func (e HookEnum) ToInt() int {
	return int(e)
}

func (e HookEnum) Value() enum.Type { return enum.Type(e) }

func (e HookEnum) String() string {
	switch e {
	case Undefined:
		return "undefined"
	case BeforeCreate:
		return "BeforeCreate"
	case AfterCreate:
		return "AfterCreate"
	case BeforeRead:
		return "BeforeRead"
	case AfterRead:
		return "AfterRead"
	case BeforeUpdate:
		return "BeforeUpdate"
	case AfterUpdate:
		return "AfterUpdate"
	case BeforeDelete:
		return "BeforeDelete"
	case AfterDelete:
		return "AfterDelete"
	}

	return "undefined"
}
