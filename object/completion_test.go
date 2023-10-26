package object_test

import (
	"github.com/darksubmarine/torpedo-lib-go/object"
	"github.com/darksubmarine/torpedo-lib-go/ptr"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type ObjectBase struct {
	Id_      *string `json:"id"`
	Created_ *int64  `json:"created"`
	Updated_ *int64  `json:"updated"`
	ObjectPartial
}

type ObjectPartial struct {
	String_  *string `json:"string,omitempty" read_method:"GetString"`
	Int_     *int    `json:"number,omitempty"`
	Boolean_ *bool   `json:"boolean,omitempty"`
	Slice_   []int   `json:"slice"`
}

type ObjectFinal struct {
	ObjectBase
	Name_ *string `json:"name,omitempty"`
}

func TestObject_IsCompleteFalse(t *testing.T) {
	obj := ObjectFinal{}
	ok, err := object.IsComplete(&obj)
	assert.False(t, ok)
	assert.Nil(t, err)
}

func TestObject_IsCompleteTrue(t *testing.T) {
	obj := ObjectFinal{
		ObjectBase: ObjectBase{
			Id_:      ptr.String("some id"),
			Created_: ptr.Int64(time.Now().UnixMilli()),
			Updated_: ptr.Int64(time.Now().UnixMilli()),
			ObjectPartial: ObjectPartial{
				Boolean_: ptr.Bool(true),
				Int_:     ptr.Int(123),
				String_:  ptr.String("some string"),
				Slice_:   []int{},
			},
		},
		Name_: ptr.String("some value"),
	}
	ok, err := object.IsComplete(&obj)
	assert.True(t, ok)
	assert.Nil(t, err)
}

type ObjectFinalBadDefinition struct {
	ObjectBase
	Name_ string `json:"name,omitempty"`
}

func TestObject_IsCompleteError(t *testing.T) {
	obj := ObjectFinalBadDefinition{
		ObjectBase: ObjectBase{
			Id_:      ptr.String("some id"),
			Created_: ptr.Int64(time.Now().UnixMilli()),
			Updated_: ptr.Int64(time.Now().UnixMilli()),
			ObjectPartial: ObjectPartial{
				Boolean_: ptr.Bool(true),
				Int_:     ptr.Int(123),
				String_:  ptr.String("some string"),
				Slice_:   []int{},
			},
		},
	}
	ok, err := object.IsComplete(&obj)
	assert.False(t, ok)
	assert.Error(t, err)
}
