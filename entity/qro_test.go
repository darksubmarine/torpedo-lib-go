package entity_test

import (
	torpedo_lib "github.com/darksubmarine/torpedo-lib-go"
	"github.com/darksubmarine/torpedo-lib-go/entity"
	"github.com/darksubmarine/torpedo-lib-go/tql"
)

type QRO struct {
	/* your custom fields goes here */
	Name_ *string `json:"name,omitempty"`
}

type EntityQRO struct {
	Id_      *string `json:"id,omitempty"`
	Created_ *int64  `json:"created,omitempty"`
	Updated_ *int64  `json:"updated,omitempty"`

	Message_ *string `json:"message,omitempty"`
	PubDate_ *int64  `json:"pubDate,omitempty"`

	QRO
}

func (qro *EntityQRO) FieldValue(field string) interface{} {
	return tql.FieldValue(qro, entity.FieldNameToCode(field))
}

func (qro *EntityQRO) HydrateFromEntity(ety *Entity, fields ...string) error {

	if ety == nil {
		return torpedo_lib.ErrNilEntity
	}

	return entity.To(ety, qro, fields...)
}
