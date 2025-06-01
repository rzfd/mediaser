# 🎥 MediaShar Social Media Integration Guide

Panduan lengkap untuk mengintegrasikan MediaShar dengan platform streaming eksternal seperti YouTube dan TikTok.

## 🎯 **Overview**

MediaShar dapat diintegrasikan dengan platform streaming populer untuk memungkinkan donasi langsung ke konten creator di berbagai platform. Integrasi ini mendukung:

- **YouTube**: Live streams, videos, shorts
- **TikTok**: Live streams, videos
- **Validasi URL**: Otomatis memverifikasi link yang valid
- **Metadata Extraction**: Mengambil informasi creator dan konten
- **Cross-Platform Donations**: Donasi unified untuk semua platform

## 🔗 **Supported URL Formats**

### **YouTube URLs**
```
# Live Streams
https://www.youtube.com/watch?v=VIDEO_ID
https://youtu.be/VIDEO_ID
https://www.youtube.com/live/VIDEO_ID

# Channel URLs
https://www.youtube.com/channel/CHANNEL_ID
https://www.youtube.com/@username
https://www.youtube.com/c/channelname

# YouTube Shorts
https://www.youtube.com/shorts/VIDEO_ID
```

### **TikTok URLs**
```
# TikTok Videos
https://www.tiktok.com/@username/video/VIDEO_ID
https://vm.tiktok.com/SHORT_CODE

# TikTok Live
https://www.tiktok.com/@username/live

# TikTok Profile
https://www.tiktok.com/@username
```

## 🛠 **API Implementation**

### **1. Database Schema Updates**

Tambahkan tabel untuk menyimpan informasi platform eksternal:

```sql
-- Tabel untuk menyimpan informasi platform streaming
CREATE TABLE streaming_platforms (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    platform_type VARCHAR(20) NOT NULL CHECK (platform_type IN ('youtube', 'tiktok', 'twitch')),
    platform_user_id VARCHAR(255) NOT NULL,
    platform_username VARCHAR(255) NOT NULL,
    channel_url TEXT NOT NULL,
    channel_name VARCHAR(255),
    profile_image_url TEXT,
    follower_count INTEGER DEFAULT 0,
    is_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, platform_type, platform_user_id)
);

-- Tabel untuk menyimpan konten streaming aktif
CREATE TABLE streaming_content (
    id SERIAL PRIMARY KEY,
    platform_id INTEGER REFERENCES streaming_platforms(id) ON DELETE CASCADE,
    content_type VARCHAR(20) NOT NULL CHECK (content_type IN ('live', 'video', 'short')),
    content_id VARCHAR(255) NOT NULL,
    content_url TEXT NOT NULL,
    title VARCHAR(500),
    description TEXT,
    thumbnail_url TEXT,
    duration INTEGER, -- dalam detik, NULL untuk live stream
    view_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    is_live BOOLEAN DEFAULT FALSE,
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(platform_id, content_id)
);

-- Tabel untuk tracking donasi per konten
CREATE TABLE content_donations (
    id SERIAL PRIMARY KEY,
    donation_id INTEGER REFERENCES donations(id) ON DELETE CASCADE,
    content_id INTEGER REFERENCES streaming_content(id) ON DELETE SET NULL,
    platform_type VARCHAR(20) NOT NULL,
    content_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk performa
CREATE INDEX idx_streaming_platforms_user_platform ON streaming_platforms(user_id, platform_type);
CREATE INDEX idx_streaming_content_platform_live ON streaming_content(platform_id, is_live);
CREATE INDEX idx_content_donations_content ON content_donations(content_id);
CREATE INDEX idx_content_donations_donation ON content_donations(donation_id);
```

### **2. API Endpoints**

#### **Platform Management**

```yaml
# Tambahkan ke swagger.yaml

paths:
  # Platform Integration Endpoints
  /platforms/connect:
    post:
      tags:
        - Platform Integration
      summary: Connect social media platform
      description: Connect YouTube or TikTok account to MediaShar
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ConnectPlatformRequest'
            examples:
              youtube:
                summary: Connect YouTube channel
                value:
                  platform_type: "youtube"
                  channel_url: "https://www.youtube.com/@username"
                  platform_username: "username"
              tiktok:
                summary: Connect TikTok account
                value:
                  platform_type: "tiktok"
                  channel_url: "https://www.tiktok.com/@username"
                  platform_username: "username"
      responses:
        '201':
          description: Platform connected successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/PlatformResponse'

  /platforms:
    get:
      tags:
        - Platform Integration
      summary: Get connected platforms
      description: Get list of connected social media platforms
      responses:
        '200':
          description: Connected platforms retrieved successfully
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: "success"
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Platform'

  /platforms/{id}:
    put:
      tags:
        - Platform Integration
      summary: Update platform connection
      description: Update connected platform information
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdatePlatformRequest'
      responses:
        '200':
          description: Platform updated successfully

    delete:
      tags:
        - Platform Integration
      summary: Disconnect platform
      description: Remove platform connection
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Platform disconnected successfully

  # Content Management
  /content/sync:
    post:
      tags:
        - Content Management
      summary: Sync content from platform
      description: Sync latest content from connected platforms
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                platform_id:
                  type: integer
                  description: Platform ID to sync (optional, syncs all if not provided)
                content_url:
                  type: string
                  description: Specific content URL to sync (optional)
              example:
                platform_id: 1
                content_url: "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
      responses:
        '200':
          description: Content synced successfully
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: "success"
                  data:
                    type: object
                    properties:
                      synced_count:
                        type: integer
                      content:
                        type: array
                        items:
                          $ref: '#/components/schemas/StreamingContent'

  /content/validate:
    post:
      tags:
        - Content Management
      summary: Validate streaming URL
      description: Validate and extract metadata from YouTube/TikTok URL
      security: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - url
              properties:
                url:
                  type: string
                  format: uri
                  description: YouTube or TikTok URL to validate
              example:
                url: "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
      responses:
        '200':
          description: URL validated successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/URLValidationResponse'
        '400':
          description: Invalid URL or unsupported platform

  # Enhanced Donation Endpoints
  /donations/to-content:
    post:
      tags:
        - Donations
      summary: Create donation to streaming content
      description: Create donation directly to YouTube/TikTok content
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ContentDonationRequest'
            examples:
              youtube_live:
                summary: Donation to YouTube live stream
                value:
                  amount: 50000
                  currency: "IDR"
                  message: "Great stream! Keep it up!"
                  content_url: "https://www.youtube.com/watch?v=LIVE_ID"
                  display_name: "Anonymous Supporter"
                  is_anonymous: false
              tiktok_video:
                summary: Donation to TikTok video
                value:
                  amount: 25.50
                  currency: "USD"
                  message: "Amazing content!"
                  content_url: "https://www.tiktok.com/@username/video/123456"
                  display_name: "Fan"
                  is_anonymous: true
      responses:
        '201':
          description: Donation created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ContentDonationResponse'

components:
  schemas:
    # Platform Integration Schemas
    ConnectPlatformRequest:
      type: object
      required:
        - platform_type
        - channel_url
        - platform_username
      properties:
        platform_type:
          type: string
          enum: [youtube, tiktok]
          description: Social media platform type
        channel_url:
          type: string
          format: uri
          description: Channel or profile URL
        platform_username:
          type: string
          description: Username on the platform
        auto_sync:
          type: boolean
          default: true
          description: Automatically sync new content

    Platform:
      type: object
      properties:
        id:
          type: integer
        platform_type:
          type: string
          enum: [youtube, tiktok]
        platform_user_id:
          type: string
        platform_username:
          type: string
        channel_url:
          type: string
        channel_name:
          type: string
        profile_image_url:
          type: string
        follower_count:
          type: integer
        is_verified:
          type: boolean
        is_active:
          type: boolean
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

    StreamingContent:
      type: object
      properties:
        id:
          type: integer
        platform_id:
          type: integer
        content_type:
          type: string
          enum: [live, video, short]
        content_id:
          type: string
        content_url:
          type: string
        title:
          type: string
        description:
          type: string
        thumbnail_url:
          type: string
        duration:
          type: integer
          nullable: true
        view_count:
          type: integer
        like_count:
          type: integer
        is_live:
          type: boolean
        started_at:
          type: string
          format: date-time
          nullable: true
        ended_at:
          type: string
          format: date-time
          nullable: true

    URLValidationResponse:
      type: object
      properties:
        status:
          type: string
          example: "success"
        data:
          type: object
          properties:
            is_valid:
              type: boolean
            platform:
              type: string
              enum: [youtube, tiktok]
            content_type:
              type: string
              enum: [live, video, short, channel]
            metadata:
              type: object
              properties:
                title:
                  type: string
                creator:
                  type: string
                thumbnail:
                  type: string
                duration:
                  type: integer
                  nullable: true
                is_live:
                  type: boolean
                view_count:
                  type: integer

    ContentDonationRequest:
      type: object
      required:
        - amount
        - currency
        - content_url
      properties:
        amount:
          type: number
          format: float
          minimum: 0.01
        currency:
          type: string
          enum: [USD, IDR, EUR, GBP]
        message:
          type: string
          maxLength: 500
        content_url:
          type: string
          format: uri
          description: YouTube or TikTok content URL
        display_name:
          type: string
          maxLength: 100
        is_anonymous:
          type: boolean
          default: false
        payment_method:
          type: string
          enum: [qris, paypal, stripe]
          default: "qris"

    ContentDonationResponse:
      type: object
      properties:
        status:
          type: string
          example: "success"
        message:
          type: string
          example: "Donation created successfully"
        data:
          type: object
          properties:
            donation:
              $ref: '#/components/schemas/Donation'
            content:
              $ref: '#/components/schemas/StreamingContent'
            payment_info:
              type: object
              properties:
                qr_code:
                  type: string
                  description: Base64 QR code (for QRIS)
                transaction_id:
                  type: string
                payment_url:
                  type: string
                  description: Payment URL (for PayPal/Stripe)

tags:
  - name: Platform Integration
    description: Social media platform integration endpoints
  - name: Content Management
    description: Streaming content management endpoints
```

### **3. Go Implementation**

#### **URL Validation Service**

```go
// internal/services/platform_service.go
package services

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "regexp"
    "strings"
    "time"
)

type PlatformService struct {
    httpClient *http.Client
}

type URLValidationResult struct {
    IsValid     bool                   `json:"is_valid"`
    Platform    string                 `json:"platform"`
    ContentType string                 `json:"content_type"`
    Metadata    map[string]interface{} `json:"metadata"`
}

type YouTubeMetadata struct {
    VideoID     string `json:"video_id"`
    Title       string `json:"title"`
    Creator     string `json:"creator"`
    Thumbnail   string `json:"thumbnail"`
    Duration    *int   `json:"duration"`
    IsLive      bool   `json:"is_live"`
    ViewCount   int    `json:"view_count"`
    ChannelID   string `json:"channel_id"`
}

type TikTokMetadata struct {
    VideoID     string `json:"video_id"`
    Username    string `json:"username"`
    Title       string `json:"title"`
    Thumbnail   string `json:"thumbnail"`
    ViewCount   int    `json:"view_count"`
    LikeCount   int    `json:"like_count"`
    IsLive      bool   `json:"is_live"`
}

func NewPlatformService() *PlatformService {
    return &PlatformService{
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (s *PlatformService) ValidateURL(inputURL string) (*URLValidationResult, error) {
    parsedURL, err := url.Parse(inputURL)
    if err != nil {
        return &URLValidationResult{IsValid: false}, nil
    }

    // Detect platform
    platform := s.detectPlatform(parsedURL.Host)
    if platform == "" {
        return &URLValidationResult{IsValid: false}, nil
    }

    switch platform {
    case "youtube":
        return s.validateYouTubeURL(inputURL, parsedURL)
    case "tiktok":
        return s.validateTikTokURL(inputURL, parsedURL)
    default:
        return &URLValidationResult{IsValid: false}, nil
    }
}

func (s *PlatformService) detectPlatform(host string) string {
    host = strings.ToLower(host)
    
    if strings.Contains(host, "youtube.com") || strings.Contains(host, "youtu.be") {
        return "youtube"
    }
    
    if strings.Contains(host, "tiktok.com") || strings.Contains(host, "vm.tiktok.com") {
        return "tiktok"
    }
    
    return ""
}

func (s *PlatformService) validateYouTubeURL(inputURL string, parsedURL *url.URL) (*URLValidationResult, error) {
    // YouTube URL patterns
    patterns := map[string]*regexp.Regexp{
        "video":   regexp.MustCompile(`(?:youtube\.com/watch\?v=|youtu\.be/)([a-zA-Z0-9_-]{11})`),
        "live":    regexp.MustCompile(`youtube\.com/live/([a-zA-Z0-9_-]{11})`),
        "shorts":  regexp.MustCompile(`youtube\.com/shorts/([a-zA-Z0-9_-]{11})`),
        "channel": regexp.MustCompile(`youtube\.com/(?:channel/|c/|@)([a-zA-Z0-9_-]+)`),
    }

    for contentType, pattern := range patterns {
        if matches := pattern.FindStringSubmatch(inputURL); len(matches) > 1 {
            videoID := matches[1]
            
            // Get metadata from YouTube API (simplified)
            metadata, err := s.getYouTubeMetadata(videoID, contentType)
            if err != nil {
                return &URLValidationResult{
                    IsValid:     true,
                    Platform:    "youtube",
                    ContentType: contentType,
                    Metadata:    map[string]interface{}{"video_id": videoID},
                }, nil
            }

            return &URLValidationResult{
                IsValid:     true,
                Platform:    "youtube",
                ContentType: contentType,
                Metadata:    s.structToMap(metadata),
            }, nil
        }
    }

    return &URLValidationResult{IsValid: false}, nil
}

func (s *PlatformService) validateTikTokURL(inputURL string, parsedURL *url.URL) (*URLValidationResult, error) {
    // TikTok URL patterns
    patterns := map[string]*regexp.Regexp{
        "video": regexp.MustCompile(`tiktok\.com/@([^/]+)/video/(\d+)`),
        "live":  regexp.MustCompile(`tiktok\.com/@([^/]+)/live`),
        "profile": regexp.MustCompile(`tiktok\.com/@([^/]+)$`),
    }

    for contentType, pattern := range patterns {
        if matches := pattern.FindStringSubmatch(inputURL); len(matches) > 1 {
            username := matches[1]
            var videoID string
            if len(matches) > 2 {
                videoID = matches[2]
            }

            // Get metadata from TikTok (simplified - in production use official API)
            metadata, err := s.getTikTokMetadata(username, videoID, contentType)
            if err != nil {
                return &URLValidationResult{
                    IsValid:     true,
                    Platform:    "tiktok",
                    ContentType: contentType,
                    Metadata:    map[string]interface{}{"username": username, "video_id": videoID},
                }, nil
            }

            return &URLValidationResult{
                IsValid:     true,
                Platform:    "tiktok",
                ContentType: contentType,
                Metadata:    s.structToMap(metadata),
            }, nil
        }
    }

    return &URLValidationResult{IsValid: false}, nil
}

func (s *PlatformService) getYouTubeMetadata(videoID, contentType string) (*YouTubeMetadata, error) {
    // Simplified implementation - in production, use YouTube Data API v3
    // This would require API key and proper OAuth
    
    // For demo purposes, return mock data
    return &YouTubeMetadata{
        VideoID:   videoID,
        Title:     "Sample YouTube Video",
        Creator:   "Sample Creator",
        Thumbnail: fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoID),
        Duration:  nil, // Will be populated from API
        IsLive:    contentType == "live",
        ViewCount: 1000,
        ChannelID: "UC_sample_channel_id",
    }, nil
}

func (s *PlatformService) getTikTokMetadata(username, videoID, contentType string) (*TikTokMetadata, error) {
    // Simplified implementation - in production, use TikTok API
    // This would require proper API access and authentication
    
    return &TikTokMetadata{
        VideoID:   videoID,
        Username:  username,
        Title:     "Sample TikTok Video",
        Thumbnail: "https://example.com/thumbnail.jpg",
        ViewCount: 5000,
        LikeCount: 500,
        IsLive:    contentType == "live",
    }, nil
}

func (s *PlatformService) structToMap(obj interface{}) map[string]interface{} {
    result := make(map[string]interface{})
    data, _ := json.Marshal(obj)
    json.Unmarshal(data, &result)
    return result
}
```

#### **Platform Handler**

```go
// internal/handlers/platform_handler.go
package handlers

import (
    "net/http"
    "strconv"

    "github.com/labstack/echo/v4"
    "github.com/rzfd/mediashar/internal/models"
    "github.com/rzfd/mediashar/internal/services"
)

type PlatformHandler struct {
    platformService *services.PlatformService
    // Add database service here
}

func NewPlatformHandler(platformService *services.PlatformService) *PlatformHandler {
    return &PlatformHandler{
        platformService: platformService,
    }
}

func (h *PlatformHandler) ValidateURL(c echo.Context) error {
    var req struct {
        URL string `json:"url" validate:"required,url"`
    }

    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "status":  "error",
            "message": "Invalid request body",
        })
    }

    if err := c.Validate(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "status":  "error",
            "message": "URL is required and must be valid",
        })
    }

    result, err := h.platformService.ValidateURL(req.URL)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "status":  "error",
            "message": "Failed to validate URL",
        })
    }

    if !result.IsValid {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "status":  "error",
            "message": "Invalid or unsupported URL",
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "status": "success",
        "data":   result,
    })
}

func (h *PlatformHandler) ConnectPlatform(c echo.Context) error {
    userID := getUserIDFromContext(c) // Implement this helper

    var req struct {
        PlatformType     string `json:"platform_type" validate:"required,oneof=youtube tiktok"`
        ChannelURL       string `json:"channel_url" validate:"required,url"`
        PlatformUsername string `json:"platform_username" validate:"required"`
        AutoSync         bool   `json:"auto_sync"`
    }

    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "status":  "error",
            "message": "Invalid request body",
        })
    }

    if err := c.Validate(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "status":  "error",
            "message": "Validation failed",
        })
    }

    // Validate the URL first
    validation, err := h.platformService.ValidateURL(req.ChannelURL)
    if err != nil || !validation.IsValid {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "status":  "error",
            "message": "Invalid platform URL",
        })
    }

    // Create platform connection in database
    platform := &models.StreamingPlatform{
        UserID:           userID,
        PlatformType:     req.PlatformType,
        PlatformUsername: req.PlatformUsername,
        ChannelURL:       req.ChannelURL,
        IsActive:         true,
    }

    // Extract additional metadata from validation result
    if metadata := validation.Metadata; metadata != nil {
        if channelName, ok := metadata["creator"].(string); ok {
            platform.ChannelName = channelName
        }
        if profileImage, ok := metadata["thumbnail"].(string); ok {
            platform.ProfileImageURL = profileImage
        }
    }

    // Save to database (implement this)
    // err = h.platformRepo.Create(platform)
    // if err != nil {
    //     return c.JSON(http.StatusInternalServerError, ...)
    // }

    return c.JSON(http.StatusCreated, map[string]interface{}{
        "status":  "success",
        "message": "Platform connected successfully",
        "data":    platform,
    })
}

func (h *PlatformHandler) CreateContentDonation(c echo.Context) error {
    var req struct {
        Amount      float64 `json:"amount" validate:"required,min=0.01"`
        Currency    string  `json:"currency" validate:"required,oneof=USD IDR EUR GBP"`
        Message     string  `json:"message" validate:"max=500"`
        ContentURL  string  `json:"content_url" validate:"required,url"`
        DisplayName string  `json:"display_name" validate:"max=100"`
        IsAnonymous bool    `json:"is_anonymous"`
        PaymentMethod string `json:"payment_method" validate:"oneof=qris paypal stripe"`
    }

    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "status":  "error",
            "message": "Invalid request body",
        })
    }

    if err := c.Validate(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "status":  "error",
            "message": "Validation failed",
        })
    }

    // Validate content URL
    validation, err := h.platformService.ValidateURL(req.ContentURL)
    if err != nil || !validation.IsValid {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "status":  "error",
            "message": "Invalid content URL",
        })
    }

    // Find or create streamer based on content URL
    // This would involve:
    // 1. Extract creator info from URL validation
    // 2. Find existing MediaShar user with this platform connection
    // 3. If not found, create a placeholder or return error

    // For now, return success with mock data
    return c.JSON(http.StatusCreated, map[string]interface{}{
        "status":  "success",
        "message": "Donation created successfully",
        "data": map[string]interface{}{
            "donation_id":    123,
            "content_url":    req.ContentURL,
            "platform":       validation.Platform,
            "content_type":   validation.ContentType,
            "amount":         req.Amount,
            "currency":       req.Currency,
            "payment_method": req.PaymentMethod,
        },
    })
}

func getUserIDFromContext(c echo.Context) int {
    // Extract user ID from JWT token
    // This should be implemented based on your auth middleware
    return 1 // Mock user ID
}
```

### **4. Frontend Integration**

#### **URL Input Component**

```javascript
// components/URLValidator.js
import React, { useState, useEffect } from 'react';
import axios from 'axios';

const URLValidator = ({ onValidURL, onInvalidURL }) => {
    const [url, setUrl] = useState('');
    const [isValidating, setIsValidating] = useState(false);
    const [validation, setValidation] = useState(null);
    const [error, setError] = useState('');

    const validateURL = async (inputURL) => {
        if (!inputURL) {
            setValidation(null);
            setError('');
            return;
        }

        setIsValidating(true);
        setError('');

        try {
            const response = await axios.post('/api/content/validate', {
                url: inputURL
            });

            if (response.data.status === 'success') {
                setValidation(response.data.data);
                onValidURL && onValidURL(response.data.data);
            }
        } catch (err) {
            const errorMsg = err.response?.data?.message || 'Invalid URL';
            setError(errorMsg);
            setValidation(null);
            onInvalidURL && onInvalidURL(errorMsg);
        } finally {
            setIsValidating(false);
        }
    };

    useEffect(() => {
        const timeoutId = setTimeout(() => {
            validateURL(url);
        }, 500);

        return () => clearTimeout(timeoutId);
    }, [url]);

    const getPlatformIcon = (platform) => {
        switch (platform) {
            case 'youtube':
                return '🎥';
            case 'tiktok':
                return '🎵';
            default:
                return '📺';
        }
    };

    const getContentTypeLabel = (type) => {
        switch (type) {
            case 'live':
                return 'Live Stream';
            case 'video':
                return 'Video';
            case 'short':
                return 'Short';
            case 'channel':
                return 'Channel';
            default:
                return type;
        }
    };

    return (
        <div className="url-validator">
            <div className="input-group">
                <input
                    type="url"
                    value={url}
                    onChange={(e) => setUrl(e.target.value)}
                    placeholder="Paste YouTube or TikTok URL here..."
                    className={`url-input ${validation ? 'valid' : ''} ${error ? 'error' : ''}`}
                />
                {isValidating && <div className="spinner">⏳</div>}
            </div>

            {error && (
                <div className="error-message">
                    ❌ {error}
                </div>
            )}

            {validation && validation.is_valid && (
                <div className="validation-result">
                    <div className="platform-info">
                        <span className="platform-icon">
                            {getPlatformIcon(validation.platform)}
                        </span>
                        <div className="content-details">
                            <div className="platform-name">
                                {validation.platform.toUpperCase()}
                            </div>
                            <div className="content-type">
                                {getContentTypeLabel(validation.content_type)}
                            </div>
                        </div>
                    </div>

                    {validation.metadata && (
                        <div className="metadata">
                            {validation.metadata.title && (
                                <div className="title">
                                    📝 {validation.metadata.title}
                                </div>
                            )}
                            {validation.metadata.creator && (
                                <div className="creator">
                                    👤 {validation.metadata.creator}
                                </div>
                            )}
                            {validation.metadata.is_live && (
                                <div className="live-indicator">
                                    🔴 LIVE
                                </div>
                            )}
                        </div>
                    )}
                </div>
            )}
        </div>
    );
};

export default URLValidator;
```

#### **Donation Form with URL Support**

```javascript
// components/ContentDonationForm.js
import React, { useState } from 'react';
import URLValidator from './URLValidator';
import axios from 'axios';

const ContentDonationForm = () => {
    const [validatedContent, setValidatedContent] = useState(null);
    const [donationData, setDonationData] = useState({
        amount: '',
        currency: 'IDR',
        message: '',
        display_name: '',
        is_anonymous: false,
        payment_method: 'qris'
    });
    const [isSubmitting, setIsSubmitting] = useState(false);

    const handleValidURL = (validation) => {
        setValidatedContent(validation);
    };

    const handleInvalidURL = () => {
        setValidatedContent(null);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        
        if (!validatedContent) {
            alert('Please enter a valid YouTube or TikTok URL');
            return;
        }

        setIsSubmitting(true);

        try {
            const response = await axios.post('/api/donations/to-content', {
                ...donationData,
                content_url: validatedContent.url
            });

            if (response.data.status === 'success') {
                // Handle successful donation
                console.log('Donation created:', response.data.data);
                
                // Redirect to payment or show QR code
                if (response.data.data.payment_info?.qr_code) {
                    // Show QR code for QRIS payment
                    showQRCode(response.data.data.payment_info.qr_code);
                } else if (response.data.data.payment_info?.payment_url) {
                    // Redirect to payment URL
                    window.location.href = response.data.data.payment_info.payment_url;
                }
            }
        } catch (error) {
            console.error('Donation failed:', error);
            alert('Failed to create donation. Please try again.');
        } finally {
            setIsSubmitting(false);
        }
    };

    const showQRCode = (qrCodeData) => {
        // Implement QR code display modal
        console.log('Show QR code:', qrCodeData);
    };

    return (
        <div className="content-donation-form">
            <h2>Donate to Content Creator</h2>
            
            <form onSubmit={handleSubmit}>
                {/* URL Input */}
                <div className="form-group">
                    <label>Content URL</label>
                    <URLValidator 
                        onValidURL={handleValidURL}
                        onInvalidURL={handleInvalidURL}
                    />
                </div>

                {/* Donation Amount */}
                <div className="form-group">
                    <label>Amount</label>
                    <div className="amount-input">
                        <input
                            type="number"
                            value={donationData.amount}
                            onChange={(e) => setDonationData({
                                ...donationData,
                                amount: e.target.value
                            })}
                            placeholder="Enter amount"
                            min="0.01"
                            step="0.01"
                            required
                        />
                        <select
                            value={donationData.currency}
                            onChange={(e) => setDonationData({
                                ...donationData,
                                currency: e.target.value
                            })}
                        >
                            <option value="IDR">IDR</option>
                            <option value="USD">USD</option>
                            <option value="EUR">EUR</option>
                            <option value="GBP">GBP</option>
                        </select>
                    </div>
                </div>

                {/* Message */}
                <div className="form-group">
                    <label>Message (Optional)</label>
                    <textarea
                        value={donationData.message}
                        onChange={(e) => setDonationData({
                            ...donationData,
                            message: e.target.value
                        })}
                        placeholder="Leave a message for the creator..."
                        maxLength="500"
                    />
                </div>

                {/* Display Name */}
                <div className="form-group">
                    <label>Display Name</label>
                    <input
                        type="text"
                        value={donationData.display_name}
                        onChange={(e) => setDonationData({
                            ...donationData,
                            display_name: e.target.value
                        })}
                        placeholder="How should we display your name?"
                        maxLength="100"
                    />
                </div>

                {/* Anonymous Option */}
                <div className="form-group">
                    <label className="checkbox-label">
                        <input
                            type="checkbox"
                            checked={donationData.is_anonymous}
                            onChange={(e) => setDonationData({
                                ...donationData,
                                is_anonymous: e.target.checked
                            })}
                        />
                        Donate anonymously
                    </label>
                </div>

                {/* Payment Method */}
                <div className="form-group">
                    <label>Payment Method</label>
                    <select
                        value={donationData.payment_method}
                        onChange={(e) => setDonationData({
                            ...donationData,
                            payment_method: e.target.value
                        })}
                    >
                        <option value="qris">QRIS</option>
                        <option value="paypal">PayPal</option>
                        <option value="stripe">Credit Card</option>
                    </select>
                </div>

                {/* Submit Button */}
                <button
                    type="submit"
                    disabled={!validatedContent || isSubmitting}
                    className="donate-button"
                >
                    {isSubmitting ? 'Processing...' : 'Donate Now'}
                </button>
            </form>
        </div>
    );
};

export default ContentDonationForm;
```

## 🔧 **Implementation Steps**

### **Phase 1: Basic URL Validation**
1. ✅ Create URL validation service
2. ✅ Add validation API endpoint
3. ✅ Implement basic YouTube/TikTok URL parsing
4. ✅ Create frontend URL input component

### **Phase 2: Platform Integration**
1. 🔄 Add database schema for platforms
2. 🔄 Implement platform connection API
3. 🔄 Add OAuth integration for YouTube/TikTok
4. 🔄 Create platform management dashboard

### **Phase 3: Content Synchronization**
1. ⏳ Implement YouTube Data API integration
2. ⏳ Add TikTok API integration
3. ⏳ Create content sync service
4. ⏳ Add real-time content updates

### **Phase 4: Enhanced Donations**
1. ⏳ Implement content-based donations
2. ⏳ Add creator discovery system
3. ⏳ Create donation analytics
4. ⏳ Add notification system

## 🚀 **Quick Start**

### **1. Update Database**
```bash
# Add migration for platform tables
psql -U postgres -d donation_system -f migrations/add_platform_tables.sql
```

### **2. Update Swagger Documentation**
```bash
# Add new endpoints to swagger.yaml
vim docs/swagger.yaml

# Restart Swagger UI
make swagger-restart
```

### **3. Test URL Validation**
```bash
# Test YouTube URL
curl -X POST http://localhost:8080/api/content/validate \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}'

# Test TikTok URL
curl -X POST http://localhost:8080/api/content/validate \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.tiktok.com/@username/video/123456"}'
```

## 📊 **Benefits**

### **For Content Creators**
- ✅ **Multi-Platform Support**: Accept donations from any platform
- ✅ **Easy Integration**: Just share your content URL
- ✅ **Unified Dashboard**: Manage all donations in one place
- ✅ **Real-time Notifications**: Get notified of new donations

### **For Donators**
- ✅ **Simple Process**: Paste URL and donate
- ✅ **Content Context**: See what you're supporting
- ✅ **Multiple Payment Options**: QRIS, PayPal, Credit Card
- ✅ **Anonymous Options**: Donate privately if preferred

### **For Platform**
- ✅ **Wider Reach**: Support creators on any platform
- ✅ **Increased Engagement**: More donation opportunities
- ✅ **Data Insights**: Analytics across platforms
- ✅ **Competitive Advantage**: Unique cross-platform feature

---

**🎉 Ready to integrate with YouTube and TikTok!**

Start with URL validation and gradually add more features as needed. The modular design allows for easy expansion to other platforms in the future. 