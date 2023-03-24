package routes

import (
	"github.com/Niromash/niromash-api/api"
	"github.com/gin-gonic/gin"
)

func HealthRoute() *api.Route {
	return &api.Route{
		Path:   "/health",
		Method: api.MethodGet,
		Handler: func(c *gin.Context, user api.User, mainService api.MainService) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		},
	}
}
