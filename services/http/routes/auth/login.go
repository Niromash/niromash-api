package auth

import (
	"github.com/Niromash/niromash-api/api"
	"github.com/Niromash/niromash-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type loginBody struct {
	Email    string `json:"email" xml:"email" yaml:"email" validate:"required"`
	Password string `json:"password" xml:"password" yaml:"password" validate:"required,min=3"`
}

func (l loginBody) IsFulfilled() bool {
	return len(l.Email) != 0 && len(l.Password) != 0
}

func LoginRoute() *api.Route {
	return &api.Route{
		Path:   "/login",
		Method: api.MethodPost,
		Handler: func(c *gin.Context, _ api.User, mainService api.MainService) {
			var body loginBody
			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			if !body.IsFulfilled() {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "you need to fulfill the form",
				})
				return
			}

			user, err := mainService.Users().GetUserByEmail(body.Email, true)
			if err != nil {
				if err == api.ErrUserNotFound {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "invalid username or password",
					})
					return
				}
				c.String(http.StatusInternalServerError, err.Error())
				return
			}

			if !utils.Compare([]byte(user.GetPassword()), []byte(body.Password)) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid username or password",
				})
				return
			}

			tokenPair, err := utils.GenerateTokenPair(user)
			if err != nil {
				c.String(http.StatusInternalServerError, "Unable to generate token.")
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"accessToken":  tokenPair[0],
				"refreshToken": tokenPair[1],
			})
		},
	}
}
