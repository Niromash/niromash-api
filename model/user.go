package model

import (
	"github.com/lib/pq"
	"strings"
)

type User struct {
	ID          uint           `json:"id" gorm:"primary_key;auto_increment;not null"`
	Username    string         `json:"username" gorm:"type:varchar(255);not null;unique"`
	Email       string         `json:"email" gorm:"type:varchar(255);not null;unique"`
	Password    string         `json:"password" gorm:"type:varchar(255);not null"`
	Permissions pq.StringArray `json:"permissions" gorm:"type:varchar(255)[]"`
}

func (u *User) GetId() uint {
	return u.ID
}

func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetPassword() string {
	return u.Password
}

func (u *User) GetPermissions() []string {
	return u.Permissions
}

// HasPermission examples:
// permission for all actions: *
// permission for all messages actions: messages.*
// permission for create message action: messages.create
// permission for delete message action: messages.delete
func (u *User) HasPermission(permission string) bool {
	for _, p := range u.Permissions {
		if p == permission || p == "*" {
			return true
		} else if strings.HasSuffix(p, ".*") {
			if strings.HasPrefix(permission, strings.TrimSuffix(p, ".*")) {
				return true
			}
		}
	}
	return false
}
