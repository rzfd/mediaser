# MediaShar Frontend Testing Interface

Frontend visual testing interface untuk MediaShar donation system dengan integrasi Midtrans payment gateway.

## 🚀 Features

- **Visual Health Monitoring**: Real-time status API, database, dan Midtrans
- **Complete Authentication Flow**: Register, login, logout dengan visual feedback
- **Donation Management**: Create dan manage donations dengan UI yang intuitif
- **Midtrans Payment Integration**: Test Snap payment dengan sandbox credentials
- **Real-time API Logging**: Monitor semua API requests dan responses
- **Session Management**: Track JWT tokens, donation IDs, dan payment data
- **Automated Testing**: Full flow testing dengan satu klik
- **Keyboard Shortcuts**: Quick actions untuk efficient testing

## 🖥️ Interface Components

### Status Dashboard
- **API Status**: Connection status ke backend
- **Database Status**: PostgreSQL connection status
- **Midtrans Status**: Payment gateway configuration status

### Authentication Section
- **Register Form**: Create new users (donator/streamer)
- **Login Form**: Authentication dengan JWT tokens
- **Quick Login**: Pre-configured test accounts
- **Session Display**: Current user information

### Donation Testing
- **Create Donation**: Amount, currency, message, streamer selection
- **Anonymous Options**: Toggle anonymous donations
- **Auto-fill Fields**: Pre-configured test values

### Payment Testing
- **Midtrans Integration**: Direct Snap payment testing
- **Payment Creation**: Generate payment tokens
- **Payment Execution**: Open Midtrans Snap UI
- **Payment Callbacks**: Handle success/pending/error states

### Monitoring & Logs
- **API Response Log**: Real-time request/response logging
- **Session Data**: JWT tokens, donation IDs, payment tokens
- **Test Results**: Visual test execution results
- **Notification System**: Toast notifications untuk user feedback

## 🛠️ How to Use

### 1. Start Backend Services
```bash
# Dari root directory
make up
# atau
docker-compose up -d
```

### 2. Open Frontend
```bash
# Serve frontend (option 1 - simple)
cd frontend
python3 -m http.server 8000

# Serve frontend (option 2 - Node.js)
cd frontend
npx serve -s . -l 8000

# Buka browser
open http://localhost:8000
```

### 3. Testing Workflow

#### Quick Start Testing
1. **Check Health**: Klik "Check Health" atau Ctrl+H
2. **Load Test Data**: Klik "Load Test Data" untuk create test users
3. **Full Flow Test**: Klik "Full Flow Test" atau Ctrl+Enter untuk automated testing
4. **Test Payment**: Klik "Open Snap Payment" untuk test Midtrans

#### Manual Testing
1. **Register User**: Fill form dan create new user
2. **Login**: Use credentials atau quick login buttons
3. **Create Donation**: Set amount, streamer ID, message
4. **Create Payment**: Generate Midtrans payment token
5. **Execute Payment**: Test Snap payment interface

## 🔧 Configuration

### API Configuration
Edit `script.js` untuk mengubah endpoints:
```javascript
const API_BASE_URL = 'http://localhost:8080/api';
const HEALTH_URL = 'http://localhost:8080/health';
```

### Midtrans Configuration
Sandbox credentials di HTML head:
```html
<script src="https://app.sandbox.midtrans.com/snap/snap.js" 
        data-client-key="SB-Mid-client-Yy6kDu1A1cTYWiYy"></script>
```

## ⌨️ Keyboard Shortcuts

- `Ctrl+Enter`: Run full flow test
- `Ctrl+H`: Check health status
- `Ctrl+L`: Clear logs dan test results

## 🧪 Test Scenarios

### Automated Full Flow Test
1. Health check backend services
2. Register new test user (auto-generated)
3. Login dengan test credentials
4. Create test donation (75,000 IDR)
5. Generate Midtrans payment token
6. Ready untuk Snap payment testing

### Manual Test Cases

#### Authentication Tests
- ✅ Register new donator
- ✅ Register new streamer
- ✅ Login dengan valid credentials
- ✅ Login dengan invalid credentials
- ✅ JWT token persistence
- ✅ Session management

#### Donation Tests
- ✅ Create donation dengan required fields
- ✅ Create anonymous donation
- ✅ Different amounts dan currencies
- ✅ Validation error handling
- ✅ Authorization requirements

#### Payment Tests
- ✅ Generate Midtrans payment token
- ✅ Open Snap payment interface
- ✅ Handle payment success callback
- ✅ Handle payment pending callback
- ✅ Handle payment error callback
- ✅ Payment cancellation

## 🎨 UI Features

### Visual Feedback
- **Color-coded status indicators**: Green (success), Red (error), Yellow (warning)
- **Real-time notifications**: Toast messages untuk user actions
- **Progress indicators**: Loading states dan connection status
- **Responsive design**: Mobile-friendly interface

### Data Visualization
- **API logs dengan syntax highlighting**
- **JSON response formatting**
- **Session data display**
- **Test results tracking**

## 🔍 Troubleshooting

### Common Issues

#### CORS Errors
```
Access to fetch at 'http://localhost:8080/api' from origin 'http://localhost:8000' has been blocked by CORS policy
```
**Solution**: Pastikan backend CORS configuration allow frontend origin.

#### Connection Refused
```
Failed to fetch: Connection refused
```
**Solution**: Pastikan backend services running dengan `make status` atau `docker-compose ps`.

#### Midtrans Snap Not Loading
```
snap is not defined
```
**Solution**: Check internet connection untuk load Midtrans CDN atau check browser console untuk errors.

### Debug Tips
1. **Check Browser Console**: Press F12 untuk detailed error messages
2. **Check API Logs**: Monitor API response log panel
3. **Check Network Tab**: Monitor HTTP requests di browser DevTools
4. **Check Backend Logs**: `make logs` atau `docker-compose logs`

## 📱 Mobile Responsive

Interface designed untuk testing di berbagai devices:
- **Desktop**: Full feature set dengan side-by-side layout
- **Tablet**: Responsive grid layout
- **Mobile**: Stacked layout dengan touch-friendly buttons

## 🔒 Security Notes

- **Sandbox Environment**: Uses Midtrans sandbox credentials
- **Test Data**: All data adalah test data, bukan production
- **Local Storage**: JWT tokens stored di browser localStorage
- **HTTPS Not Required**: Testing environment allows HTTP

## 🚀 Advanced Usage

### Custom Test Data
Edit default values di script.js:
```javascript
// Default donation amount
document.getElementById('donation-amount').value = '100000';

// Default test credentials
const TEST_CREDENTIALS = {
    streamer: { email: 'streamer@test.com', password: 'password123' },
    donator: { email: 'donator@test.com', password: 'password123' }
};
```

### Extended Logging
Enable verbose logging:
```javascript
// Add to script.js
const DEBUG_MODE = true;
console.log('Debug mode enabled');
```

## 📊 Performance

### Frontend Metrics
- **Load Time**: < 2 seconds
- **API Response Time**: Tracked dan displayed
- **Real-time Updates**: 30 second health check intervals
- **Memory Usage**: Optimized dengan log rotation

### Backend Integration
- **Connection Pooling**: Efficient API connections
- **Error Handling**: Graceful degradation
- **Retry Logic**: Auto-retry failed requests
- **Timeout Handling**: Request timeouts configured

---

**Happy Testing! 🎉**

Untuk questions atau issues, check backend logs atau browser console untuk detailed debugging information. 