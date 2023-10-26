package entity_test

import (
	"github.com/darksubmarine/torpedo-lib-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToFieldMap(t *testing.T) {
	ety := NewEntity()

	fm := entity.ToFieldMap(ety)
	/*
		id      string
			created int64
			updated int64

			// schema fields
			_string  string
			_int     int
			_boolean bool
			_slice   []int
		name string
	*/
	var kind string
	assert.True(t, fm.HasField("id"))
	kind, _ = fm.FieldType("id")
	assert.EqualValues(t, "string", kind)

	assert.True(t, fm.HasField("created"))
	kind, _ = fm.FieldType("created")
	assert.EqualValues(t, "int64", kind)

	assert.True(t, fm.HasField("updated"))
	kind, _ = fm.FieldType("updated")
	assert.EqualValues(t, "int64", kind)

	assert.True(t, fm.HasField("_string"))
	kind, _ = fm.FieldType("_string")
	assert.EqualValues(t, "string", kind)

	assert.True(t, fm.HasField("_int"))
	kind, _ = fm.FieldType("_int")
	assert.EqualValues(t, "int", kind)

	assert.True(t, fm.HasField("_boolean"))
	kind, _ = fm.FieldType("_boolean")
	assert.EqualValues(t, "bool", kind)

	assert.True(t, fm.HasField("_slice"))
	kind, _ = fm.FieldType("_slice")
	assert.EqualValues(t, "[]int", kind)

	assert.True(t, fm.HasField("name"))
	kind, _ = fm.FieldType("name")
	assert.EqualValues(t, "string", kind)

}
