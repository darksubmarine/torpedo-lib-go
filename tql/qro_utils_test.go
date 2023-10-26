package tql_test

import (
	"github.com/darksubmarine/torpedo-lib-go/ptr"
	"github.com/darksubmarine/torpedo-lib-go/tql"
	"github.com/stretchr/testify/assert"
	"testing"
)

type QRO struct {
	Custom_ *string `json:"custom"`
}

type EntityQRO struct {
	Id_      *string `json:"id,omitempty"`
	Created_ *int64  `json:"created,omitempty"`
	Updated_ *int64  `json:"updated,omitempty"`

	Message_ *string `json:"message,omitempty"`

	PubDate_ *int64 `json:"pubDate,omitempty"`

	PostId_ *string `json:"postId,omitempty"`

	QRO
}

func TestFieldValue(t *testing.T) {
	qro := &EntityQRO{
		Id_:      ptr.String("qwerty"),
		Created_: ptr.Int64(123),
		Updated_: ptr.Int64(456),
		Message_: ptr.String("in a bottle"),
		PubDate_: ptr.Int64(789),
		PostId_:  ptr.String("asdfgh"),
		QRO: QRO{
			Custom_: ptr.String("some custom value!"),
		},
	}

	val := tql.FieldValue(qro, "Custom_")
	switch v := val.(type) {
	case string:
		assert.EqualValues(t, "some custom value!", v)
	default:
		assert.Fail(t, "invalid data type")
	}

}
