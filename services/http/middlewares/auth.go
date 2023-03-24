package middlewares

import (
	"fmt"
	"github.com/Niromash/niromash-api/api"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

func AuthMiddleware(service api.MainService, c *gin.Context, funcc func(user api.User)) {
	token := ExtractToken(c)
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token required"})
		return
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if claims["type"] == "refresh" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "refresh token not allowed"})
			return
		}

		user, err := service.Users().GetUserByEmail(claims["email"].(string))
		if err != nil {
			if err == api.ErrUserNotFound {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unable to find user relative to this account!"})
			}
			return
		}
		funcc(user)
	}
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
