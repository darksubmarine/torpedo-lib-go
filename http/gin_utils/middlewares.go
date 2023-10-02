package gin_utils

import "github.com/gin-gonic/gin"

type TorpedoMiddleware struct {
	Type MiddlewareTypeEnum
	Fn   gin.HandlerFunc
}

func WithCreateMiddleware(fn gin.HandlerFunc) *TorpedoMiddleware {
	return &TorpedoMiddleware{Type: Create, Fn: fn}
}

func WithReadMiddleware(fn gin.HandlerFunc) *TorpedoMiddleware {
	return &TorpedoMiddleware{Type: Read, Fn: fn}
}

func WithUpdateMiddleware(fn gin.HandlerFunc) *TorpedoMiddleware {
	return &TorpedoMiddleware{Type: Update, Fn: fn}
}

func WithDeleteMiddleware(fn gin.HandlerFunc) *TorpedoMiddleware {
	return &TorpedoMiddleware{Type: Delete, Fn: fn}
}

type CRUDMiddleware struct {
	Create []gin.HandlerFunc
	Read   []gin.HandlerFunc
	Update []gin.HandlerFunc
	Delete []gin.HandlerFunc
}

func NewCRUDMiddleware() *CRUDMiddleware {
	return &CRUDMiddleware{
		Create: []gin.HandlerFunc{},
		Read:   []gin.HandlerFunc{},
		Update: []gin.HandlerFunc{},
		Delete: []gin.HandlerFunc{},
	}
}

func (m *CRUDMiddleware) hasMiddleware(list []gin.HandlerFunc) bool {
	if len(list) > 0 {
		return true
	}
	return false
}

func (m *CRUDMiddleware) HasCreate() bool {
	return m.hasMiddleware(m.Create)
}

func (m *CRUDMiddleware) HasRead() bool {
	return m.hasMiddleware(m.Read)
}

func (m *CRUDMiddleware) HasUpdate() bool {
	return m.hasMiddleware(m.Update)
}

func (m *CRUDMiddleware) HasDelete() bool {
	return m.hasMiddleware(m.Delete)
}

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
