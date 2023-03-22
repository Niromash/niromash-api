package projects

import (
	"github.com/gin-gonic/gin"
	"niromash-api/api"
	"strconv"
)

func GetProjectRoute() *api.Route {
	return &api.Route{
		Path:   "/:id",
		Method: api.MethodGet,
		Handler: func(c *gin.Context, user api.User, mainService api.MainService) {
			id := c.Param("id")
			projectId, err := strconv.Atoi(id)
			if err != nil {
				c.AbortWithStatusJSON(400, gin.H{
					"message": "invalid id",
				})
				return
			}
			project, err := mainService.Projects().GetProject(uint(projectId))
			if err != nil {
				if err == api.ErrProjectNotFound {
					c.AbortWithStatusJSON(404, gin.H{
						"message": "project not found",
					})
					return
				}
				c.AbortWithStatusJSON(500, gin.H{
					"message": "internal server error",
					"error":   err.Error(),
				})
			}

			c.JSON(200, project)
		},
	}
}
