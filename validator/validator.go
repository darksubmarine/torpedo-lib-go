package validator

type IsValidInterface interface {
	IsValid() bool
}

type ValueInterface interface {
	Value(val interface{}) IsValidInterface
}

type IValidator interface {
	IsValidInterface
	ValueInterface
}
