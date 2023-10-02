package trn_test

import (
	"fmt"
	"github.com/darksubmarine/torpedo-lib-go/trn"
	"testing"
)

func TestTRN_String(t *testing.T) {
	kind := "entity"
	name := "user"
	id := "1234567890QWERTYUIOP"

	rn := trn.New(kind, name, id)

	if rn.Kind() != kind {
		t.Error("invalid kind")
	}

	if rn.Name() != name {
		t.Error("invalid name")
	}

	if rn.ID() != id {
		t.Error("invalid id")
	}

	if rn.String() != fmt.Sprintf("trn::%s::%s::%s", kind, name, id) {
		t.Error("invalid string representation")
	}
	fmt.Println(rn.String())
}

func TestTRN_Equals(t *testing.T) {
	kind := "entity"
	name := "user"
	id1 := "1234567890QWERTYUIOP"
	id2 := "asd123"

	rn1 := trn.New(kind, name, id1)
	rn2 := trn.New(kind, name, id2)

	if rn1.Equals(*rn2) {
		t.Error("wrong equals method")
	}

	if !rn1.Equals(*rn1) {
		t.Error("wrong equals method")
	}
}
