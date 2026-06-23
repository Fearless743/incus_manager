package service

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"incus-manager/internal/model"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) Login(username, password string) (*model.User, error) {
	var user model.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	// TODO: Verify password hash
	if user.ID == 0 {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

func (s *AuthService) GenerateToken(user *model.User) (string, error) {
	// TODO: Implement JWT token generation
	return "", errors.New("not implemented")
}
