package enum_test

import (
	"github.com/darksubmarine/torpedo-lib-go/enum"
	"strings"
	"testing"
)

/* ENUM */
type CreditCardEnum enum.Type

const (
	Undefined CreditCardEnum = iota
	_
	Visa
	MasterCard
)

func NewCreditCardEnumFromString(s string) CreditCardEnum {
	switch strings.ToLower(s) {
	case "visa":
		return Visa
	case "mastercard":
		return MasterCard
	default:
		return Undefined
	}
}

func (c CreditCardEnum) ToInt() int {
	return int(c)
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

func TestCreditCardEnum(t *testing.T) {

	var visa = Visa
	var master = MasterCard

	if NewCreditCardEnumFromString("visa") != visa {
		t.Error("invalid enum from constructor")
	}

	if visa.ToInt() != 2 {
		t.Error("invalid int cardinal for VISA enumerator. Expected 2 Current:", visa.ToInt())
	}

	if master.ToInt() != 3 {
		t.Error("invalid int cardinal for MASTER enumerator. Expected 3 Current:", master.ToInt())
	}
}
