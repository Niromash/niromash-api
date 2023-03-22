package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"niromash-api/api"
	"niromash-api/model"
	"niromash-api/utils"
)

type registerBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r registerBody) IsFulfilled() bool {
	return len(r.Username) != 0 && len(r.Email) != 0 && len(r.Password) != 0
}

func RegisterRoute() *api.Route {
	return &api.Route{
		Path:   "/register",
		Method: api.MethodPost,
		Handler: func(c *gin.Context, user api.User, mainService api.MainService) {
			var body registerBody
			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Unable to bind JSON to RegisterBody struct.",
					"error":   err.Error(),
				})
				return
			}

			if !body.IsFulfilled() {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "You need to fulfill the form.",
				})
				return
			}

			if mainService.Users().IsExist(body.Email) {
				c.JSON(http.StatusConflict, gin.H{
					"message": "An account with this email already exists!",
				})
				return
			}

			hash, err := utils.Crypt(body.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Unable to crypt given password!",
				})
				return
			}

			if err = mainService.Users().Register(&model.User{
				Username: body.Username,
				Email:    body.Email,
				Password: string(hash),
			}); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "unable to create the account",
					"error":   err.Error(),
				})
				return
			}

			c.Status(http.StatusCreated)
		},
	}
}
