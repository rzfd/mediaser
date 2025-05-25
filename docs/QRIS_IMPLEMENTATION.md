# QRIS Implementation Guide

## 🏗️ **Arsitektur QRIS System**

Sistem QRIS (Quick Response Code Indonesian Standard) terintegrasi dengan backend donation system untuk memungkinkan pembayaran menggunakan QR Code.

### **Backend Implementation (Recommended)**

✅ **Keuntungan Backend Generate QRIS:**
- Lebih aman (merchant credentials di backend)
- Kontrol penuh atas payment flow
- Logging dan monitoring terpusat
- Validasi server-side
- Konsisten dengan arsitektur existing

## 📋 **API Endpoints**

### **1. Create Donation with QRIS (Public/Optional Auth)**
```http
POST /api/qris/donate
```

**Request Body:**
```json
{
  "amount": 50000,
  "currency": "IDR",
  "message": "Semangat streaming!",
  "streamer_id": 1,
  "display_name": "Anonymous Supporter",
  "is_anonymous": true
}
```

**Response:**
```json
{
  "success": true,
  "message": "Donation created with QRIS",
  "data": {
    "donation": {
      "id": 123,
      "amount": 50000,
      "currency": "IDR",
      "message": "Semangat streaming!",
      "streamer_id": 1,
      "display_name": "Anonymous Supporter",
      "is_anonymous": true,
      "status": "pending",
      "created_at": "2024-01-15T10:30:00Z"
    },
    "qris": {
      "qris_string": "00020101021126280009ID.LINKAJA.WWW...",
      "qr_code_base64": "iVBORw0KGgoAAAANSUhEUgAAAQAAAAEA...",
      "expiry_time": "2024-01-15T10:45:00Z",
      "amount": 50000,
      "transaction_id": "DON-123-1705312200"
    }
  }
}
```

### **2. Generate QRIS for Existing Donation (Protected)**
```http
POST /api/qris/donations/:id/generate
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "success": true,
  "message": "QRIS generated successfully",
  "data": {
    "qris_string": "00020101021126280009ID.LINKAJA.WWW...",
    "qr_code_base64": "iVBORw0KGgoAAAANSUhEUgAAAQAAAAEA...",
    "expiry_time": "2024-01-15T10:45:00Z",
    "amount": 50000,
    "transaction_id": "DON-123-1705312200"
  }
}
```

### **3. Check Payment Status (Protected)**
```http
GET /api/qris/status/:transaction_id
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "success": true,
  "message": "Payment status retrieved",
  "data": {
    "status": "paid",
    "transaction_id": "DON-123-1705312200",
    "amount": 50000,
    "paid_at": "2024-01-15T10:35:00Z"
  }
}
```

### **4. QRIS Payment Webhook (Public)**
```http
POST /api/webhooks/qris
```

## 🔧 **Frontend Implementation Options**

### **Option A: Display QR Code Image (Recommended)**

```javascript
// 1. Create donation with QRIS
const createDonationWithQRIS = async (donationData) => {
  try {
    const response = await fetch('/api/qris/donate', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        // Optional: 'Authorization': `Bearer ${token}` for logged-in users
      },
      body: JSON.stringify(donationData)
    });
    
    const result = await response.json();
    
    if (result.success) {
      // Display QR code image
      displayQRCode(result.data.qris);
      
      // Start polling for payment status
      pollPaymentStatus(result.data.qris.transaction_id);
    }
  } catch (error) {
    console.error('Error creating donation:', error);
  }
};

// 2. Display QR Code
const displayQRCode = (qrisData) => {
  const qrContainer = document.getElementById('qr-container');
  
  qrContainer.innerHTML = `
    <div class="qr-payment">
      <h3>Scan QR Code untuk Donasi</h3>
      <img src="data:image/png;base64,${qrisData.qr_code_base64}" 
           alt="QRIS QR Code" 
           class="qr-code-image" />
      <p>Jumlah: Rp ${qrisData.amount.toLocaleString('id-ID')}</p>
      <p>Berlaku sampai: ${new Date(qrisData.expiry_time).toLocaleString('id-ID')}</p>
      <div class="payment-status">
        <span class="status-indicator">⏳</span>
        <span>Menunggu pembayaran...</span>
      </div>
    </div>
  `;
};

// 3. Poll payment status
const pollPaymentStatus = (transactionId) => {
  const interval = setInterval(async () => {
    try {
      const response = await fetch(`/api/qris/status/${transactionId}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      
      const result = await response.json();
      
      if (result.success) {
        const status = result.data.status;
        
        if (status === 'paid') {
          clearInterval(interval);
          showPaymentSuccess();
        } else if (status === 'expired' || status === 'failed') {
          clearInterval(interval);
          showPaymentFailed();
        }
      }
    } catch (error) {
      console.error('Error checking payment status:', error);
    }
  }, 3000); // Check every 3 seconds
  
  // Stop polling after 15 minutes
  setTimeout(() => clearInterval(interval), 15 * 60 * 1000);
};
```

### **Option B: Generate QR Code di Frontend**

```javascript
// Install: npm install qrcode

import QRCode from 'qrcode';

const generateQRCodeOnFrontend = async (qrisString) => {
  try {
    // Generate QR code as data URL
    const qrCodeDataURL = await QRCode.toDataURL(qrisString, {
      width: 256,
      margin: 2,
      color: {
        dark: '#000000',
        light: '#FFFFFF'
      }
    });
    
    // Display QR code
    const qrImage = document.getElementById('qr-image');
    qrImage.src = qrCodeDataURL;
    
  } catch (error) {
    console.error('Error generating QR code:', error);
  }
};

// Usage
const response = await fetch('/api/qris/donate', { /* ... */ });
const result = await response.json();

if (result.success) {
  // Use QRIS string to generate QR code on frontend
  await generateQRCodeOnFrontend(result.data.qris.qris_string);
}
```

## 🎨 **UI/UX Best Practices**

### **1. QR Code Display Component**

```html
<div class="qris-payment-modal">
  <div class="modal-header">
    <h2>💝 Donasi via QRIS</h2>
    <button class="close-btn">&times;</button>
  </div>
  
  <div class="modal-body">
    <div class="donation-info">
      <div class="streamer-info">
        <img src="streamer-avatar.jpg" alt="Streamer" class="avatar">
        <div>
          <h3>Nama Streamer</h3>
          <p>Gaming Content Creator</p>
        </div>
      </div>
      
      <div class="amount-display">
        <span class="currency">Rp</span>
        <span class="amount">50.000</span>
      </div>
    </div>
    
    <div class="qr-section">
      <div class="qr-container">
        <img id="qr-code" src="" alt="QRIS QR Code" />
      </div>
      
      <div class="instructions">
        <h4>📱 Cara Pembayaran:</h4>
        <ol>
          <li>Buka aplikasi mobile banking atau e-wallet</li>
          <li>Pilih menu "Scan QR" atau "QRIS"</li>
          <li>Arahkan kamera ke QR code di atas</li>
          <li>Konfirmasi pembayaran</li>
        </ol>
      </div>
    </div>
    
    <div class="payment-status">
      <div class="status-indicator loading">
        <div class="spinner"></div>
        <span>Menunggu pembayaran...</span>
      </div>
      
      <div class="timer">
        <span>⏰ Berlaku: </span>
        <span id="countdown">14:59</span>
      </div>
    </div>
  </div>
</div>
```

### **2. CSS Styling**

```css
.qris-payment-modal {
  max-width: 400px;
  margin: 0 auto;
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0,0,0,0.1);
  overflow: hidden;
}

.qr-container {
  text-align: center;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.qr-container img {
  width: 200px;
  height: 200px;
  border-radius: 12px;
  background: white;
  padding: 10px;
}

.amount-display {
  font-size: 2rem;
  font-weight: bold;
  color: #2d3748;
  text-align: center;
  margin: 20px 0;
}

.status-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 15px;
  border-radius: 8px;
  margin: 20px 0;
}

.status-indicator.loading {
  background: #fef5e7;
  color: #d69e2e;
}

.status-indicator.success {
  background: #f0fff4;
  color: #38a169;
}

.status-indicator.failed {
  background: #fed7d7;
  color: #e53e3e;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid #d69e2e;
  border-top: 2px solid transparent;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
```

## 🔒 **Security Considerations**

### **1. QRIS String Validation**
- Validasi format QRIS sesuai standar Bank Indonesia
- Implementasi CRC16-CCITT untuk checksum
- Validasi merchant ID dan amount

### **2. Transaction Security**
- Expiry time untuk QR code (15 menit)
- Unique transaction ID untuk setiap donasi
- Webhook signature validation

### **3. Rate Limiting**
- Limit pembuatan QRIS per user/IP
- Cooldown period untuk donation creation

## 🚀 **Production Deployment**

### **1. Environment Variables**
```bash
# QRIS Configuration
QRIS_MERCHANT_ID=your_merchant_id
QRIS_MERCHANT_NAME="Your Business Name"
QRIS_WEBHOOK_SECRET=your_webhook_secret

# Payment Provider Settings
QRIS_PROVIDER_API_URL=https://api.provider.com
QRIS_PROVIDER_API_KEY=your_api_key
```

### **2. Payment Provider Integration**

Untuk production, integrasikan dengan payment provider seperti:
- **Midtrans**: QRIS payment gateway
- **Xendit**: QRIS payment solution
- **DANA**: Direct QRIS integration
- **GoPay**: QRIS merchant integration
- **OVO**: QRIS payment API

### **3. Webhook Configuration**

Setup webhook endpoint di payment provider:
```
POST https://yourdomain.com/api/webhooks/qris
```

## 📱 **Mobile App Integration**

Untuk mobile app (React Native/Flutter):

```javascript
// React Native example
import { Linking } from 'react-native';

const openQRISApp = (qrisString) => {
  // Deep link ke aplikasi mobile banking
  const deepLinks = [
    `dana://qr?data=${encodeURIComponent(qrisString)}`,
    `gopay://qr?data=${encodeURIComponent(qrisString)}`,
    `ovo://qr?data=${encodeURIComponent(qrisString)}`,
  ];
  
  // Try to open preferred app
  deepLinks.forEach(link => {
    Linking.canOpenURL(link).then(supported => {
      if (supported) {
        Linking.openURL(link);
      }
    });
  });
};
```

## 🧪 **Testing**

### **1. Unit Tests**
```bash
go test ./internal/service -v -run TestQRISService
go test ./internal/handler -v -run TestQRISHandler
```

### **2. Integration Tests**
```bash
# Test QRIS generation
curl -X POST http://localhost:8080/api/qris/donate \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 50000,
    "currency": "IDR",
    "message": "Test donation",
    "streamer_id": 1,
    "display_name": "Test User",
    "is_anonymous": true
  }'

# Test payment status
curl -X GET http://localhost:8080/api/qris/status/DON-123-1705312200 \
  -H "Authorization: Bearer your-jwt-token"
```

## 📊 **Monitoring & Analytics**

Track QRIS payment metrics:
- QR code generation rate
- Payment success rate
- Average payment time
- Popular payment methods
- Failed payment reasons

Implementasi ini memberikan solusi QRIS yang lengkap dan production-ready untuk sistem donasi streamer Anda! 🎯 