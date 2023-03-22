package controllers

import (
	"niromash-api/api"
	"niromash-api/services/http/routes/stats"
)

func StatsController() *api.Controller {
	return &api.Controller{
		Path: "/stats",
		Routes: []*api.Route{
			stats.StatsRoute(),
		},
	}
}
