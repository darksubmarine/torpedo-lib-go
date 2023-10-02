package validator

import (
	"github.com/darksubmarine/torpedo-lib-go/enum"
	"golang.org/x/exp/constraints"
)

type Comparable interface {
	constraints.Ordered | enum.Type | ~bool
}

type Value[T Comparable] struct {
	val       T
	toCompare T
}

func NewValue[T Comparable](v T) *Value[T] {
	return &Value[T]{val: v}
}

func (r *Value[T]) Value(v interface{}) IsValidInterface {
	if val, ok := v.(T); ok {
		r.toCompare = val
	}

	return r
}

func (r *Value[T]) IsValid() bool {
	return r.val == r.toCompare
}
