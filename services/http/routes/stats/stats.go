package stats

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"niromash-api/api"
	"time"
)

func StatsRoute() *api.Route {
	return &api.Route{
		Path:   "/",
		Method: api.MethodGet,
		Handler: func(c *gin.Context, user api.User, mainService api.MainService) {
			totalDevTimeResp, err := mainService.Stats().GetTotalDevTime()
			if err != nil {
				c.JSON(500, gin.H{
					"message": "Error getting total dev time",
					"error":   err.Error(),
				})
				return
			}

			bestDevTimeDay, err := mainService.Stats().GetBestDevTimeDay()
			if err != nil {
				c.JSON(500, gin.H{
					"message": "Error getting best dev time day",
					"error":   err.Error(),
				})
				return
			}

			isDeveloping, err := mainService.Stats().IsDeveloping()
			if err != nil {
				c.JSON(500, gin.H{
					"message": "Error getting developing state",
					"error":   err.Error(),
				})
				return
			}

			visitorCount, err := mainService.Stats().GetVisitorCount()
			if err != nil && err != redis.Nil {
				c.JSON(500, gin.H{
					"message": "Error getting visitor count",
					"error":   err.Error(),
				})
				return
			}

			storedRepositories, err := mainService.Stats().ListRepositories()
			if err != nil {
				c.JSON(500, gin.H{
					"message": "Error getting repositories",
					"error":   err.Error(),
				})
				return
			}

			c.JSON(200, gin.H{
				"totalDevTime": time.Duration(totalDevTimeResp.Data.TotalSeconds).Seconds(),
				"totalSince":   totalDevTimeResp.Data.Range.Start,
				"bestDevTime":  time.Duration(bestDevTimeDay).Seconds(),
				"developing":   isDeveloping,
				"visitorCount": visitorCount,
				"repositories": storedRepositories,
			})
		},
	}
}
