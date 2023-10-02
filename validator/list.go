package validator

type List[T Comparable] struct {
	values    map[T]struct{}
	toCompare T
}

func NewList[T Comparable](list []T) *List[T] {
	var m = make(map[T]struct{})
	for _, v := range list {
		m[v] = struct{}{}
	}

	return &List[T]{values: m}
}

func (r *List[T]) Value(v interface{}) IsValidInterface {
	if val, ok := v.(T); ok {
		r.toCompare = val
	}
	return r
}

func (r *List[T]) IsValid() bool {
	_, ok := r.values[r.toCompare]
	return ok
}
