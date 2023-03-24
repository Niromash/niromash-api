package messages

import (
	"github.com/Niromash/niromash-api/api"
	"github.com/Niromash/niromash-api/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thoas/go-funk"
	"strconv"
	"strings"
)

func CreateMessageRoute() *api.Route {
	return &api.Route{
		Path:            "/",
		Method:          api.MethodPost,
		IsAuthenticated: true,
		AuthenticateMiddleware: func(c *gin.Context, user api.User, mainService api.MainService) {
			if !user.HasPermission("message.create") {
				c.AbortWithStatusJSON(403, gin.H{
					"message": "You do not have the permission!",
				})
				return
			}

			return
		},
		Handler: func(c *gin.Context, user api.User, mainService api.MainService) {
			var msgs []*model.Message
			if err := c.ShouldBindJSON(&msgs); err != nil {
				c.JSON(400, gin.H{
					"message": "You need to give translations to be saved!",
					"error":   err.Error(),
				})
				return
			}

			if len(msgs) == 0 {
				c.JSON(400, gin.H{
					"message": "You need to give translations to be saved!",
				})
				return
			}

			scope := msgs[0].Scope

			if scope.ID == uuid.Nil && msgs[0].Scope.Name != "" {
				scopeFromName, err := mainService.Messages().GetScopeFromName(msgs[0].Scope.Name)
				if err != nil {
					c.JSON(500, gin.H{
						"message": "An error occurred while trying to get the scope!",
						"error":   err.Error(),
					})
					return
				}

				scope = *scopeFromName
			}

			for _, msg := range msgs {
				msg.Scope = scope
			}

			var allMessages []*model.Message
			var err error
			if allMessages, err = mainService.Messages().GetAllMessagesFromScopeId(scope.ID); err != nil {
				for _, msg := range msgs {
					if err = mainService.Messages().AddMessage(msg); err != nil {
						c.JSON(500, gin.H{
							"message": "An error occurred while trying to add the translations!",
							"error":   err.Error(),
						})
						return
					}
				}

				c.JSON(200, gin.H{
					"message": "The scope: " + scope.ID.String() + " has been created with the translations (" + strconv.Itoa(len(msgs)) + "): " + strings.Join(funk.Map(msgs, func(t *model.Message) string {
						return t.Key
					}).([]string), ", "),
				})
				return
			}

			keys := funk.Map(allMessages, func(t *model.Message) string {
				return t.Key
			}).([]string)

			translationsToAdd := funk.Filter(msgs, func(t *model.Message) bool {
				return !funk.ContainsString(keys, t.Key)
			}).([]*model.Message)

			if len(translationsToAdd) == 0 {
				c.JSON(409, gin.H{
					"message": "The translations given already exist!",
				})
				return
			}

			for _, msg := range msgs {
				if err = mainService.Messages().AddMessage(msg); err != nil {
					c.JSON(500, gin.H{
						"message": "An error occurred while trying to add the translations!",
						"error":   err.Error(),
					})
					return
				}
			}

			c.JSON(200, gin.H{
				"message": "Successfully added translations (" + strconv.Itoa(len(translationsToAdd)) + "): " + strings.Join(funk.Map(translationsToAdd, func(t *model.Message) string {
					return t.Key
				}).([]string), ", ") + " to the scope: " + scope.ID.String(),
			})
		},
	}
}
