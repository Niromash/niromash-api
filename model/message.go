package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Scope struct {
	ID        uuid.UUID      `json:"id" gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"unique;not null"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type MessageArgument struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	MessageID   uuid.UUID `json:"-" gorm:"type:uuid;"`
	Message     Message   `json:"-" gorm:"foreignkey:MessageID;association_foreignkey:ID;constraint:OnUpdate:Cascade,OnDelete:Cascade;"`
	Name        string    `json:"name" validate:"required,min=3"`
	Description string    `json:"description" validate:"required,min=3"`
}

type MessageTranslation struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	MessageID uuid.UUID `json:"-" gorm:"type:uuid;"`
	Message   Message   `json:"-" gorm:"foreignkey:MessageID;association_foreignkey:ID;constraint:OnUpdate:Cascade,OnDelete:Cascade;"`
	Locale    string    `json:"locale" validate:"required,min=3"`
	Value     string    `json:"value" validate:"required,min=3"`
}

type Message struct {
	ID           uuid.UUID             `json:"id" gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	ScopeID      uuid.UUID             `json:"-" gorm:"type:uuid;"`
	Scope        Scope                 `json:"scope" gorm:"foreignkey:ScopeID;association_foreignkey:ID;constraint:OnUpdate:Cascade,OnDelete:Cascade;"`
	Key          string                `json:"key" validate:"required,min=3"`
	Description  string                `json:"description" validate:"required,min=3"`
	Arguments    []*MessageArgument    `json:"arguments" gorm:"foreignkey:MessageID;association_foreignkey:ID;constraint:OnUpdate:Cascade,OnDelete:Cascade;"`
	Translations []*MessageTranslation `json:"translations" gorm:"foreignkey:MessageID;association_foreignkey:ID;constraint:OnUpdate:Cascade,OnDelete:Cascade;"`
}
