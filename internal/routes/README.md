# Routes Documentation

Folder ini berisi konfigurasi routing untuk aplikasi donation system dengan JWT authentication.

## 📁 **Struktur File Routes**

Routes diorganisir berdasarkan handler untuk maintainability yang lebih baik:

```
internal/routes/
├── routes.go           # Main routes setup (orchestrator)
├── auth_routes.go      # Authentication routes
├── user_routes.go      # User management routes  
├── donation_routes.go  # Donation management routes
├── qris_routes.go      # QRIS payment routes
├── webhook_routes.go   # Payment webhook routes
└── README.md          # Documentation
```

## 🏗️ **Arsitektur Routes**

### **Main Routes File (`routes.go`)**
File utama yang mengatur semua routes dengan memanggil setup function dari setiap handler:

```go
func SetupRoutes(e *echo.Echo, handlers..., jwtSecret string) {
    api := e.Group("/api")
    
    // Setup routes by handler
    SetupAuthRoutes(api, authHandler, jwtSecret)
    SetupUserRoutes(api, userHandler, jwtSecret)
    SetupDonationRoutes(api, donationHandler, jwtSecret)
    SetupQRISRoutes(api, qrisHandler, jwtSecret)
    SetupWebhookRoutes(api, webhookHandler, qrisHandler)
}
```

## 📋 **Detail Routes per Handler**

### **1. Authentication Routes (`auth_routes.go`)**

**File:** `internal/routes/auth_routes.go`

**Public Routes:**
- `POST /api/auth/register` - Registrasi user baru
- `POST /api/auth/login` - Login dan mendapatkan JWT token
- `POST /api/auth/refresh` - Refresh JWT token

**Protected Routes (JWT Required):**
- `GET /api/auth/profile` - Mendapatkan profile user saat ini
- `PUT /api/auth/profile` - Update profile user saat ini
- `POST /api/auth/change-password` - Ganti password
- `POST /api/auth/logout` - Logout

### **2. User Routes (`user_routes.go`)**

**File:** `internal/routes/user_routes.go`

**Public Routes:**
- `GET /api/users/:id` - Mendapatkan detail user
- `GET /api/streamers` - Mendapatkan daftar streamer

**Protected Routes (JWT Required):**
- `POST /api/users` - Membuat user baru (admin)
- `PUT /api/users/:id` - Update informasi user (self atau admin)
- `GET /api/users/:id/donations` - Mendapatkan donasi yang dibuat oleh user

### **3. Donation Routes (`donation_routes.go`)**

**File:** `internal/routes/donation_routes.go`

**Protected Routes (JWT Required):**
- `POST /api/donations` - Membuat donasi baru
- `GET /api/donations` - Mendapatkan daftar semua donasi
- `GET /api/donations/:id` - Mendapatkan detail donasi
- `GET /api/donations/latest` - Mendapatkan donasi terbaru

**Payment Processing:**
- `POST /api/payments/process` - Memproses pembayaran donasi

**Streamer-Only Routes (JWT + Streamer Role):**
- `GET /api/streamers/:id/donations` - Mendapatkan donasi untuk streamer tertentu
- `GET /api/streamers/:id/total` - Mendapatkan total donasi streamer

### **4. QRIS Routes (`qris_routes.go`)**

**File:** `internal/routes/qris_routes.go`

**Public Routes (Optional JWT):**
- `POST /api/qris/donate` - Membuat donasi dengan QRIS (anonymous atau authenticated)

**Protected Routes (JWT Required):**
- `POST /api/qris/donations/:id/generate` - Generate QRIS untuk donasi existing
- `GET /api/qris/status/:transaction_id` - Cek status pembayaran QRIS

### **5. Webhook Routes (`webhook_routes.go`)**

**File:** `internal/routes/webhook_routes.go`

**Public Routes (Secured by webhook secrets):**
- `POST /api/webhooks/paypal` - Webhook untuk PayPal
- `POST /api/webhooks/stripe` - Webhook untuk Stripe
- `POST /api/webhooks/crypto` - Webhook untuk cryptocurrency
- `POST /api/webhooks/qris` - Webhook untuk QRIS payment

## 🔧 **Middleware Usage**

### **JWT Middleware Types**
- **JWTMiddleware**: Memvalidasi JWT token dan mengekstrak informasi user
- **OptionalJWTMiddleware**: Validasi JWT opsional (tidak gagal jika tidak ada token)
- **StreamerOnlyMiddleware**: Memastikan hanya streamer yang bisa mengakses endpoint tertentu

### **Contoh Penggunaan Middleware**
```go
// Protected routes dengan JWT
protectedAuth := api.Group("/auth", middleware.JWTMiddleware(jwtSecret))

// Optional JWT untuk anonymous donations
qrisPublic := api.Group("", middleware.OptionalJWTMiddleware(jwtSecret))

// Streamer-only routes
streamerDonations := api.Group("/streamers", 
    middleware.JWTMiddleware(jwtSecret), 
    middleware.StreamerOnlyMiddleware())
```

## 🏗️ **Keuntungan Struktur Modular**

### **1. Maintainability**
- ✅ Setiap handler memiliki file routes terpisah
- ✅ Mudah menemukan dan mengubah routes spesifik
- ✅ Mengurangi konflik saat multiple developer bekerja

### **2. Scalability**
- ✅ Mudah menambah routes baru per handler
- ✅ Tidak perlu mengubah file routes utama
- ✅ Struktur yang konsisten untuk semua handler

### **3. Readability**
- ✅ Kode lebih terorganisir dan mudah dibaca
- ✅ Separation of concerns yang jelas
- ✅ Dokumentasi yang lebih fokus per handler

## 📝 **Cara Menambah Routes Baru**

### **1. Untuk Handler Existing**
Edit file routes yang sesuai, contoh untuk auth routes:

```go
// internal/routes/auth_routes.go
func SetupAuthRoutes(api *echo.Group, authHandler *handler.AuthHandler, jwtSecret string) {
    // ... existing routes ...
    
    // Tambah route baru
    protectedAuth.POST("/verify-email", authHandler.VerifyEmail)
}
```

### **2. Untuk Handler Baru**
1. Buat file routes baru: `internal/routes/new_handler_routes.go`
2. Implementasi setup function:
```go
func SetupNewHandlerRoutes(api *echo.Group, newHandler *handler.NewHandler, jwtSecret string) {
    // Setup routes untuk handler baru
}
```
3. Tambahkan ke main routes file:
```go
// internal/routes/routes.go
func SetupRoutes(e *echo.Echo, ..., newHandler *handler.NewHandler, jwtSecret string) {
    api := e.Group("/api")
    
    // ... existing setup calls ...
    SetupNewHandlerRoutes(api, newHandler, jwtSecret)
}
```

## Authentication Flow

1. **Register/Login**: User mendaftar atau login melalui `/api/auth/register` atau `/api/auth/login`
2. **Receive Token**: Server mengembalikan JWT token
3. **Use Token**: Client menyertakan token di header `Authorization: Bearer <token>`
4. **Access Protected Routes**: Server memvalidasi token dan memberikan akses

## Security Features

- JWT token dengan expiry time
- Password hashing dengan bcrypt
- Role-based access control (Streamer vs Donator)
- CORS protection
- Input validation
- SQL injection protection via GORM 