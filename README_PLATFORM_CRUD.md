# 🎥 MediaShar Platform Integration - CREATE & READ Operations

Panduan lengkap untuk menggunakan operasi CREATE dan READ pada integrasi platform MediaShar dengan YouTube dan TikTok.

## 📋 **Overview**

Implementasi ini menyediakan operasi dasar CREATE dan READ untuk:

- **✅ StreamingPlatform**: Koneksi platform user (YouTube/TikTok)
- **✅ StreamingContent**: Konten streaming (videos, live streams, shorts)
- **✅ ContentDonation**: Tracking donasi per konten

## 🗄️ **Database Models**

### **StreamingPlatform**
```go
type StreamingPlatform struct {
    Base
    UserID           uint   `json:"user_id"`
    PlatformType     string `json:"platform_type"`     // youtube, tiktok
    PlatformUserID   string `json:"platform_user_id"`
    PlatformUsername string `json:"platform_username"`
    ChannelURL       string `json:"channel_url"`
    ChannelName      string `json:"channel_name"`
    ProfileImageURL  string `json:"profile_image_url"`
    FollowerCount    int    `json:"follower_count"`
    IsVerified       bool   `json:"is_verified"`
    IsActive         bool   `json:"is_active"`
}
```

### **StreamingContent**
```go
type StreamingContent struct {
    Base
    PlatformID   uint      `json:"platform_id"`
    ContentType  string    `json:"content_type"`  // live, video, short
    ContentID    string    `json:"content_id"`
    ContentURL   string    `json:"content_url"`
    Title        string    `json:"title"`
    Description  string    `json:"description"`
    ThumbnailURL string    `json:"thumbnail_url"`
    Duration     *int      `json:"duration"`      // seconds, NULL for live
    ViewCount    int       `json:"view_count"`
    LikeCount    int       `json:"like_count"`
    IsLive       bool      `json:"is_live"`
    StartedAt    *time.Time `json:"started_at"`
    EndedAt      *time.Time `json:"ended_at"`
}
```

### **ContentDonation**
```go
type ContentDonation struct {
    ID           uint      `json:"id"`
    DonationID   uint      `json:"donation_id"`
    ContentID    *uint     `json:"content_id"`    // nullable
    PlatformType string    `json:"platform_type"`
    ContentURL   string    `json:"content_url"`
    CreatedAt    time.Time `json:"created_at"`
}
```

## 🔧 **Repository Operations**

### **StreamingPlatform Operations**

#### **CREATE Platform**
```go
platform := &models.StreamingPlatform{
    UserID:           1,
    PlatformType:     "youtube",
    PlatformUserID:   "UC_channel_123",
    PlatformUsername: "gaming_creator",
    ChannelURL:       "https://www.youtube.com/@gaming_creator",
    ChannelName:      "Gaming Creator Channel",
    IsActive:         true,
}

err := platformRepo.CreatePlatform(platform)
```

#### **READ Platform Operations**
```go
// Get by ID
platform, err := platformRepo.GetPlatformByID(1)

// Get by User ID
platforms, err := platformRepo.GetPlatformsByUserID(1)

// Get by User and Type
platform, err := platformRepo.GetPlatformByUserAndType(1, "youtube")

// Get Active Platforms
platforms, err := platformRepo.GetActivePlatforms(0, 10) // offset, limit
```

### **StreamingContent Operations**

#### **CREATE Content**
```go
content := &models.StreamingContent{
    PlatformID:   1,
    ContentType:  "live",
    ContentID:    "live_stream_123",
    ContentURL:   "https://www.youtube.com/watch?v=live_stream_123",
    Title:        "Epic Gaming Live Stream",
    Description:  "Playing latest games!",
    IsLive:       true,
}

err := platformRepo.CreateContent(content)
```

#### **READ Content Operations**
```go
// Get by ID
content, err := platformRepo.GetContentByID(1)

// Get by Platform ID
contents, err := platformRepo.GetContentByPlatformID(1, 0, 10)

// Get by URL
content, err := platformRepo.GetContentByURL("https://www.youtube.com/watch?v=123")

// Get Live Content
liveContents, err := platformRepo.GetLiveContent(0, 10)

// Get by Content Type
videos, err := platformRepo.GetContentByType("video", 0, 10)
```

### **ContentDonation Operations**

#### **CREATE Content Donation**
```go
contentDonation := &models.ContentDonation{
    DonationID:   1,
    ContentID:    &contentID, // optional
    PlatformType: "youtube",
    ContentURL:   "https://www.youtube.com/watch?v=123",
}

err := platformRepo.CreateContentDonation(contentDonation)
```

#### **READ Content Donation Operations**
```go
// Get by ID
contentDonation, err := platformRepo.GetContentDonationByID(1)

// Get by Donation ID
contentDonations, err := platformRepo.GetContentDonationsByDonationID(1)

// Get by Content ID
contentDonations, err := platformRepo.GetContentDonationsByContentID(1)

// Get by Platform Type
youtubeDonations, err := platformRepo.GetContentDonationsByPlatform("youtube", 0, 10)
```

## 🌐 **API Endpoints**

### **Platform Management**

#### **Connect Platform**
```http
POST /api/platforms/connect
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "platform_type": "youtube",
  "channel_url": "https://www.youtube.com/@gaming_creator",
  "platform_username": "gaming_creator"
}
```

#### **Get Connected Platforms**
```http
GET /api/platforms
Authorization: Bearer JWT_TOKEN
```

#### **Get Platform by ID**
```http
GET /api/platforms/{id}
Authorization: Bearer JWT_TOKEN
```

### **Content Management**

#### **Create Content**
```http
POST /api/content
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "platform_id": 1,
  "content_type": "live",
  "content_id": "live_stream_123",
  "content_url": "https://www.youtube.com/watch?v=live_stream_123",
  "title": "Epic Gaming Live Stream",
  "is_live": true
}
```

#### **Get Content by ID**
```http
GET /api/content/{id}
```

#### **Get Content by URL**
```http
POST /api/content/by-url
Content-Type: application/json

{
  "url": "https://www.youtube.com/watch?v=123"
}
```

#### **Get Live Content**
```http
GET /api/content/live?page=1&pageSize=10
```

### **Content Donations**

#### **Create Content Donation**
```http
POST /api/content-donations
Authorization: Bearer JWT_TOKEN
Content-Type: application/json

{
  "donation_id": 1,
  "content_url": "https://www.youtube.com/watch?v=123",
  "platform_type": "youtube"
}
```

#### **Get Content Donations by Donation**
```http
GET /api/content-donations/donation/{donationId}
Authorization: Bearer JWT_TOKEN
```

## 🧪 **Testing Examples**

### **1. Setup Database Connection**
```go
dsn := "host=localhost user=postgres password=postgres dbname=donation_system port=5432 sslmode=disable"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
platformRepo := repository.NewPlatformRepository(db)
```

### **2. Create YouTube Platform**
```bash
curl -X POST http://localhost:8080/api/platforms/connect \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "platform_type": "youtube",
    "channel_url": "https://www.youtube.com/@gaming_creator",
    "platform_username": "gaming_creator"
  }'
```

### **3. Create TikTok Platform**
```bash
curl -X POST http://localhost:8080/api/platforms/connect \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "platform_type": "tiktok",
    "channel_url": "https://www.tiktok.com/@tiktok_creator",
    "platform_username": "tiktok_creator"
  }'
```

### **4. Get Connected Platforms**
```bash
curl -X GET http://localhost:8080/api/platforms \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### **5. Create Live Content**
```bash
curl -X POST http://localhost:8080/api/content \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "platform_id": 1,
    "content_type": "live",
    "content_id": "live_123",
    "content_url": "https://www.youtube.com/watch?v=live_123",
    "title": "Epic Gaming Stream",
    "is_live": true
  }'
```

### **6. Get Live Content**
```bash
curl -X GET http://localhost:8080/api/content/live?page=1&pageSize=5
```

## 📁 **File Structure**

```
internal/
├── models/
│   └── platform.go              # Platform models
├── repository/
│   └── platform_repository.go   # Repository interface & implementation
├── handler/
│   └── platform_handler_v2.go   # HTTP handlers with repository integration
└── service/
    └── platform_service.go      # URL validation service (existing)

examples/
└── platform_demo.go             # Complete demo (Server + Database tests)

migrations/
└── add_platform_tables.sql     # Database schema
```

## 🚀 **Quick Start**

### **1. Run Migration**
```bash
# Apply database migration
./scripts/setup-platform-integration.sh --migrate
```

### **2. Test Repository Operations**
```go
// Run the example code
go run examples/platform_demo.go
```

### **3. Test API Endpoints**
```bash
# Start the server
make docker-up

# Test URL validation
curl -X POST http://localhost:8080/api/content/validate \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}'

# Test platform connection (requires JWT)
curl -X POST http://localhost:8080/api/platforms/connect \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "platform_type": "youtube",
    "channel_url": "https://www.youtube.com/@test",
    "platform_username": "test"
  }'
```

## 🔍 **Available Queries**

### **Platform Queries**
- ✅ Get platform by ID
- ✅ Get platforms by user ID
- ✅ Get platform by user and type
- ✅ Get active platforms with pagination

### **Content Queries**
- ✅ Get content by ID
- ✅ Get content by platform ID
- ✅ Get content by URL
- ✅ Get live content
- ✅ Get content by type (live/video/short)

### **Content Donation Queries**
- ✅ Get content donation by ID
- ✅ Get content donations by donation ID
- ✅ Get content donations by content ID
- ✅ Get content donations by platform type

## 📊 **Response Examples**

### **Platform Response**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "user_id": 1,
    "platform_type": "youtube",
    "platform_username": "gaming_creator",
    "channel_url": "https://www.youtube.com/@gaming_creator",
    "channel_name": "Gaming Creator Channel",
    "follower_count": 15000,
    "is_verified": true,
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### **Content Response**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "platform_id": 1,
    "content_type": "live",
    "content_url": "https://www.youtube.com/watch?v=live_123",
    "title": "Epic Gaming Live Stream",
    "view_count": 1250,
    "like_count": 89,
    "is_live": true,
    "platform": {
      "channel_name": "Gaming Creator Channel",
      "platform_type": "youtube"
    }
  }
}
```

## ⚠️ **Important Notes**

### **Database Constraints**
- Platform type must be: `youtube`, `tiktok`, or `twitch`
- Content type must be: `live`, `video`, or `short`
- Unique constraint: `(user_id, platform_type, platform_user_id)`
- Unique constraint: `(platform_id, content_id)`

### **Relationships**
- `StreamingPlatform` belongs to `User`
- `StreamingContent` belongs to `StreamingPlatform`
- `ContentDonation` belongs to `Donation` and optionally to `StreamingContent`

### **Preloading**
Repository methods automatically preload related data:
- Platform queries preload `User` and `StreamingContent`
- Content queries preload `Platform` and `ContentDonations`
- Content donation queries preload `Donation` and `Content`

## 🔄 **Next Steps**

### **Immediate**
1. Integrate handlers with main application routes
2. Add proper JWT authentication
3. Add validation middleware
4. Add error handling middleware

### **Future Enhancements**
1. Add UPDATE and DELETE operations
2. Add bulk operations
3. Add search and filtering
4. Add caching layer
5. Add real-time synchronization

## 📚 **Documentation**

- **API Documentation**: http://localhost:8083 (Swagger UI)
- **Database Schema**: `migrations/add_platform_tables.sql`
- **Integration Guide**: `docs/SOCIAL_MEDIA_INTEGRATION.md`
- **Setup Script**: `./scripts/setup-platform-integration.sh`

---

**🎉 Ready to use CREATE & READ operations for platform integration!**

Start testing with the examples above or explore the Swagger UI for interactive API documentation.

**Happy Coding! 🚀** 