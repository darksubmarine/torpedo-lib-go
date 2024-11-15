package gin_utils

import (
	"github.com/darksubmarine/torpedo-lib-go/context"
	"github.com/gin-gonic/gin"
)

const (
	KeyDataContext = "TORPEDO_DATA_CONTEXT"
)

// WithDataContext sets a context.DataMap as context application useful to share data within request execution flow.
func WithDataContext() func(*gin.Context) {
	return func(c *gin.Context) {
		// Setting context.DataMap as context application
		c.Set(KeyDataContext, context.NewDataMap())

		c.Next()
	}
}

// SetDataContext sets a key-value pair into the DataContext map.
func SetDataContext(c *gin.Context, key string, val interface{}) {
	dCtx, ok := c.Get(KeyDataContext)
	if !ok {
		dCtx = context.NewDataMap()
	}

	dCtx.(*context.DataMap).Set(key, val)
	c.Set(KeyDataContext, dCtx)
}

// GetDataContext returns the request context.DataMap and a boolean to check if it exists.
// If the data map has not been set a context.NoOpDataMap is returned instead
func GetDataContext(c *gin.Context) (context.IDataMap, bool) {
	if dCtx, ok := c.Get(KeyDataContext); ok {
		cast, ok := dCtx.(*context.DataMap)
		return cast, ok
	}

	return context.NoOpDataMap, false
}
