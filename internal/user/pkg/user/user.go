package user

import (
	"crypto/md5"
	"errors"
	"fmt"

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

func (s *UserService) CreateUser(username, password string) (int64, error) {
	if password == "" || username == "" {
		return 0, errors.New("用户名或密码不能为空")
	}
	if len(password) < 6 {
		return 0, errors.New("密码长度不能小于6位")
	}
	if len(password) > 32 {
		return 0, errors.New("密码长度不能大于32位")
	}

	db := s.dbInstance()
	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	user := User{
		Username: username,
		Password: password,
	}
	err := db.Create(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, errors.New("用户名已被占用")
		}
		return 0, err
	}
	return user.Id, nil
}

func (s *UserService) Login(username, password string) (int64, error) {
	db := s.dbInstance()
	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	user := User{}
	err := db.Where("username = ? AND password = ?", username, password).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("用户名或密码错误")
		}
		return 0, err
	}
	return user.Id, nil

}

func (s *UserService) GetUserById(id int64) (*User, error) {
	db := s.dbInstance()
	user := User{}
	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByIds(ids []int64) ([]*User, error) {
	db := s.dbInstance()
	users := make([]*User, 0)
	err := db.Where("id IN ?", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
