package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/pkg/utils"
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID     uint   `json:"user_id"`
	Email      string `json:"email"`
	IsStreamer bool   `json:"is_streamer"`
	jwt.RegisteredClaims
}

// JWTMiddleware creates a JWT middleware
func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Missing authorization header", nil))
			}

			// Check if token starts with "Bearer "
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid authorization format", nil))
			}

			// Parse and validate token
			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid token", err))
			}

			if !token.Valid {
				return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Token is not valid", nil))
			}

			// Extract claims
			claims, ok := token.Claims.(*JWTClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid token claims", nil))
			}

			// Set user info in context
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("is_streamer", claims.IsStreamer)

			return next(c)
		}
	}
}

// OptionalJWTMiddleware creates an optional JWT middleware (doesn't fail if no token)
func OptionalJWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return next(c)
			}

			// Check if token starts with "Bearer "
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return next(c)
			}

			// Parse and validate token
			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return next(c)
			}

			// Extract claims
			claims, ok := token.Claims.(*JWTClaims)
			if ok {
				// Set user info in context
				c.Set("user_id", claims.UserID)
				c.Set("user_email", claims.Email)
				c.Set("is_streamer", claims.IsStreamer)
			}

			return next(c)
		}
	}
}

// StreamerOnlyMiddleware ensures only streamers can access the endpoint
func StreamerOnlyMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isStreamer, ok := c.Get("is_streamer").(bool)
			if !ok || !isStreamer {
				return c.JSON(http.StatusForbidden, utils.ErrorResponse("Access denied: Streamers only", nil))
			}
			return next(c)
		}
	}
} 