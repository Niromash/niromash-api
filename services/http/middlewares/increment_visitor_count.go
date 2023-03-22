package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"niromash-api/api"
	"strings"
)

func IncrementVisitorCount(c *gin.Context, user api.User, mainService api.MainService) (next bool) {
	if !strings.Contains(c.Request.UserAgent(), "Mozilla") {
		return true
	}

	if err := mainService.Databases().Redis().GetClient().Base().Incr(context.Background(), "personal:states:visitors").Err(); err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": "internal server error",
			"error":   err.Error(),
		})
		return false
	}

	return true
}
