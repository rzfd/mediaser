# Midtrans Payment Integration Guide

## Overview
Integrasi Midtrans memungkinkan aplikasi MediaShar untuk menerima pembayaran donation menggunakan berbagai metode pembayaran yang didukung oleh Midtrans.

## Setup Configuration

### 1. Environment Variables
Tambahkan konfigurasi berikut ke file `.env`:

```bash
# Midtrans Configuration
MIDTRANS_MERCHANT_ID=G454372620
MIDTRANS_CLIENT_KEY=SB-Mid-client-Yy6kDu1A1cTYWiYy
MIDTRANS_SERVER_KEY=SB-Mid-server-Zz8uCQ5-zrUcEbes_eijanu
MIDTRANS_ENVIRONMENT=sandbox
```

### 2. YAML Configuration
Atau konfigurasi melalui `configs/config.yaml`:

```yaml
midtrans:
  merchantID: "G454372620"
  clientKey: "SB-Mid-client-Yy6kDu1A1cTYWiYy"
  serverKey: "SB-Mid-server-Zz8uCQ5-zrUcEbes_eijanu"
  environment: "sandbox"  # sandbox atau production
  webhookSecret: ""  # Opsional untuk webhook validation
```

## API Endpoints

### 1. Create Payment
**POST** `/api/midtrans/payment/:donationId`

Headers:
```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

Response:
```json
{
  "status": "success",
  "data": {
    "token": "snap-token-here",
    "redirect_url": "https://app.sandbox.midtrans.com/snap/v2/vtweb/token",
    "order_id": "DONATION-1-1642751234"
  }
}
```

### 2. Webhook Handler
**POST** `/api/midtrans/webhook`

Endpoint ini digunakan oleh Midtrans untuk mengirim notifikasi status pembayaran.

### 3. Get Transaction Status
**GET** `/api/midtrans/status/:orderId`

Response:
```json
{
  "status": "success",
  "data": {
    "transaction_status": "settlement",
    "status_code": "200",
    "transaction_id": "12345",
    "order_id": "DONATION-1-1642751234",
    "gross_amount": "50000.00",
    "payment_type": "bank_transfer"
  }
}
```

## Frontend Integration

### 1. Membuat Payment
```javascript
// 1. Buat donation terlebih dahulu
const donation = await fetch('/api/donations', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    amount: 50000,
    currency: "IDR",
    message: "Support untuk streaming!",
    streamer_id: 1,
    display_name: "Anonymous Donator",
    is_anonymous: false
  })
});

const donationData = await donation.json();

// 2. Buat Midtrans payment
const payment = await fetch(`/api/midtrans/payment/${donationData.data.id}`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const paymentData = await payment.json();

// 3. Redirect ke Midtrans Snap
window.open(paymentData.data.redirect_url, '_blank');
```

### 2. Menggunakan Snap.js (Recommended)
```html
<!-- Include Snap.js -->
<script src="https://app.sandbox.midtrans.com/snap/snap.js" data-client-key="SB-Mid-client-Yy6kDu1A1cTYWiYy"></script>

<script>
async function processPayment(donationId) {
  try {
    // Get snap token
    const response = await fetch(`/api/midtrans/payment/${donationId}`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    const data = await response.json();
    
    // Open Snap payment popup
    snap.pay(data.data.token, {
      onSuccess: function(result) {
        console.log('Payment success:', result);
        // Redirect to success page
        window.location.href = '/donation/success';
      },
      onPending: function(result) {
        console.log('Payment pending:', result);
        // Show pending message
      },
      onError: function(result) {
        console.log('Payment error:', result);
        // Show error message
      },
      onClose: function() {
        console.log('Payment popup closed');
        // Handle popup close
      }
    });
  } catch (error) {
    console.error('Error creating payment:', error);
  }
}
</script>
```

## Webhook Configuration

### 1. Midtrans Dashboard Setup
1. Login ke Midtrans Dashboard
2. Pergi ke **Settings > Configuration > Notification URL**
3. Tambahkan URL: `https://yourdomain.com/api/midtrans/webhook`

### 2. Webhook Security
Webhook menggunakan signature verification untuk memastikan request berasal dari Midtrans. Server key digunakan untuk verifikasi signature.

## Payment Flow

1. **User membuat donation** → POST `/api/donations`
2. **User memilih Midtrans payment** → POST `/api/midtrans/payment/:donationId`
3. **User diarahkan ke Midtrans Snap** → Snap popup/redirect
4. **User melakukan pembayaran** → Di platform Midtrans
5. **Midtrans mengirim webhook** → POST `/api/midtrans/webhook`
6. **Server update status donation** → Status berubah ke 'completed'

## Error Handling

### Common Errors
- **Invalid donation ID**: Donation tidak ditemukan
- **Invalid signature**: Webhook signature tidak valid
- **Transaction not found**: Order ID tidak ditemukan
- **Midtrans API error**: Error dari Midtrans API

### Status Mapping
- `settlement`, `capture` → `completed`
- `pending` → `pending`
- `deny`, `expire`, `cancel` → `failed`

## Testing

### 1. Test Cards (Sandbox)
```
Card Number: 4811 1111 1111 1114
Expiry: 01/25
CVV: 123
```

### 2. Test Bank Transfer
Gunakan virtual account number yang diberikan Midtrans untuk testing.

### 3. Test E-wallet
Gunakan akun test yang disediakan Midtrans untuk masing-masing e-wallet.

## Production Deployment

### 1. Ganti Environment
```bash
MIDTRANS_ENVIRONMENT=production
```

### 2. Ganti API Keys
Gunakan production keys dari Midtrans Dashboard:
- Production Server Key
- Production Client Key

### 3. Update Webhook URL
Pastikan webhook URL mengarah ke production server.

## Monitoring & Logging

Server akan mencatat semua aktivitas payment di log aplikasi. Monitor webhook notifications untuk memastikan semua pembayaran berhasil diproses.

## Support

Untuk troubleshooting lebih lanjut, hubungi:
- Midtrans Support: https://support.midtrans.com
- Midtrans Documentation: https://docs.midtrans.com 