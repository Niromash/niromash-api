package controllers

import (
	"github.com/Niromash/niromash-api/api"
	"github.com/Niromash/niromash-api/services/http/routes/messages"
)

func MessageController() *api.Controller {
	return &api.Controller{
		Path: "/messages",
		Routes: []*api.Route{
			messages.CreateMessageRoute(),
			messages.ListMessageFromScopeRoute(),
			{Path: "/:scope/:key", Method: api.MethodGet, Handler: messages.ListMessageFromScopeRoute().Handler},
			messages.UpdateMessageRoute(),
			messages.DeleteMessageRoute(),
			{Path: "/:scope/:key", Method: api.MethodDelete, Handler: messages.DeleteMessageRoute().Handler},
		},
	}
}
