package validator

import (
	"golang.org/x/exp/constraints"
)

type Range[T constraints.Ordered] struct {
	inf       T
	sup       T
	toCompare T
}

func NewRange[T constraints.Ordered](inf, sup T) *Range[T] {
	return &Range[T]{inf: inf, sup: sup}
}

func (r *Range[T]) Value(v interface{}) IsValidInterface {
	if val, ok := v.(T); ok {
		r.toCompare = val
	}
	return r
}

func (r *Range[T]) IsValid() bool {
	return (r.toCompare >= r.inf) && (r.toCompare <= r.sup)
}
