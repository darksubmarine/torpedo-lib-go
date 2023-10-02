package enum

type Type = uint8

type Enum interface {
	Value() Type
	String() string
}
