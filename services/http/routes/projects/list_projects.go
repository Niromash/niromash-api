package projects

import (
	"github.com/Niromash/niromash-api/api"
	"github.com/Niromash/niromash-api/model"
	"github.com/Niromash/niromash-api/services/http/middlewares"
	"github.com/gin-gonic/gin"
)

func ListProjectsRoute() *api.Route {
	return &api.Route{
		Path:        "/",
		Method:      api.MethodGet,
		Middlewares: []api.MiddlewareHandler{middlewares.IncrementVisitorCount},
		Handler: func(c *gin.Context, user api.User, mainService api.MainService) {
			projects, err := mainService.Projects().ListProjects()
			if err != nil {
				c.AbortWithStatusJSON(500, gin.H{
					"message": "internal server error",
					"error":   err.Error(),
				})
			}

			if len(projects) == 0 {
				projects = []*model.Project{} // return empty array instead of null
			}

			c.JSON(200, projects)
		},
	}
}
