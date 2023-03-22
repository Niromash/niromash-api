package api

import (
	"github.com/gin-gonic/gin"
)

type HttpMethod int

const (
	MethodGet HttpMethod = iota
	MethodPost
	MethodDelete
	MethodPut
	MethodPatch
	MethodAny
)

type HttpService interface {
	ServiceInitializer
	ServiceStarter
	Engine() *gin.Engine
	RegisterRoutes(...*Route)
	RegisterControllers(...*Controller)
}

type MiddlewareHandler func(c *gin.Context, user User, mainService MainService) (next bool)
type Handler func(c *gin.Context, user User, mainService MainService)

type Route struct {
	Path                   string
	Method                 HttpMethod
	IsAuthenticated        bool
	AuthenticateMiddleware func(c *gin.Context, user User, mainService MainService)
	Middlewares            []MiddlewareHandler
	Handler                Handler
}

type Controller struct {
	Path   string
	Routes []*Route
}

func (m HttpMethod) String() string {
	switch m {
	case MethodGet:
		return "GET"
	case MethodPost:
		return "POST"
	case MethodDelete:
		return "DELETE"
	case MethodPut:
		return "PUT"
	case MethodPatch:
		return "PATCH"
	case MethodAny:
		return "ANY"
	}
	return ""
}

func (m HttpMethod) ToFunctionFromEngine(app *gin.Engine) func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return []func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes{
		app.GET, app.POST, app.DELETE, app.PUT, app.PATCH, app.Any,
	}[m]
}
