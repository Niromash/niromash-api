package messages

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"niromash-api/api"
	"niromash-api/model"
)

func UpdateMessageRoute() *api.Route {
	return &api.Route{
		Path:   "/:scopeId",
		Method: api.MethodPut,
		AuthenticateMiddleware: func(c *gin.Context, user api.User, service api.MainService) {
			if !user.HasPermission("message.update") {
				c.AbortWithStatusJSON(403, gin.H{
					"message": "You do not have the permission!",
				})
				return
			}

			return
		},
		Handler: func(c *gin.Context, user api.User, service api.MainService) {
			var msgs []*model.Message
			if err := c.ShouldBindJSON(&msgs); err != nil {
				return
			}

			if len(msgs) == 0 {
				c.JSON(400, gin.H{
					"message": "You need to give translations to be saved!",
				})
				return
			}

			scopeId, err := uuid.Parse(c.Param("scopeId"))
			if err != nil {
				c.JSON(400, gin.H{
					"message": "Invalid scope id!",
				})
				return
			}

			for _, msg := range msgs {
				if msg.ID == uuid.Nil {
					c.JSON(400, gin.H{
						"message": "You need to give an id for each message!",
					})
					return
				}
				if msg.Scope.ID != scopeId {
					c.JSON(400, gin.H{
						"message": "You can only give translations for one scope at a time!",
					})
					return
				}

				msg.ScopeID = msg.Scope.ID
				if err = service.Messages().UpdateMessage(msg); err != nil {
					return
				}
			}

			c.Status(204)
		},
	}
}
