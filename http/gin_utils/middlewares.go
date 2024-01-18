package gin_utils

import (
	"github.com/gin-gonic/gin"
)

// TorpedoMiddleware struct to create middlewares via parameter
type TorpedoMiddleware struct {
	Type MiddlewareTypeEnum
	Fn   gin.HandlerFunc
}

// WithCreateMiddleware attach the given middleware to the Create request
func WithCreateMiddleware(fn gin.HandlerFunc) *TorpedoMiddleware {
	return &TorpedoMiddleware{Type: Create, Fn: fn}
}

// WithReadMiddleware attach the given middleware to the Read request
func WithReadMiddleware(fn gin.HandlerFunc) *TorpedoMiddleware {
	return &TorpedoMiddleware{Type: Read, Fn: fn}
}

// WithUpdateMiddleware attach the given middleware to the Update request
func WithUpdateMiddleware(fn gin.HandlerFunc) *TorpedoMiddleware {
	return &TorpedoMiddleware{Type: Update, Fn: fn}
}

// WithDeleteMiddleware attach the given middleware to the Delete request
func WithDeleteMiddleware(fn gin.HandlerFunc) *TorpedoMiddleware {
	return &TorpedoMiddleware{Type: Delete, Fn: fn}
}

// WithQueryMiddleware attach the given middleware to the Query request
func WithQueryMiddleware(fn gin.HandlerFunc) *TorpedoMiddleware {
	return &TorpedoMiddleware{Type: Query, Fn: fn}
}

// CRUDMiddleware contains all the given CRUD middlewares
type CRUDMiddleware struct {
	Create []gin.HandlerFunc
	Read   []gin.HandlerFunc
	Update []gin.HandlerFunc
	Delete []gin.HandlerFunc
}

// NewCRUDMiddleware constructor of CRUDMiddleware
func NewCRUDMiddleware() *CRUDMiddleware {
	return &CRUDMiddleware{
		Create: []gin.HandlerFunc{},
		Read:   []gin.HandlerFunc{},
		Update: []gin.HandlerFunc{},
		Delete: []gin.HandlerFunc{},
	}
}

// hasMiddleware simple function to check if the given list is empty or not
func (m *CRUDMiddleware) hasMiddleware(list []gin.HandlerFunc) bool {
	if len(list) > 0 {
		return true
	}
	return false
}

// HasCreate check if the Create middleware list is empty or not
func (m *CRUDMiddleware) HasCreate() bool {
	return m.hasMiddleware(m.Create)
}

// HasRead check if the Read middleware list is empty or not
func (m *CRUDMiddleware) HasRead() bool {
	return m.hasMiddleware(m.Read)
}

// HasUpdate check if the Update middleware list is empty or not
func (m *CRUDMiddleware) HasUpdate() bool {
	return m.hasMiddleware(m.Update)
}

// HasDelete check if the Delete middleware list is empty or not
func (m *CRUDMiddleware) HasDelete() bool {
	return m.hasMiddleware(m.Delete)
}

// CRUDQMiddleware extends the CRUDMiddleware to support the Query middleware as well
type CRUDQMiddleware struct {
	CRUDMiddleware
	Query []gin.HandlerFunc
}

// NewCRUDQMiddleware constructor of CRUDQMiddleware
func NewCRUDQMiddleware() *CRUDQMiddleware {
	crudmid := NewCRUDMiddleware()
	return &CRUDQMiddleware{
		CRUDMiddleware: *crudmid,
		Query:          []gin.HandlerFunc{},
	}
}

// HasQuery check if the Query middleware list is empty or not
func (m *CRUDQMiddleware) HasQuery() bool {
	return m.hasMiddleware(m.Query)
}

// ToCRUDMiddlewares given a list of TorpedoMiddleware creates and returns the CRUDMiddleware instance
func ToCRUDMiddlewares(mm ...*TorpedoMiddleware) *CRUDMiddleware {
	var middlewares *CRUDMiddleware
	if len(mm) > 0 {
		middlewares = NewCRUDMiddleware()
		for _, m := range mm {
			switch m.Type {
			case Create:
				middlewares.Create = append(middlewares.Create, m.Fn)
			case Read:
				middlewares.Read = append(middlewares.Read, m.Fn)
			case Update:
				middlewares.Update = append(middlewares.Update, m.Fn)
			case Delete:
				middlewares.Delete = append(middlewares.Delete, m.Fn)
			}
		}
	}

	return middlewares
}

// ToCRUDQMiddlewares given a list of TorpedoMiddleware creates and returns the CRUDQMiddleware instance
func ToCRUDQMiddlewares(mm ...*TorpedoMiddleware) *CRUDQMiddleware {
	var middlewares *CRUDQMiddleware
	if len(mm) > 0 {
		middlewares = NewCRUDQMiddleware()
		for _, m := range mm {
			switch m.Type {
			case Create:
				middlewares.Create = append(middlewares.Create, m.Fn)
			case Read:
				middlewares.Read = append(middlewares.Read, m.Fn)
			case Update:
				middlewares.Update = append(middlewares.Update, m.Fn)
			case Delete:
				middlewares.Delete = append(middlewares.Delete, m.Fn)
			case Query:
				middlewares.Query = append(middlewares.Query, m.Fn)
			}
		}
	}

	return middlewares
}
