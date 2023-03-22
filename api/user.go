package api

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
)

type UsersService interface {
	ServiceInitializer
	GetUser(id uint, withPassword ...bool) (User, error)
	GetUserByUsername(username string, withPassword ...bool) (User, error)
	GetUserByEmail(email string, withPassword ...bool) (User, error)
	ListUsers() ([]User, error)
	Register(user User) error
	IsExist(email string) bool
}

type User interface {
	GetId() uint
	GetUsername() string
	GetEmail() string
	GetPassword() string
	GetPermissions() []string
	HasPermission(permission string) bool
}
