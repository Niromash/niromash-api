package controllers

import (
	"niromash-api/api"
	"niromash-api/services/http/routes/projects"
)

func ProjectsController() *api.Controller {
	return &api.Controller{
		Path: "/projects",
		Routes: []*api.Route{
			projects.ListProjectsRoute(),
			projects.GetProjectRoute(),
		},
	}
}
