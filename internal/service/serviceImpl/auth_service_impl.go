package serviceImpl

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rzfd/mediashar/internal/middleware"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/service"
)

type authService struct {
	jwtSecret   string
	tokenExpiry time.Duration
}

func NewAuthService(jwtSecret string, tokenExpiryHours int) service.AuthService {
	return &authService{
		jwtSecret:   jwtSecret,
		tokenExpiry: time.Duration(tokenExpiryHours) * time.Hour,
	}
}

// GenerateToken generates a JWT token for the user
func (s *authService) GenerateToken(user *models.User) (string, error) {
	claims := &middleware.JWTClaims{
		UserID:     user.ID,
		Email:      user.Email,
		IsStreamer: user.IsStreamer,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mediashar",
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateToken validates a JWT token and returns the claims
func (s *authService) ValidateToken(tokenString string) (*middleware.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &middleware.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenNotValidYet
	}

	claims, ok := token.Claims.(*middleware.JWTClaims)
	if !ok {
		return nil, jwt.ErrTokenMalformed
	}

	return claims, nil
}

// RefreshToken generates a new token from an existing valid token
func (s *authService) RefreshToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Create new claims with extended expiry
	newClaims := &middleware.JWTClaims{
		UserID:     claims.UserID,
		Email:      claims.Email,
		IsStreamer: claims.IsStreamer,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mediashar",
			Subject:   claims.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return token.SignedString([]byte(s.jwtSecret))
} 