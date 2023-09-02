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

func (c *UserService) Register(username string, password string) (*User, error) {
	db := c.dbInstance()
	user := &User{
		Username: username,
		Password: password,
	}
	err := db.Create(user).Error
	return user, err
}

func (c *UserService) Login(username string, password string) (*User, error) {
	db := c.dbInstance()
	user := &User{}
	err := db.Where("username=? and password=?", username, password).First(user).Error
	return user, err
}

func (c *UserService) GetUserInfo(Id int64) (*User, error) {
	db := c.dbInstance()
	user := &User{}
	err := db.Where("id = ?", Id).First(user).Error
	return user, err
}
