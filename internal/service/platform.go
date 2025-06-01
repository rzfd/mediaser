package service

import "github.com/rzfd/mediashar/internal/models"

type PlatformService interface {
	ValidateURL(url string) (*models.PlatformValidationResult, error)
	IsLiveStream(url string) bool
	GetPlatformFromURL(url string) string
	GetContentTypeFromURL(url string) string
} 