package serviceImpl

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/service"
)

type platformService struct {
	// Add any dependencies here (e.g., HTTP client for real API calls)
}

func NewPlatformService() service.PlatformService {
	return &platformService{}
}

// ValidateURL validates and extracts metadata from YouTube/TikTok URLs
func (s *platformService) ValidateURL(url string) (*models.PlatformValidationResult, error) {
	// YouTube URL patterns
	youtubePatterns := []string{
		`^https?://(?:www\.)?youtube\.com/watch\?v=([a-zA-Z0-9_-]{11})`,
		`^https?://youtu\.be/([a-zA-Z0-9_-]{11})`,
		`^https?://(?:www\.)?youtube\.com/live/([a-zA-Z0-9_-]{11})`,
		`^https?://(?:www\.)?youtube\.com/shorts/([a-zA-Z0-9_-]{11})`,
		`^https?://(?:www\.)?youtube\.com/@([a-zA-Z0-9_.-]+)`,
		`^https?://(?:www\.)?youtube\.com/channel/([a-zA-Z0-9_-]+)`,
		`^https?://(?:www\.)?youtube\.com/c/([a-zA-Z0-9_-]+)`,
	}

	// TikTok URL patterns
	tiktokPatterns := []string{
		`^https?://(?:www\.)?tiktok\.com/@([a-zA-Z0-9_.]+)/video/(\d+)`,
		`^https?://vm\.tiktok\.com/([a-zA-Z0-9]+)`,
		`^https?://(?:www\.)?tiktok\.com/@([a-zA-Z0-9_.]+)/live`,
		`^https?://(?:www\.)?tiktok\.com/@([a-zA-Z0-9_.]+)`,
	}

	// Check YouTube patterns
	for _, pattern := range youtubePatterns {
		if matched, _ := regexp.MatchString(pattern, url); matched {
			return s.extractYouTubeMetadata(url, pattern)
		}
	}

	// Check TikTok patterns
	for _, pattern := range tiktokPatterns {
		if matched, _ := regexp.MatchString(pattern, url); matched {
			return s.extractTikTokMetadata(url, pattern)
		}
	}
	
	return &models.PlatformValidationResult{
		IsValid:     false,
		Platform:    "",
		ContentType: "",
		Metadata:    nil,
	}, nil
}

// extractYouTubeMetadata extracts metadata from YouTube URLs
func (s *platformService) extractYouTubeMetadata(url, pattern string) (*models.PlatformValidationResult, error) {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(url)

	contentType := "video"
	if strings.Contains(url, "/live/") {
		contentType = "live"
	} else if strings.Contains(url, "/shorts/") {
		contentType = "short"
	} else if strings.Contains(url, "/@") || strings.Contains(url, "/channel/") || strings.Contains(url, "/c/") {
		contentType = "channel"
	}

	metadata := map[string]interface{}{
		"platform": "youtube",
		"url":      url,
	}

	// Mock metadata - in production, use YouTube Data API
	if len(matches) > 1 {
		if contentType == "channel" {
			metadata["channel_id"] = matches[1]
			metadata["creator"] = matches[1]
			metadata["title"] = fmt.Sprintf("Channel: %s", matches[1])
			metadata["thumbnail"] = "https://yt3.ggpht.com/sample_channel_avatar.jpg"
			metadata["subscriber_count"] = 15000
		} else {
			metadata["video_id"] = matches[1]
			metadata["title"] = "Sample YouTube Video Title"
			metadata["creator"] = "Sample Creator"
			metadata["thumbnail"] = fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", matches[1])
			metadata["duration"] = 300 // 5 minutes
			metadata["view_count"] = 1250
			metadata["like_count"] = 89
		}
	}

	if contentType == "live" {
		metadata["is_live"] = true
		metadata["viewer_count"] = 45
	}

	return &models.PlatformValidationResult{
		IsValid:     true,
		Platform:    "youtube",
		ContentType: contentType,
		Metadata:    metadata,
	}, nil
}

// extractTikTokMetadata extracts metadata from TikTok URLs
func (s *platformService) extractTikTokMetadata(url, pattern string) (*models.PlatformValidationResult, error) {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(url)

	contentType := "video"
	if strings.Contains(url, "/live") {
		contentType = "live"
	} else if strings.Contains(url, "/@") && !strings.Contains(url, "/video/") {
		contentType = "profile"
	}

	metadata := map[string]interface{}{
		"platform": "tiktok",
		"url":      url,
	}
	
	// Mock metadata - in production, use TikTok API
	if len(matches) > 1 {
		if contentType == "profile" {
			metadata["username"] = matches[1]
			metadata["creator"] = matches[1]
			metadata["title"] = fmt.Sprintf("TikTok Profile: @%s", matches[1])
			metadata["thumbnail"] = "https://p16-sign-va.tiktokcdn.com/sample_avatar.jpeg"
			metadata["follower_count"] = 8500
		} else {
			metadata["username"] = matches[1]
			if len(matches) > 2 {
				metadata["video_id"] = matches[2]
			}
			metadata["title"] = "Sample TikTok Video"
			metadata["creator"] = matches[1]
			metadata["thumbnail"] = "https://p16-sign-va.tiktokcdn.com/sample_video_cover.jpeg"
			metadata["duration"] = 30 // 30 seconds
			metadata["view_count"] = 2500
			metadata["like_count"] = 150
		}
	}
	
	if contentType == "live" {
		metadata["is_live"] = true
		metadata["viewer_count"] = 25
	}

	return &models.PlatformValidationResult{
		IsValid:     true,
		Platform:    "tiktok",
		ContentType: contentType,
		Metadata:    metadata,
	}, nil
}

// IsLiveStream checks if a URL represents a live stream
func (s *platformService) IsLiveStream(url string) bool {
	livePatterns := []string{
		`youtube\.com/live/`,
		`tiktok\.com/@[^/]+/live`,
	}
	
	for _, pattern := range livePatterns {
		if matched, _ := regexp.MatchString(pattern, url); matched {
			return true
		}
	}
	return false
}

// GetPlatformFromURL extracts platform type from URL
func (s *platformService) GetPlatformFromURL(url string) string {
	if strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be") {
		return "youtube"
	}
	if strings.Contains(url, "tiktok.com") || strings.Contains(url, "vm.tiktok.com") {
		return "tiktok"
	}
	return ""
}

// GetContentTypeFromURL extracts content type from URL
func (s *platformService) GetContentTypeFromURL(url string) string {
	if strings.Contains(url, "/live/") || strings.Contains(url, "/live") {
		return "live"
	}
	if strings.Contains(url, "/shorts/") {
		return "short"
	}
	if strings.Contains(url, "/@") || strings.Contains(url, "/channel/") || strings.Contains(url, "/c/") {
		return "channel"
	}
	return "video"
} 