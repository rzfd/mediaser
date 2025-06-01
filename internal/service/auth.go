package service

import (
	"github.com/rzfd/mediashar/internal/middleware"
	"github.com/rzfd/mediashar/internal/models"
)

type AuthService interface {
	GenerateToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (*middleware.JWTClaims, error)
	RefreshToken(tokenString string) (string, error)
} 