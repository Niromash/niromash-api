package services

import (
	"gorm.io/gorm"
	"niromash-api/api"
	"niromash-api/model"
)

var _ api.UsersService = (*UsersService)(nil)

type UsersService struct {
	service api.MainService
}

func NewUsersService() *UsersService {
	return &UsersService{}
}

func (u *UsersService) Init(service api.MainService) error {
	u.service = service
	return nil
}

func (u *UsersService) GetUser(id uint, withPassword ...bool) (api.User, error) {
	var user model.User
	statement := u.service.Databases().Postgres().GetClient()
	if len(withPassword) == 0 || !withPassword[0] {
		statement = statement.Omit("password")
	}
	if err := statement.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, api.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *UsersService) GetUserByUsername(username string, withPassword ...bool) (api.User, error) {
	var user model.User
	statement := u.service.Databases().Postgres().GetClient()
	if len(withPassword) == 0 || !withPassword[0] {
		statement = statement.Omit("password")
	}
	if err := statement.First(&user, "username = ?", username).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, api.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *UsersService) GetUserByEmail(email string, withPassword ...bool) (api.User, error) {
	var user model.User
	statement := u.service.Databases().Postgres().GetClient()
	if len(withPassword) == 0 || !withPassword[0] {
		statement = statement.Omit("password")
	}
	if err := statement.First(&user, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, api.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *UsersService) ListUsers() (usersResult []api.User, err error) {
	rows, err := u.service.Databases().Postgres().GetClient().Omit("password").Rows()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, api.ErrUserNotFound
		}
		return
	}

	for rows.Next() {
		var user model.User
		if err = u.service.Databases().Postgres().GetClient().ScanRows(rows, &user); err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, api.ErrUserNotFound
			}
			return
		}
		usersResult = append(usersResult, &user)
	}

	return
}

func (u *UsersService) Register(user api.User) error {
	return u.service.Databases().Postgres().GetClient().Create(user).Error
}

func (u *UsersService) IsExist(email string) bool {
	var user model.User
	if err := u.service.Databases().Postgres().GetClient().First(&user, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
	}

	return true
}
