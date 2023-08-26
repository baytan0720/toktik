package user

import (
	"gorm.io/gorm"

	"toktik/pkg/db/model"
)

type User = model.User

type UserService struct {
	dbInstance func() *gorm.DB
}

func NewUserService(db func() *gorm.DB) *UserService {
	return &UserService{
		dbInstance: db,
	}
}
