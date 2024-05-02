package context

import (
	"sync"
	"time"
)

// DataMap wrapper of data_struct.KeyValueMap and implements an empty context.Context interface to match with go standards
type DataMap struct {
	sync.Map
}

func NewDataMap() *DataMap {
	return &DataMap{}
}

func (d *DataMap) Set(key string, val interface{}) {
	d.Store(key, val)
}

func (d *DataMap) Get(key string) interface{} {
	if v, ok := d.Load(key); ok {
		return v
	}
	return nil
}

func (d *DataMap) GetOrElse(key string, val interface{}) interface{} {
	if v, ok := d.Load(key); ok {
		return v
	}
	return val
}

func (d *DataMap) GetStringOrElse(key string, val string) string {
	if v, ok := d.Load(key); ok {
		if va, ok := v.(string); ok {
			return va
		}
	}
	return val
}

func (d *DataMap) GetInt64OrElse(key string, val int64) int64 {
	if v, ok := d.Load(key); ok {
		if va, ok := v.(int64); ok {
			return va
		}
	}
	return val
}

func (d *DataMap) GetIntOrElse(key string, val int) int {
	if v, ok := d.Load(key); ok {
		if va, ok := v.(int); ok {
			return va
		}
	}
	return val
}

func (d *DataMap) GetFloat64OrElse(key string, val float64) float64 {
	if v, ok := d.Load(key); ok {
		if va, ok := v.(float64); ok {
			return va
		}
	}
	return val
}

func (d *DataMap) GetBoolOrElse(key string, val bool) bool {
	if v, ok := d.Load(key); ok {
		if va, ok := v.(bool); ok {
			return va
		}
	}
	return val
}

/* CONTEXT INTERFACE */

func (d *DataMap) Deadline() (deadline time.Time, ok bool) {
	return
}

func (d *DataMap) Done() <-chan struct{} {
	return nil
}

func (d *DataMap) Err() error {
	return nil
}

func (d *DataMap) Value(key any) any {
	return nil
}
