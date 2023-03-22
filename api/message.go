package api

import (
	"github.com/google/uuid"
	"niromash-api/model"
)

type MessagesService interface {
	ServiceInitializer
	GetAllMessagesFromScopeId(scopeId uuid.UUID) (messages []*model.Message, err error)
	GetAllMessagesFromScope(scope string) (messages []*model.Message, err error)
	GetMessage(scope, key string) (message *model.Message, err error)
	AddMessage(message *model.Message) (err error)
	UpdateMessage(message *model.Message) (err error)
	AddTranslations(message *model.Message) (err error)
	GetScopeFromName(scopeName string) (scope *model.Scope, err error)
}
