package context_test

import (
	"github.com/darksubmarine/torpedo-lib-go/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataMap_SyncMap(t *testing.T) {
	dc := context.NewDataMap()

	dc.Set("k1", 123123)
	assert.Equal(t, 123123, dc.GetIntOrElse("k1", 87))
}

func TestDataMap_Set_Get(t *testing.T) {
	m := context.NewDataMap()
	m.Set("key1", 89)
	m.Set("key2", "hola")

	v1 := m.Get("key1")
	assert.IsType(t, 0, v1)
	assert.EqualValues(t, 89, v1)

	v2 := m.Get("key2")
	assert.IsType(t, "", v2)
	assert.EqualValues(t, "hola", v2)
}

func TestDataMap_GetOrElse(t *testing.T) {
	m := context.NewDataMap()

	type defaultStruct struct {
		value int
	}

	v := m.GetOrElse("no-valid", defaultStruct{value: 123})
	assert.IsType(t, defaultStruct{}, v)
	assert.EqualValues(t, v.(defaultStruct).value, 123)
}

func TestDataMap_GetBoolOrElse(t *testing.T) {
	m := context.NewDataMap()
	m.Set("bool", true)

	v := m.GetBoolOrElse("bool", false)
	assert.True(t, v)

	v1 := m.GetBoolOrElse("invalid", true)
	assert.True(t, v1)
}

func TestDataMap_GetFloat64OrElse(t *testing.T) {
	m := context.NewDataMap()
	m.Set("key", 12.45)

	v := m.GetFloat64OrElse("key", 456.98)
	assert.EqualValues(t, 12.45, v)

	v1 := m.GetFloat64OrElse("invalid", 23.5)
	assert.EqualValues(t, 23.5, v1)
}

func TestDataMap_GetInt64OrElse(t *testing.T) {
	m := context.NewDataMap()
	m.Set("key", int64(1245))

	v := m.GetInt64OrElse("key", 45698)
	assert.EqualValues(t, int64(1245), v)

	v1 := m.GetInt64OrElse("invalid", 235)
	assert.EqualValues(t, 235, v1)
}

func TestDataMap_GetIntOrElse(t *testing.T) {
	m := context.NewDataMap()
	m.Set("key", 1245)

	v := m.GetIntOrElse("key", 45698)
	assert.EqualValues(t, 1245, v)

	v1 := m.GetIntOrElse("invalid", 235)
	assert.EqualValues(t, 235, v1)
}

func TestDataMap_GetStringOrElse(t *testing.T) {
	m := context.NewDataMap()
	m.Set("key", "hola")

	v := m.GetStringOrElse("key", "just in case")
	assert.EqualValues(t, "hola", v)

	v1 := m.GetStringOrElse("invalid", "default string")
	assert.EqualValues(t, "default string", v1)
}

func TestDataMap_Value(t *testing.T) {
	m := context.NewDataMap()
	m.Set("keyInt", 89)
	m.Set("keyString", "torpedo")
	m.Set("keyBool", true)

	v1 := m.Value("keyInt")
	assert.IsType(t, 0, v1)
	assert.EqualValues(t, 89, v1)

	v2 := m.Value("keyString")
	assert.IsType(t, "", v2)
	assert.EqualValues(t, "torpedo", v2)

	v3 := m.Value("keyBool")
	assert.IsType(t, false, v3)
	assert.EqualValues(t, true, v3)
}

func TestNewNoopDataMap(t *testing.T) {
	m := context.NewNoopDataMap()
	assert.Nil(t, m.Value("key"))
}
