package controllers

import (
	"niromash-api/api"
	"niromash-api/services/http/routes/auth"
)

func AuthController() *api.Controller {
	return &api.Controller{
		Path: "/auth",
		Routes: []*api.Route{
			auth.RegisterRoute(),
			auth.LoginRoute(),
		},
	}
}
