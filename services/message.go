package services

import (
	"context"
	"fmt"
	"github.com/Niromash/niromash-api/api"
	"github.com/Niromash/niromash-api/model"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

var _ api.MessagesService = (*MessageService)(nil)

type MessageService struct {
	service api.MainService
}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (m *MessageService) Init(service api.MainService) error {
	m.service = service
	return nil
}

func (m *MessageService) GetAllMessagesFromScopeId(scopeId uuid.UUID) (messages []*model.Message, err error) {
	tx := m.service.Databases().Postgres().GetClient().
		Preload(clause.Associations).
		Joins("JOIN scopes ON scopes.id = messages.scope_id").
		Find(&messages, "scopes.id = ?", scopeId)

	if tx.Error != nil {
		err = tx.Error
	}
	return
}

func (m *MessageService) GetAllMessagesFromScope(scope string) (messages []*model.Message, err error) {
	tx := m.service.Databases().Postgres().GetClient().
		Preload(clause.Associations).
		Joins("JOIN scopes ON scopes.id = messages.scope_id").
		Find(&messages, "lower(scopes.name) = ?", strings.ToLower(scope))

	if tx.Error != nil {
		err = tx.Error
	}
	return
}

func (m *MessageService) GetMessage(scope, key string) (message *model.Message, err error) {
	redisKey := fmt.Sprintf("message:%s:%s", scope, key)
	var msg model.Message
	if err = m.service.Databases().Redis().GetClient().GetJsonObject(redisKey, ".", &msg); err != nil {
		if err == redis.Nil {
			if err = m.service.Databases().Postgres().GetClient().
				Preload(clause.Associations).
				Joins("JOIN scopes ON scopes.id = messages.scope_id").
				First(&message, "lower(scopes.name) = ? AND lower(key) = ?", strings.ToLower(scope), strings.ToLower(key)).Error; err != nil {
				return
			}

			if _, err = m.service.Databases().Redis().GetClient().ReJson().JSONSet(redisKey, ".", message); err != nil {
				return
			}
			if err = m.service.Databases().Redis().GetClient().Base().Expire(context.TODO(), redisKey, 15*time.Minute).Err(); err != nil {
				return
			}

			return
		}

		return nil, err
	}

	return &msg, nil
}

func (m *MessageService) AddMessage(message *model.Message) (err error) {
	return m.service.Databases().Postgres().GetClient().Create(&message).Error
}

func (m *MessageService) UpdateMessage(message *model.Message) (err error) {
	tx := m.service.Databases().Postgres().GetClient().Begin()
	defer tx.Rollback()
	if err = tx.Omit("Scope").Updates(&m).Error; err != nil {
		return
	}

	session := tx.Session(&gorm.Session{FullSaveAssociations: true})
	if err = session.Model(&model.Message{}).Association("Arguments").Replace(message.Arguments); err != nil {
		return
	}

	if err = session.Model(&model.Message{}).Association("Translations").Replace(message.Translations); err != nil {
		return
	}

	return tx.Commit().Error
}

func (m *MessageService) AddTranslations(message *model.Message) (err error) {
	return m.service.Databases().Postgres().GetClient().Model(&m).Association("Translations").Append(message.Translations)
}

func (m *MessageService) GetScopeFromName(scopeName string) (scope *model.Scope, err error) {
	tx := m.service.Databases().Postgres().GetClient().First(&scope, "lower(name) = ?", strings.ToLower(scopeName))
	if tx.Error != nil {
		err = tx.Error
	}
	return
}
