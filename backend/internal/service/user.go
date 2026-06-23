package service

import (
	"errors"

	"gorm.io/gorm"
	"incus-manager/internal/model"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) CreateUser(username, email, password string) (*model.User, error) {
	user := &model.User{
		Username: username,
		Email:    email,
		Role:     "user",
	}

	// TODO: Hash password before saving
	if err := s.DB.Create(user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
