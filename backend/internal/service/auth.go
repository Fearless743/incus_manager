package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"incus-manager/internal/model"
)

type AuthService struct {
	DB        *gorm.DB
	JWTSecret []byte
}

func NewAuthService(db *gorm.DB, jwtSecret string) *AuthService {
	return &AuthService{
		DB:        db,
		JWTSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) Login(username, password string) (*model.User, string, error) {
	var user model.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// TODO: Verify password hash with bcrypt
	if user.ID == 0 {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.GenerateToken(&user)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *AuthService) GenerateToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"username":  user.Username,
		"role":      user.Role,
		"expires_at": time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(s.JWTSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.JWTSecret, nil
	})
}
