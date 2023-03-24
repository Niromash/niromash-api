package controllers

import (
	"github.com/Niromash/niromash-api/api"
	"github.com/Niromash/niromash-api/services/http/routes/stats"
)

func StatsController() *api.Controller {
	return &api.Controller{
		Path: "/stats",
		Routes: []*api.Route{
			stats.StatsRoute(),
		},
	}
}
