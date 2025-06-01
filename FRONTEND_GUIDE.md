# 🌐 MediaShar Frontend Testing Guide

Panduan lengkap untuk menggunakan frontend visual testing interface untuk mengetes integrasi Midtrans dan backend MediaShar.

## 🚀 Quick Start

### 1. Start Development Environment
```bash
# Option 1: Start everything at once
make dev-full

# Option 2: Manual step-by-step
make up              # Start backend services
make frontend        # Start frontend (in another terminal)
```

### 2. Access Frontend
- **Frontend Interface**: http://localhost:8000
- **Backend API**: http://localhost:8080  
- **Swagger Documentation**: http://localhost:8083
- **pgAdmin**: http://localhost:8082

## 📊 Interface Overview

### Header Section
- **MediaShar Logo & Title**: Branding dan identification
- **Connection Status**: Real-time indicator (green=connected, red=error, yellow=checking)
- **Check Health Button**: Manual health check trigger

### Status Dashboard (3 Cards)
1. **API Status**: Backend connection health
2. **Database Status**: PostgreSQL connection status  
3. **Midtrans Status**: Payment gateway configuration

### Main Interface (2 Columns)

#### Left Column: Testing Controls
1. **Authentication Section**
   - Register new users (donator/streamer)
   - Login dengan credentials
   - Quick login buttons untuk test accounts
   - User session display

2. **Donation Section**
   - Create donation dengan amount, currency, message
   - Streamer ID selection
   - Anonymous donation option
   - Display name customization

3. **Payment Section**
   - Generate Midtrans payment token
   - Open Snap payment interface
   - Payment process execution

4. **Quick Actions**
   - Full Flow Test (automated end-to-end)
   - Load Test Data (create test accounts)
   - View Donations (list all donations)
   - Clear Logs (reset interface)

#### Right Column: Monitoring & Results
1. **API Response Log**
   - Real-time API request/response logging
   - Color-coded success/error indicators
   - JSON formatted responses
   - Scrollable history

2. **Session Data**
   - Current JWT token
   - Last donation ID created
   - Last Snap token generated
   - Last order ID from Midtrans

3. **Test Results**
   - Visual test execution results
   - Success/failure indicators
   - Detailed error messages

## 🧪 Testing Workflows

### Automated Full Flow Test (Recommended)
1. **Click "Full Flow Test"** atau tekan `Ctrl+Enter`
2. **Watch the automation**:
   - ✅ Health check
   - ✅ Auto user registration
   - ✅ Auto login
   - ✅ Create test donation (75,000 IDR)
   - ✅ Generate Midtrans payment token
3. **Manual payment test**: Click "Open Snap Payment"
4. **Complete payment**: Use Midtrans test cards

### Manual Step-by-Step Testing

#### Step 1: Prepare Test Data
```bash
# Click "Load Test Data" button
# Creates: streamer@test.com & donator@test.com (password: password123)
```

#### Step 2: Authentication Test
1. **Register new user** (optional):
   - Fill username, email, password
   - Select Donator/Streamer
   - Click "Register User"

2. **Login**:
   - Manual: Enter email & password
   - Quick Login: Click "Quick Login Donator" atau "Quick Login Streamer"

#### Step 3: Create Donation
1. **Set donation details**:
   - Amount: 50,000 (default) atau custom amount
   - Currency: IDR atau USD
   - Streamer ID: 1 (default)
   - Message: Custom donation message
   - Display name: How name appears

2. **Create donation**: Click "Create Donation"
3. **Note donation ID**: Will auto-populate in payment section

#### Step 4: Process Payment
1. **Create payment token**: Click "Create Payment"
2. **Open Snap payment**: Click "Open Snap Payment" (enabled after token creation)
3. **Test payment scenarios**:
   - **Success**: Use card `4811 1111 1111 1114`
   - **Pending**: Use card `4911 1111 1111 1113`
   - **Failed**: Use card `4411 1111 1111 1118`

## 💳 Midtrans Testing Cards

### Test Credit Cards (Sandbox)
```
✅ SUCCESS: 4811 1111 1111 1114
⏳ PENDING: 4911 1111 1111 1113  
❌ FAILED:  4411 1111 1111 1118

CVV: 123
Exp: 12/25
```

### Test Scenarios
1. **Successful Payment**:
   - Card: 4811 1111 1111 1114
   - Expected: Green notification "Payment successful!"
   - Log: Success callback dengan transaction details

2. **Pending Payment**:
   - Card: 4911 1111 1111 1113
   - Expected: Yellow notification "Payment pending"
   - Log: Pending callback dengan transaction status

3. **Failed Payment**:
   - Card: 4411 1111 1111 1118
   - Expected: Red notification "Payment failed"
   - Log: Error callback dengan failure reason

## 📋 Monitoring & Debugging

### API Response Log
- **Green entries**: Successful API calls
- **Red entries**: Failed API calls  
- **Blue entries**: Information/data display
- **Timestamp**: When request was made
- **Request data**: What was sent to API
- **Response data**: What API returned

### Session Data Tracking
- **JWT Token**: Current authentication token (truncated for security)
- **Donation ID**: Last created donation ID
- **Snap Token**: Last generated Midtrans token
- **Order ID**: Last Midtrans order ID

### Real-time Notifications
- **Green**: Success messages
- **Red**: Error messages
- **Yellow**: Warning messages
- **Blue**: Information messages

## ⌨️ Keyboard Shortcuts

| Shortcut | Action | Description |
|----------|---------|-------------|
| `Ctrl+Enter` | Full Flow Test | Run complete automated test |
| `Ctrl+H` | Health Check | Check service health |
| `Ctrl+L` | Clear Logs | Clear API logs dan test results |

## 🔧 Configuration & Customization

### API Endpoints
Edit `frontend/script.js`:
```javascript
const API_BASE_URL = 'http://localhost:8080/api';
const HEALTH_URL = 'http://localhost:8080/health';
```

### Default Values
Customize default form values:
```javascript
// Default donation amount
document.getElementById('donation-amount').value = '100000';

// Default test data
const TEST_CREDENTIALS = {
    streamer: { email: 'streamer@test.com', password: 'password123' },
    donator: { email: 'donator@test.com', password: 'password123' }
};
```

### Midtrans Client Key
Update sandbox credentials di `frontend/index.html`:
```html
<script src="https://app.sandbox.midtrans.com/snap/snap.js" 
        data-client-key="SB-Mid-client-Yy6kDu1A1cTYWiYy"></script>
```

## 🐛 Troubleshooting

### Common Issues & Solutions

#### 1. CORS Error
```
Access to fetch blocked by CORS policy
```
**Solution**: 
- Check backend CORS configuration
- Pastikan frontend serve dari port 8000
- Restart backend services: `make restart`

#### 2. Connection Refused
```
Failed to fetch: Connection refused
```
**Solution**:
- Check backend status: `make status`
- Start services: `make up`
- Check ports: `docker-compose ps`

#### 3. Midtrans Snap Not Loading
```
snap is not defined
```
**Solution**:
- Check internet connection (CDN dependency)
- Check browser console for errors
- Verify client key configuration

#### 4. Authentication Issues
```
Token invalid or expired
```
**Solution**:
- Click "Logout" dan login again
- Clear browser localStorage
- Check JWT token di session data panel

#### 5. Payment Token Generation Failed
```
Failed to create payment
```
**Solution**:
- Ensure donation ID exists
- Check Midtrans configuration di environment
- Verify user authentication
- Check backend logs: `make logs`

### Debug Workflow
1. **Check Browser Console** (F12)
2. **Monitor API Response Log** (real-time di interface)
3. **Check Backend Logs**: `make logs`
4. **Verify Environment**: `make env-check`
5. **Health Check**: `make health-check`

## 📱 Mobile & Responsive Design

### Desktop (Recommended)
- Full feature set
- Side-by-side layout
- Keyboard shortcuts enabled
- Optimal logging display

### Tablet
- Responsive grid layout
- Touch-friendly buttons
- Collapsible sections
- Readable text sizes

### Mobile
- Stacked single-column layout
- Large touch targets
- Simplified navigation
- Essential features only

## 🔒 Security Notes

### Sandbox Environment
- All data adalah test data
- Menggunakan Midtrans sandbox credentials
- JWT tokens stored di browser localStorage
- No production data involved

### Best Practices
- Clear session data after testing
- Don't use production credentials
- Monitor API logs untuk security issues
- Use HTTPS di production environment

## 📊 Performance Tips

### Optimization
- Auto-refresh health status every 30 seconds
- Log rotation (max 10 entries)
- Efficient API connection pooling
- Lazy loading untuk large responses

### Memory Management
- Clear logs regularly with `Ctrl+L`
- Restart frontend periodically untuk long sessions
- Monitor browser memory usage
- Close unused browser tabs

## 🎯 Testing Best Practices

### Comprehensive Testing
1. **Start with automated test**: `Ctrl+Enter`
2. **Test each payment scenario**: Success, Pending, Failed
3. **Test error conditions**: Invalid data, network issues
4. **Test authentication flows**: Register, login, logout
5. **Monitor API responses**: Check for errors or unexpected behavior

### Data Validation
- Test dengan different amounts (small, large, decimal)
- Test dengan different currencies (IDR, USD)
- Test dengan different user types (donator, streamer)
- Test dengan anonymous vs named donations

### Edge Cases
- Empty form submissions
- Invalid email formats
- Very long messages
- Special characters di input fields
- Network disconnection scenarios

---

## 🎉 Happy Testing!

Untuk pertanyaan atau issues:
1. Check browser console untuk detailed errors
2. Monitor API response log di interface
3. Check backend logs dengan `make logs`
4. Review Makefile commands dengan `make help`

**Frontend Interface URL**: http://localhost:8000
**Backend API URL**: http://localhost:8080  
**Swagger Docs**: http://localhost:8083 