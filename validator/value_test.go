package validator

import (
	"github.com/darksubmarine/torpedo-lib-go/enum"
	"strings"
	"testing"
)

func TestValue_Not_IsValid(t *testing.T) {

	if !NewValue("hola mundo").Value("hola mundo").IsValid() {
		t.Error("string comparison error")
	}

	if !NewValue(20).Value(20).IsValid() {
		t.Error("int comparison error")
	}

	if !NewValue(2.32).Value(2.32).IsValid() {
		t.Error("float comparison error")
	}

	if !NewValue(true).Value(true).IsValid() {
		t.Error("boolean comparison error")
	}
}

func TestValue_IsValid(t *testing.T) {

	if NewValue("hola mundo").Value("XYZ").IsValid() {
		t.Error("string comparison error")
	}

	if NewValue(20).Value(10).IsValid() {
		t.Error("int comparison error")
	}

	if NewValue(2.32).Value(1).IsValid() {
		t.Error("float comparison error")
	}

	if NewValue(true).Value(false).IsValid() {
		t.Error("boolean comparison error")
	}
}

/* ENUM */
type CreditCardEnum enum.Type

const (
	Undefined CreditCardEnum = iota
	_
	Visa
	MasterCard
)

func (c CreditCardEnum) ToInt() int {
	return int(c)
}

func (c CreditCardEnum) FromString(s string) CreditCardEnum {
	switch strings.ToLower(s) {
	case "visa":
		return Visa
	case "mastercard":
		return MasterCard
	default:
		return Undefined
	}
}

func (c CreditCardEnum) Value() enum.Type { return enum.Type(c) }

func (c CreditCardEnum) String() string {
	switch c {
	case Undefined:
		return "undefined"
	case Visa:
		return "visa"
	case MasterCard:
		return "mastercard"
	}

	return "undefined"
}

func TestEnum(t *testing.T) {
	var visa = Visa
	var master = MasterCard

	if !NewValue(visa).Value(Visa).IsValid() {
		t.Error("Enum comparison error")
	}

	if NewValue(visa).Value(master).IsValid() {
		t.Error("Enum comparison error")
	}
}
