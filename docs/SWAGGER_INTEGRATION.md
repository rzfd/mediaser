# 📚 Swagger Integration Guide - MediaShar API

Panduan lengkap untuk menggunakan Swagger UI yang telah diintegrasikan dengan Docker Compose, berdasarkan konversi dari Postman collection.

## 🎯 **Overview**

Dokumentasi API MediaShar sekarang tersedia dalam format OpenAPI 3.0 dan dapat diakses melalui Swagger UI yang berjalan di Docker container. Dokumentasi ini dikonversi dari Postman collection yang sudah ada dengan semua endpoint, authentication, dan testing scenarios.

## 🚀 **Quick Start**

### **1. Start Services**
```bash
# Start semua services termasuk Swagger UI
make docker-up

# Atau manual
docker-compose up -d
```

### **2. Access Swagger UI**
Buka browser dan akses:
```
http://localhost:8081
```

### **3. Services Overview**
| Service | URL | Description |
|---------|-----|-------------|
| **API Server** | http://localhost:8080 | MediaShar API |
| **Swagger UI** | http://localhost:8081 | API Documentation |
| **PgAdmin** | http://localhost:8082 | Database Admin |
| **PostgreSQL** | localhost:5432 | Database |

## 📖 **Using Swagger UI**

### **🔐 Authentication Setup**

1. **Register User** (Optional)
   - Expand `Authentication` section
   - Try `POST /auth/register`
   - Use example payload untuk donator atau streamer

2. **Login to Get JWT Token**
   - Try `POST /auth/login`
   - Copy JWT token dari response

3. **Set Authorization**
   - Klik tombol **"Authorize"** di kanan atas
   - Paste JWT token (tanpa "Bearer ")
   - Klik **"Authorize"**

### **🧪 Testing Endpoints**

#### **Basic Flow Testing**
```
1. POST /auth/register (Streamer) → Get streamer_id
2. POST /auth/register (Donator) → Get JWT token
3. POST /auth/login → Update JWT token
4. POST /donations → Create donation
5. GET /donations → List donations
```

#### **QRIS Flow Testing**
```
1. POST /qris/donate → Create QRIS donation
2. GET /qris/status/{transaction_id} → Check status
```

#### **Streamer Dashboard Testing**
```
1. Login as streamer
2. GET /streamers/{id}/donations → View donations
3. GET /streamers/{id}/total → Check total
```

### **📝 Request Examples**

Setiap endpoint dilengkapi dengan:
- ✅ **Multiple Examples**: Donator vs Streamer scenarios
- ✅ **Request Schemas**: Validation rules dan field descriptions
- ✅ **Response Examples**: Success dan error responses
- ✅ **Authentication Info**: Required permissions

## 🔄 **Conversion from Postman**

### **What Was Converted**

| Postman Feature | Swagger Equivalent | Status |
|-----------------|-------------------|--------|
| **Collections** | Tags (Authentication, Donations, etc.) | ✅ Converted |
| **Requests** | Paths with HTTP methods | ✅ Converted |
| **Examples** | Request/Response examples | ✅ Converted |
| **Environment Variables** | Server configurations | ✅ Converted |
| **Test Scripts** | Response schemas | ✅ Converted |
| **Authentication** | Security schemes (Bearer JWT) | ✅ Converted |

### **Enhanced Features in Swagger**

1. **Interactive Testing**: Try endpoints directly dari browser
2. **Schema Validation**: Real-time request validation
3. **Auto-completion**: Smart field suggestions
4. **Response Visualization**: Formatted JSON responses
5. **Download Options**: Export as JSON/YAML

## 🛠 **Configuration**

### **Swagger UI Environment Variables**

File: `docker-compose.yml`
```yaml
swagger-ui:
  environment:
    - SWAGGER_JSON=/app/swagger.yaml
    - BASE_URL=/
    - DEEP_LINKING=true
    - DISPLAY_OPERATION_ID=true
    - DEFAULT_MODELS_EXPAND_DEPTH=1
    - DEFAULT_MODEL_EXPAND_DEPTH=1
    - DISPLAY_REQUEST_DURATION=true
    - DOC_EXPANSION=list
    - FILTER=true
    - SHOW_EXTENSIONS=true
    - SHOW_COMMON_EXTENSIONS=true
    - TRY_IT_OUT_ENABLED=true
```

### **OpenAPI Specification**

File: `docs/swagger.yaml`
- **Format**: OpenAPI 3.0.3
- **Authentication**: JWT Bearer token
- **Servers**: Development (localhost:8080) + Production ready
- **Tags**: Organized by functionality
- **Schemas**: Complete request/response models

## 📊 **API Endpoints Overview**

### **🔐 Authentication (8 endpoints)**
- `POST /auth/register` - Register user
- `POST /auth/login` - User login
- `GET /auth/profile` - Get profile
- `PUT /auth/profile` - Update profile
- `POST /auth/change-password` - Change password
- `POST /auth/refresh` - Refresh token
- `POST /auth/logout` - Logout

### **👥 User Management (4 endpoints)**
- `GET /users/{id}` - Get user by ID
- `GET /streamers` - List streamers
- `PUT /users/{id}` - Update user
- `GET /users/{id}/donations` - Get user donations

### **💰 Donations (5 endpoints)**
- `POST /donations` - Create donation
- `GET /donations` - List donations
- `GET /donations/{id}` - Get donation by ID
- `GET /donations/latest` - Get latest donations
- `POST /payments/process` - Process payment

### **🎮 Streamer Endpoints (2 endpoints)**
- `GET /streamers/{id}/donations` - Get streamer donations
- `GET /streamers/{id}/total` - Get total donations

### **💳 QRIS Payments (3 endpoints)**
- `POST /qris/donate` - Create QRIS donation
- `POST /qris/donations/{id}/generate` - Generate QRIS
- `GET /qris/status/{transaction_id}` - Check status

### **🔗 Webhooks (3 endpoints)**
- `POST /webhooks/paypal` - PayPal webhook
- `POST /webhooks/stripe` - Stripe webhook
- `POST /webhooks/qris` - QRIS webhook

## 🔧 **Development Workflow**

### **1. Update API Documentation**

Ketika menambah endpoint baru:

```bash
# 1. Update swagger.yaml
vim docs/swagger.yaml

# 2. Restart Swagger UI
docker-compose restart swagger-ui

# 3. Refresh browser
# http://localhost:8081
```

### **2. Sync with Postman**

Untuk keep Postman dan Swagger in sync:

```bash
# Export dari Swagger ke Postman
# 1. Download OpenAPI spec dari Swagger UI
# 2. Import ke Postman sebagai collection baru
# 3. Update environment variables
```

### **3. Testing Workflow**

```bash
# 1. Test di Swagger UI untuk quick validation
# 2. Export test cases ke Postman untuk automation
# 3. Run Postman collection untuk CI/CD
```

## 🎨 **Customization**

### **Custom Swagger UI Theme**

Untuk customize appearance:

```yaml
# docker-compose.yml
swagger-ui:
  environment:
    - SWAGGER_JSON=/app/swagger.yaml
    - CUSTOM_CSS_URL=/custom.css  # Add custom CSS
    - CUSTOM_JS_URL=/custom.js    # Add custom JS
```

### **Multiple API Versions**

Untuk support multiple versions:

```yaml
volumes:
  - ./docs/swagger-v1.yaml:/app/v1/swagger.yaml:ro
  - ./docs/swagger-v2.yaml:/app/v2/swagger.yaml:ro
```

## 🐛 **Troubleshooting**

### **❌ Swagger UI Not Loading**

```bash
# Check container status
docker-compose ps swagger-ui

# Check logs
docker-compose logs swagger-ui

# Restart service
docker-compose restart swagger-ui
```

### **❌ YAML Syntax Error**

```bash
# Validate YAML syntax
docker run --rm -v $(pwd)/docs:/docs mikefarah/yq eval docs/swagger.yaml

# Or online validator
# https://editor.swagger.io/
```

### **❌ Authentication Not Working**

1. **Check JWT Token Format**
   - Token harus valid JWT format
   - Tidak perlu prefix "Bearer "
   - Check expiration time

2. **Check API Server**
   ```bash
   # Test API directly
   curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/api/auth/profile
   ```

### **❌ CORS Issues**

Jika testing dari Swagger UI gagal karena CORS:

```go
// Add to Go server
func setupCORS(e *echo.Echo) {
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"http://localhost:8081"}, // Swagger UI
        AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
        AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
    }))
}
```

## 📈 **Advanced Features**

### **1. API Versioning**

```yaml
# swagger.yaml
servers:
  - url: http://localhost:8080/api/v1
    description: API Version 1
  - url: http://localhost:8080/api/v2
    description: API Version 2
```

### **2. Environment Switching**

```yaml
# Multiple environments
servers:
  - url: http://localhost:8080/api
    description: Development
  - url: https://staging-api.mediashar.com/api
    description: Staging
  - url: https://api.mediashar.com/api
    description: Production
```

### **3. Code Generation**

Generate client SDKs:

```bash
# Install swagger-codegen
npm install -g swagger-codegen-cli

# Generate JavaScript client
swagger-codegen generate -i docs/swagger.yaml -l javascript -o clients/js

# Generate Go client
swagger-codegen generate -i docs/swagger.yaml -l go -o clients/go
```

## 📚 **Resources**

### **Documentation**
- [OpenAPI 3.0 Specification](https://swagger.io/specification/)
- [Swagger UI Configuration](https://swagger.io/docs/open-source-tools/swagger-ui/usage/configuration/)
- [Swagger Editor](https://editor.swagger.io/)

### **Tools**
- **Swagger Editor**: Online YAML editor
- **Swagger Codegen**: Generate client SDKs
- **Swagger Inspector**: API testing tool
- **Postman**: Import/Export OpenAPI specs

### **Best Practices**
- Keep examples up-to-date dengan actual API
- Use meaningful descriptions untuk semua endpoints
- Include error responses untuk better debugging
- Version your API documentation
- Test documentation dengan real API calls

## 🎉 **Benefits**

### **For Developers**
- ✅ **Interactive Testing**: Test API langsung dari browser
- ✅ **Auto-completion**: Smart field suggestions
- ✅ **Schema Validation**: Real-time validation
- ✅ **Code Generation**: Generate client SDKs

### **For QA/Testers**
- ✅ **Visual Interface**: Easy-to-use testing interface
- ✅ **Request Examples**: Pre-filled request templates
- ✅ **Response Validation**: Automatic response checking
- ✅ **Authentication**: Built-in JWT token management

### **For Product/Business**
- ✅ **API Discovery**: Browse available endpoints
- ✅ **Documentation**: Always up-to-date docs
- ✅ **Integration Planning**: Understand API capabilities
- ✅ **Collaboration**: Share API specs easily

---

**🚀 Ready to explore your API with Swagger UI!**

Access: **http://localhost:8081** 