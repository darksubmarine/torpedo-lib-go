package entity_test

import "github.com/darksubmarine/torpedo-lib-go/validator"

type entityBase struct {
	// required
	id      string
	created int64
	updated int64

	// schema fields
	_string  string
	_int     int
	_boolean bool
	_slice   []int

	validators map[string]validator.IValidator
}

func (e *entityBase) SetId(id string) { e.id = id }
func (e *entityBase) Id() string      { return e.id }

func (e *entityBase) SetCreated(created int64) { e.created = created }
func (e *entityBase) Created() int64           { return e.created }

func (e *entityBase) SetUpdated(updated int64) { e.updated = updated }
func (e *entityBase) Updated() int64           { return e.updated }

func (e *entityBase) SetString(s string) { e._string = s }
func (e *entityBase) String() string     { return e._string }

func (e *entityBase) SetInt(i int) { e._int = i }
func (e *entityBase) Int() int     { return e._int }

func (e *entityBase) SetBoolean(b bool) { e._boolean = b }
func (e *entityBase) Boolean() bool     { return e._boolean }

func (e *entityBase) SetSlice(s []int) { e._slice = s }
func (e *entityBase) Slice() []int     { return e._slice }

type Entity struct {
	*entityBase

	// custom user fields (out of schema)
	name string
}

func (e *Entity) SetName(n string) { e.name = n }
func (e *Entity) Name() string     { return e.name }

func NewEntity() *Entity {
	//return &Entity{entityBase: &entityBase{_slice: []int{}}}
	return &Entity{entityBase: &entityBase{}}
}
