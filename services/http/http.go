package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"niromash-api/api"
	"niromash-api/services/http/controllers"
	"niromash-api/services/http/middlewares"
	"niromash-api/services/http/routes"
	"strings"
	"time"
)

var _ api.HttpService = (*HttpService)(nil)

type HttpService struct {
	engine  *gin.Engine
	service api.MainService
}

func NewHttpService() *HttpService {
	return &HttpService{}
}

func (h *HttpService) Init(service api.MainService) error {
	h.service = service
	h.engine = gin.New()
	h.engine.Use(h.logger(), gin.Recovery())

	h.engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	h.RegisterControllers(controllers.ProjectsController(), controllers.StatsController(), controllers.AuthController(),
		controllers.MessageController())
	h.RegisterRoutes(routes.HealthRoute())

	return nil
}

func (h *HttpService) Start() error {
	return h.engine.Run(fmt.Sprintf(":3000"))
}

func (h *HttpService) Close() error {
	return nil
}

func (h *HttpService) Engine() *gin.Engine {
	return h.engine
}

func (h *HttpService) RegisterRoutes(routes ...*api.Route) {
	for _, route := range routes {
		r := route
		path := r.Path
		if strings.HasSuffix(path, "/") {
			path = path[:len(path)-1]
		}
		r.Method.ToFunctionFromEngine(h.engine)(path, func(c *gin.Context) {
			var needToCancel bool

			if r.IsAuthenticated {
				middlewares.AuthMiddleware(h.service, c, func(user api.User) {
					if c.Writer.Written() {
						needToCancel = true
						return
					}

					r.AuthenticateMiddleware(c, user, h.service)
					if c.Writer.Written() {
						needToCancel = true
						return
					}

					for _, middleware := range r.Middlewares {
						if !middleware(c, user, h.service) {
							needToCancel = true
							return
						}
					}
					r.Handler(c, user, h.service)
					needToCancel = true
					return
				})
				if c.Writer.Written() || needToCancel {
					return
				}
			}

			for _, middleware := range r.Middlewares {
				if !middleware(c, nil, h.service) {
					return
				}
			}

			r.Handler(c, nil, h.service)
		})
	}
}

func (h *HttpService) RegisterControllers(controllers ...*api.Controller) {
	for _, controller := range controllers {
		for _, route := range controller.Routes {
			route.Path = controller.Path + route.Path
		}
		h.RegisterRoutes(controller.Routes...)
	}
}

func (h *HttpService) Settings() api.ServiceSettings {
	return api.ServiceSettings{
		MustWaitForStart: false,
		Priority:         5,
	}
}

func (h *HttpService) logger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(params gin.LogFormatterParams) string {
			if ip := net.ParseIP(params.ClientIP); ip.IsLoopback() {
				return ""
			}

			if params.Path == routes.HealthRoute().Path {
				return ""
			}

			var statusColor, methodColor, resetColor string
			if params.IsOutputColor() {
				statusColor = params.StatusCodeColor()
				methodColor = params.MethodColor()
				resetColor = params.ResetColor()
			}

			if params.Latency > time.Minute {
				params.Latency = params.Latency.Truncate(time.Second)
			}
			return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
				params.TimeStamp.Format("2006/01/02 - 15:04:05"),
				statusColor, params.StatusCode, resetColor,
				params.Latency,
				params.ClientIP,
				methodColor, params.Method, resetColor,
				params.Path,
				params.ErrorMessage,
			)
		},
	})
}
