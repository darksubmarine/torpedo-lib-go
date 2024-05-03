package context

import (
	"sync"
)

// NoOpDataMap is a NoopDataMap to reference it avoiding memory allocation creating new instances of this one.
var NoOpDataMap = NewNoopDataMap()

// IDataMap interface to defines DataMap methods.
type IDataMap interface {
	Set(key string, val interface{})
	Get(key string) interface{}
	GetOrElse(key string, val interface{}) interface{}
	GetStringOrElse(key string, val string) string
	GetInt64OrElse(key string, val int64) int64
	GetIntOrElse(key string, val int) int
	GetFloat64OrElse(key string, val float64) float64
	GetBoolOrElse(key string, val bool) bool
}

// DataMap wrapper of data_struct.KeyValueMap and implements an empty context.Context interface to match with go standards.
type DataMap struct {
	sync.Map
	NoopGoContext
}

// NewDataMap DataMap constructor.
func NewDataMap() *DataMap {
	return &DataMap{}
}

// Set sets a given value under the given key
func (d *DataMap) Set(key string, val interface{}) {
	d.Store(key, val)
}

// Get returns the value saved under the given key. Returns nil if key not found.
func (d *DataMap) Get(key string) interface{} {
	if v, ok := d.Load(key); ok {
		return v
	}
	return nil
}

// GetOrElse returns the value saved under the given key. Returns val if key not found.
func (d *DataMap) GetOrElse(key string, val interface{}) interface{} {
	if v, ok := d.Load(key); ok {
		return v
	}
	return val
}

// GetStringOrElse returns the value, string asserted, saved under the given key. Returns val if key not found or data type mismatch.
func (d *DataMap) GetStringOrElse(key string, val string) string {
	if v, ok := d.Load(key); ok {
		if va, ok := v.(string); ok {
			return va
		}
	}
	return val
}

// GetInt64OrElse returns the value, int64 asserted, saved under the given key. Returns val if key not found or data type mismatch.
func (d *DataMap) GetInt64OrElse(key string, val int64) int64 {
	if v, ok := d.Load(key); ok {
		if va, ok := v.(int64); ok {
			return va
		}
	}
	return val
}

// GetIntOrElse returns the value, int asserted, saved under the given key. Returns val if key not found or data type mismatch.
func (d *DataMap) GetIntOrElse(key string, val int) int {
	if v, ok := d.Load(key); ok {
		if va, ok := v.(int); ok {
			return va
		}
	}
	return val
}

// GetFloat64OrElse returns the value, float64 asserted, saved under the given key. Returns val if key not found or data type mismatch.
func (d *DataMap) GetFloat64OrElse(key string, val float64) float64 {
	if v, ok := d.Load(key); ok {
		if va, ok := v.(float64); ok {
			return va
		}
	}
	return val
}

// GetBoolOrElse returns the value, bool asserted, saved under the given key. Returns val if key not found or data type mismatch.
func (d *DataMap) GetBoolOrElse(key string, val bool) bool {
	if v, ok := d.Load(key); ok {
		if va, ok := v.(bool); ok {
			return va
		}
	}
	return val
}

// Value returns the value saved under the given key. Returns nil if key not found or is invalid data type.
// This method overrides the context.Context method
func (d *DataMap) Value(key any) any {
	if k, ok := key.(string); ok {
		return d.Get(k)
	}

	return nil
}

// NoopDataMap no operational Data Map
type NoopDataMap struct{ NoopGoContext }

// NewNoopDataMap constructor of NoopDataMap.
func NewNoopDataMap() *NoopDataMap {
	return new(NoopDataMap)
}

func (d *NoopDataMap) Set(key string, val interface{})                   {}
func (d *NoopDataMap) Get(key string) interface{}                        { return nil }
func (d *NoopDataMap) GetOrElse(key string, val interface{}) interface{} { return nil }
func (d *NoopDataMap) GetStringOrElse(key string, val string) string     { return val }
func (d *NoopDataMap) GetInt64OrElse(key string, val int64) int64        { return val }
func (d *NoopDataMap) GetIntOrElse(key string, val int) int              { return val }
func (d *NoopDataMap) GetFloat64OrElse(key string, val float64) float64  { return val }
func (d *NoopDataMap) GetBoolOrElse(key string, val bool) bool           { return val }
