// Configuration
const API_BASE_URL = 'http://localhost:8080/api';
const HEALTH_URL = 'http://localhost:8080/health';

// Global state
let currentUser = null;
let authToken = localStorage.getItem('authToken');
let lastDonationId = null;
let lastSnapToken = null;
let lastOrderId = null;

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    console.log('🚀 DOM Content Loaded - Initializing...');
    
    // Message handling is now managed by SnapHandler
    
    checkHealth();
    
    // If token exists, try to get user profile
    if (authToken) {
        getUserProfile();
    }
    
    // Wait a bit for DOM to be fully ready, then load streamers
    setTimeout(() => {
        console.log('⏰ Loading streamers after timeout...');
        window.loadStreamers();
    }, 1000);
    
    updateSessionData();
    
    // Auto-refresh health status every 30 seconds
    setInterval(checkHealth, 30000);
    
    // Add event listener for Load Streamers button
    const loadStreamersBtn = document.getElementById('load-streamers-btn');
    if (loadStreamersBtn) {
        loadStreamersBtn.addEventListener('click', function() {
            console.log('🔄 Load Streamers button clicked');
            window.loadStreamers();
        });
    }
    
    // Clean up event listeners on page unload
    window.addEventListener('beforeunload', function() {
        // Clean up any pending promises or callbacks
        console.log('🧹 Cleaning up before page unload');
    });
});

// Override console errors to catch and handle Snap.js related errors
const originalError = console.error;
console.error = function(...args) {
    const errorMsg = args.join(' ');
    if (errorMsg.includes('message channel closed') || 
        errorMsg.includes('asynchronous response')) {
        console.warn('Handled Snap.js communication error:', errorMsg);
        return;
    }
    originalError.apply(console, args);
};

// Utility Functions
function showNotification(message, type = 'info') {
    const notification = document.getElementById('notification');
    const icon = document.getElementById('notification-icon');
    const messageEl = document.getElementById('notification-message');
    
    notification.className = `fixed top-4 right-4 p-4 rounded-lg shadow-lg z-50 ${type}`;
    
    switch(type) {
        case 'success':
            icon.className = 'fas fa-check-circle mr-2';
            break;
        case 'error':
            icon.className = 'fas fa-exclamation-circle mr-2';
            break;
        case 'warning':
            icon.className = 'fas fa-exclamation-triangle mr-2';
            break;
        default:
            icon.className = 'fas fa-info-circle mr-2';
    }
    
    messageEl.textContent = message;
    notification.classList.remove('hidden');
    
    setTimeout(() => {
        notification.classList.add('hidden');
    }, 5000);
}

function logResponse(endpoint, method, request, response, success = true) {
    const logContainer = document.getElementById('api-log');
    const timestamp = new Date().toLocaleTimeString();
    
    const logEntry = document.createElement('div');
    logEntry.className = `mb-4 p-3 rounded ${success ? 'bg-green-50 border-l-4 border-green-400' : 'bg-red-50 border-l-4 border-red-400'}`;
    
    logEntry.innerHTML = `
        <div class="font-semibold ${success ? 'text-green-800' : 'text-red-800'} mb-2">
            [${timestamp}] ${method} ${endpoint} - ${success ? 'SUCCESS' : 'ERROR'}
        </div>
        ${request ? `<div class="mb-2"><strong>Request:</strong><pre class="text-xs">${JSON.stringify(request, null, 2)}</pre></div>` : ''}
        <div><strong>Response:</strong><pre class="text-xs">${JSON.stringify(response, null, 2)}</pre></div>
    `;
    
    // Remove first child if too many entries
    if (logContainer.children.length > 10) {
        logContainer.removeChild(logContainer.firstElementChild);
    }
    
    logContainer.appendChild(logEntry);
    logContainer.scrollTop = logContainer.scrollHeight;
}

function updateSessionData() {
    document.getElementById('current-token').textContent = authToken || 'Not logged in';
    document.getElementById('current-donation-id').textContent = lastDonationId || 'None';
    document.getElementById('current-snap-token').textContent = lastSnapToken || 'None';
    document.getElementById('current-order-id').textContent = lastOrderId || 'None';
}

function showUserInfo() {
    if (currentUser) {
        document.getElementById('user-info').classList.remove('hidden');
        document.getElementById('user-details').textContent = 
            `${currentUser.username} (${currentUser.email}) - ${currentUser.is_streamer ? 'Streamer' : 'Donator'}`;
    } else {
        document.getElementById('user-info').classList.add('hidden');
    }
}

function addTestResult(test, success, message = '') {
    const resultsContainer = document.getElementById('test-results');
    const result = document.createElement('div');
    result.className = `flex items-center p-2 rounded ${success ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`;
    result.innerHTML = `
        <i class="fas ${success ? 'fa-check' : 'fa-times'} mr-2"></i>
        <span><strong>${test}</strong> ${message ? `- ${message}` : ''}</span>
    `;
    resultsContainer.appendChild(result);
}

// Load streamers for dropdown
window.loadStreamers = async function loadStreamers() {
    console.log('🔄 Loading streamers...');
    try {
        console.log('📡 Making API request to /streamers');
        const response = await apiRequest('/streamers', 'GET');
        console.log('📥 API Response:', response);
        
        const streamers = response.data || response;
        console.log('📋 Streamers data:', streamers);
        
        const streamerSelect = document.getElementById('streamer-id');
        console.log('🎯 Streamer select element:', streamerSelect);
        
        if (!streamerSelect) {
            console.error('❌ Streamer select element not found!');
            return;
        }
        
        // Clear existing options except the first one
        streamerSelect.innerHTML = '<option value="">Select Streamer...</option>';
        
        // Add streamers to dropdown
        if (Array.isArray(streamers) && streamers.length > 0) {
            streamers.forEach((streamer, index) => {
                console.log(`➕ Adding streamer ${index + 1}:`, streamer);
                const option = document.createElement('option');
                option.value = streamer.id;
                option.textContent = `${streamer.username}${streamer.full_name ? ' - ' + streamer.full_name : ''}`;
                streamerSelect.appendChild(option);
            });
            
            console.log(`✅ Loaded ${streamers.length} streamers successfully`);
            showNotification(`Loaded ${streamers.length} streamers`, 'success');
        } else {
            console.warn('⚠️ No streamers found in response');
            showNotification('No streamers available', 'warning');
        }
        
        // Add refresh button if not exists
        window.addStreamerRefreshButton();
        
    } catch (error) {
        console.error('❌ Failed to load streamers:', error);
        showNotification(`Failed to load streamers: ${error.message}`, 'error');
        
        // Add refresh button even on error
        window.addStreamerRefreshButton();
    }
}

// Add refresh button next to streamer dropdown
window.addStreamerRefreshButton = function addStreamerRefreshButton() {
    const streamerSelect = document.getElementById('streamer-id');
    const container = streamerSelect.parentNode;
    
    // Check if refresh button already exists
    if (container.querySelector('.streamer-refresh-btn')) {
        return;
    }
    
    const refreshBtn = document.createElement('button');
    refreshBtn.className = 'streamer-refresh-btn px-3 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition text-sm';
    refreshBtn.innerHTML = '<i class="fas fa-sync-alt"></i>';
    refreshBtn.title = 'Refresh Streamers';
    refreshBtn.onclick = function(e) {
        e.preventDefault();
        window.loadStreamers();
        showNotification('Refreshing streamers...', 'info');
    };
    
    // Add refresh button to the container
    container.appendChild(refreshBtn);
}

// API Functions
async function apiRequest(endpoint, method = 'GET', data = null, requireAuth = false) {
    const url = endpoint.startsWith('http') ? endpoint : `${API_BASE_URL}${endpoint}`;
    
    const headers = {
        'Content-Type': 'application/json',
    };
    
    if (requireAuth && authToken) {
        headers['Authorization'] = `Bearer ${authToken}`;
    }
    
    const config = {
        method,
        headers,
        mode: 'cors',
        credentials: 'include',
    };
    
    if (data && method !== 'GET') {
        config.body = JSON.stringify(data);
    }
    
    try {
        const response = await fetch(url, config);
        
        // Handle different response types
        let responseData;
        const contentType = response.headers.get('content-type');
        
        if (contentType && contentType.includes('application/json')) {
            responseData = await response.json();
        } else {
            responseData = await response.text();
        }
        
        logResponse(endpoint, method, data, responseData, response.ok);
        
        // Handle specific HTTP status codes
        if (response.status === 401) {
            showNotification('Session expired. Please login again.', 'warning');
            logout(); // Clear invalid token
            throw new Error('Authentication required');
        }
        
        if (response.status === 403) {
            showNotification('Access denied. Insufficient permissions.', 'error');
            throw new Error('Access forbidden');
        }
        
        if (!response.ok) {
            const errorMessage = responseData.message || responseData.error || `HTTP ${response.status}`;
            throw new Error(errorMessage);
        }
        
        return responseData;
        
    } catch (error) {
        // Handle different types of errors
        if (error.name === 'TypeError' && error.message.includes('fetch')) {
            const networkError = {
                error: 'Network Error',
                message: 'Cannot connect to backend. Please check if backend is running.',
                details: error.message
            };
            logResponse(endpoint, method, data, networkError, false);
            throw new Error('Network Error: Cannot connect to backend');
        }
        
        if (error.message.includes('CORS')) {
            const corsError = {
                error: 'CORS Error',
                message: 'Cross-Origin Request Blocked',
                details: error.message
            };
            logResponse(endpoint, method, data, corsError, false);
            throw new Error('CORS Error: Please check backend configuration');
        }
        
        // For other errors, log and re-throw
        logResponse(endpoint, method, data, { error: error.message }, false);
        throw error;
    }
}

// Health Check Functions
async function checkHealth() {
    try {
        document.getElementById('status-indicator').className = 'w-3 h-3 bg-yellow-400 rounded-full animate-pulse';
        document.getElementById('status-text').textContent = 'Checking...';
        
        const health = await apiRequest(HEALTH_URL);
        
        // Update overall status
        document.getElementById('status-indicator').className = 'w-3 h-3 bg-green-400 rounded-full';
        document.getElementById('status-text').textContent = 'Connected';
        
        // Update individual status
        document.getElementById('api-status').textContent = health.status === 'ok' ? 'Healthy' : 'Error';
        document.getElementById('db-status').textContent = health.database?.status === 'connected' ? 'Connected' : 'Disconnected';
        document.getElementById('midtrans-status').textContent = health.midtrans === 'configured' ? 'Configured' : 'Not Configured';
        
        showNotification('Health check successful', 'success');
        
    } catch (error) {
        console.error('Health check error:', error);
        document.getElementById('status-indicator').className = 'w-3 h-3 bg-red-400 rounded-full';
        document.getElementById('status-text').textContent = 'Error';
        document.getElementById('api-status').textContent = 'Error';
        document.getElementById('db-status').textContent = 'Error';
        document.getElementById('midtrans-status').textContent = 'Error';
        
        // More specific error handling
        if (error.message.includes('CORS')) {
            showNotification('CORS Error: Please check backend CORS configuration', 'error');
        } else if (error.message.includes('fetch')) {
            showNotification('Network Error: Cannot connect to backend. Please check if backend is running.', 'error');
        } else {
            showNotification(`Health check failed: ${error.message}`, 'error');
        }
    }
}

// Authentication Functions
async function registerUser() {
    const username = document.getElementById('reg-username').value;
    const email = document.getElementById('reg-email').value;
    const password = document.getElementById('reg-password').value;
    const isStreamer = document.getElementById('reg-type').value === 'true';
    
    if (!username || !email || !password) {
        showNotification('Please fill all registration fields', 'error');
        return;
    }
    
    try {
        const response = await apiRequest('/auth/register', 'POST', {
            username,
            email,
            password,
            is_streamer: isStreamer
        });
        
        showNotification('User registered successfully!', 'success');
        addTestResult('User Registration', true, `User: ${username}`);
        
        // Auto-login after successful registration
        if (response.data?.token) {
            authToken = response.data.token;
            localStorage.setItem('authToken', authToken);
            currentUser = response.data.user;
            updateSessionData();
            showUserInfo();
            showNotification('Auto-logged in after registration!', 'info');
        }
        
        // Clear form
        document.getElementById('reg-username').value = '';
        document.getElementById('reg-email').value = '';
        document.getElementById('reg-password').value = '';
        
    } catch (error) {
        showNotification(`Registration failed: ${error.message}`, 'error');
        addTestResult('User Registration', false, error.message);
    }
}

async function loginUser() {
    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;
    
    if (!email || !password) {
        showNotification('Please fill email and password', 'error');
        return;
    }
    
    try {
        const response = await apiRequest('/auth/login', 'POST', {
            email,
            password
        });
        
        // Fix: Extract token from response.data.token
        authToken = response.data?.token || response.token;
        localStorage.setItem('authToken', authToken);
        
        // Fix: Extract user from response.data.user
        currentUser = response.data?.user || response.user;
        
        showNotification('Login successful!', 'success');
        addTestResult('User Login', true, `Email: ${email}`);
        
        updateSessionData();
        showUserInfo();
        
    } catch (error) {
        showNotification(`Login failed: ${error.message}`, 'error');
        addTestResult('User Login', false, error.message);
    }
}

async function getUserProfile() {
    if (!authToken) {
        showNotification('Please login first', 'warning');
        return;
    }
    
    try {
        const response = await apiRequest('/auth/profile', 'GET', null, true);
        currentUser = response.data || response;
        updateSessionData();
        showUserInfo();
        showNotification('Profile updated', 'success');
    } catch (error) {
        console.error('Failed to get user profile:', error);
        showNotification(`Failed to get profile: ${error.message}`, 'error');
    }
}

async function quickLogin(type) {
    // Pre-fill with test credentials
    if (type === 'streamer') {
        document.getElementById('login-email').value = 'streamer@test.com';
        document.getElementById('login-password').value = 'password123';
    } else {
        document.getElementById('login-email').value = 'donator@test.com';
        document.getElementById('login-password').value = 'password123';
    }
    
    await loginUser();
}

function logout() {
    authToken = null;
    currentUser = null;
    localStorage.removeItem('authToken');
    
    updateSessionData();
    showUserInfo();
    
    showNotification('Logged out successfully', 'info');
}

// Donation Functions
async function createDonation() {
    if (!authToken) {
        showNotification('Please login first', 'error');
        return;
    }
    
    const amount = parseFloat(document.getElementById('donation-amount').value);
    const currency = document.getElementById('donation-currency').value;
    const streamerSelect = document.getElementById('streamer-id');
    const streamerId = parseInt(streamerSelect.value);
    const message = document.getElementById('donation-message').value;
    const displayName = document.getElementById('display-name').value;
    const isAnonymous = document.getElementById('is-anonymous').checked;
    
    if (!amount || !streamerId) {
        showNotification('Please fill amount and select a streamer', 'error');
        return;
    }
    
    try {
        const response = await apiRequest('/donations', 'POST', {
            amount,
            currency,
            streamer_id: streamerId,
            message,
            display_name: displayName,
            is_anonymous: isAnonymous
        }, true);
        
        lastDonationId = response.data?.id || response.id;
        document.getElementById('payment-donation-id').value = lastDonationId;
        
        updateSessionData();
        
        showNotification('Donation created successfully!', 'success');
        addTestResult('Create Donation', true, `ID: ${lastDonationId}, Amount: ${amount} ${currency}`);
        
    } catch (error) {
        showNotification(`Failed to create donation: ${error.message}`, 'error');
        addTestResult('Create Donation', false, error.message);
    }
}

// Payment Functions
async function createPayment() {
    const donationId = parseInt(document.getElementById('payment-donation-id').value);
    
    if (!donationId) {
        showNotification('Please enter donation ID', 'error');
        return;
    }
    
    if (!authToken) {
        showNotification('Please login first', 'error');
        return;
    }
    
    try {
        const response = await apiRequest(`/midtrans/payment/${donationId}`, 'POST', null, true);
        
        lastSnapToken = response.data?.token || response.token;
        lastOrderId = response.data?.order_id || response.order_id;
        
        updateSessionData();
        
        // Enable Snap payment button
        document.getElementById('snap-button').disabled = false;
        
        showNotification('Payment created successfully!', 'success');
        addTestResult('Create Payment', true, `Order ID: ${lastOrderId}`);
        
    } catch (error) {
        showNotification(`Failed to create payment: ${error.message}`, 'error');
        addTestResult('Create Payment', false, error.message);
    }
}

async function openSnapPayment() {
    if (!lastSnapToken) {
        showNotification('No snap token available. Create payment first.', 'error');
        return;
    }
    
    // Check if SnapHandler is ready
    if (!window.snapHandler || !window.snapHandler.isReady()) {
        showNotification('Snap.js is not ready. Please wait and try again.', 'error');
        return;
    }
    
    try {
        const result = await window.snapHandler.pay(lastSnapToken, {
            onSuccess: function(result) {
                showNotification('Payment successful!', 'success');
                addTestResult('Snap Payment', true, 'Payment completed successfully');
                console.log('Payment Success:', result);
                logResponse('Snap Payment Success', 'CALLBACK', null, result, true);
            },
            onPending: function(result) {
                showNotification('Payment pending', 'warning');
                addTestResult('Snap Payment', true, 'Payment is pending');
                console.log('Payment Pending:', result);
                logResponse('Snap Payment Pending', 'CALLBACK', null, result, true);
            },
            onError: function(result) {
                showNotification('Payment failed', 'error');
                addTestResult('Snap Payment', false, 'Payment failed');
                console.log('Payment Error:', result);
                logResponse('Snap Payment Error', 'CALLBACK', null, result, false);
            },
            onClose: function() {
                showNotification('Payment popup closed', 'info');
                console.log('Payment popup closed');
            }
        });
        
        // Handle the result
        switch (result.type) {
            case 'success':
                showNotification('Payment completed successfully!', 'success');
                addTestResult('Snap Payment Result', true, 'Payment successful');
                break;
            case 'pending':
                showNotification('Payment is being processed', 'warning');
                addTestResult('Snap Payment Result', true, 'Payment pending');
                break;
            case 'closed':
                showNotification('Payment popup was closed', 'info');
                break;
        }
        
    } catch (error) {
        console.error('Error opening Snap payment:', error);
        showNotification(`Failed to process payment: ${error.message}`, 'error');
        addTestResult('Snap Payment', false, error.message);
    }
}

// Test Functions
async function fullFlowTest() {
    clearLogs();
    showNotification('Starting full flow test...', 'info');
    
    try {
        // 1. Health Check
        await checkHealth();
        addTestResult('Health Check', true);
        
        // 2. Register if not logged in
        if (!authToken) {
            const testEmail = `test-${Date.now()}@example.com`;
            const testUsername = `testuser${Date.now()}`;
            
            await apiRequest('/auth/register', 'POST', {
                username: testUsername,
                email: testEmail,
                password: 'password123',
                is_streamer: false
            });
            addTestResult('Auto Registration', true, testEmail);
            
            // 3. Login
            const loginResponse = await apiRequest('/auth/login', 'POST', {
                email: testEmail,
                password: 'password123'
            });
            
            authToken = loginResponse.token;
            localStorage.setItem('authToken', authToken);
            await getUserProfile();
            updateSessionData();
            showUserInfo();
            addTestResult('Auto Login', true);
        }
        
        // 4. Get available streamer ID from dropdown
        const streamerSelect = document.getElementById('streamer-id');
        let streamerId = parseInt(streamerSelect.value);
        
        // If no streamer selected, try to use the first available one
        if (!streamerId && streamerSelect.options.length > 1) {
            streamerId = parseInt(streamerSelect.options[1].value);
        }
        
        if (!streamerId) {
            throw new Error('No streamers available. Please create a streamer first.');
        }
        
        // 4. Create Donation
        const donationResponse = await apiRequest('/donations', 'POST', {
            amount: 75000,
            currency: 'IDR',
            streamer_id: streamerId,
            message: 'Full flow test donation',
            display_name: 'Test Supporter',
            is_anonymous: false
        }, true);
        
        lastDonationId = donationResponse.data?.id || donationResponse.id;
        addTestResult('Auto Create Donation', true, `ID: ${lastDonationId}`);
        
        // 5. Create Payment
        const paymentResponse = await apiRequest(`/payments/${lastDonationId}`, 'POST', null, true);
        
        lastSnapToken = paymentResponse.data?.snap_token || paymentResponse.snap_token;
        lastOrderId = paymentResponse.data?.order_id || paymentResponse.order_id;
        
        updateSessionData();
        addTestResult('Auto Create Payment', true, `Order: ${lastOrderId}`);
        
        // Enable payment button
        document.getElementById('snap-button').disabled = false;
        
        showNotification('Full flow test completed successfully! You can now test payment.', 'success');
        
    } catch (error) {
        showNotification(`Full flow test failed: ${error.message}`, 'error');
        addTestResult('Full Flow Test', false, error.message);
    }
}

async function loadTestData() {
    try {
        showNotification('Loading test data...', 'info');
        
        // Create test streamer if not exists
        try {
            await apiRequest('/auth/register', 'POST', {
                username: 'teststreamer',
                email: 'streamer@test.com',
                password: 'password123',
                is_streamer: true
            });
        } catch (e) {
            // User might already exist
        }
        
        // Create test donator if not exists
        try {
            await apiRequest('/auth/register', 'POST', {
                username: 'testdonator',
                email: 'donator@test.com',
                password: 'password123',
                is_streamer: false
            });
        } catch (e) {
            // User might already exist
        }
        
        // Reload streamers dropdown after creating test data
        await window.loadStreamers();
        
        showNotification('Test data loaded successfully!', 'success');
        addTestResult('Load Test Data', true, 'Test users created');
        
    } catch (error) {
        showNotification(`Failed to load test data: ${error.message}`, 'error');
        addTestResult('Load Test Data', false, error.message);
    }
}

async function viewDonations() {
    try {
        const donations = await apiRequest('/donations', 'GET', null, true);
        
        const logContainer = document.getElementById('api-log');
        const logEntry = document.createElement('div');
        logEntry.className = 'mb-4 p-3 rounded bg-blue-50 border-l-4 border-blue-400';
        
        logEntry.innerHTML = `
            <div class="font-semibold text-blue-800 mb-2">Donations List (${donations.data?.length || donations.length || 0} items)</div>
            <pre class="text-xs">${JSON.stringify(donations, null, 2)}</pre>
        `;
        
        logContainer.appendChild(logEntry);
        logContainer.scrollTop = logContainer.scrollHeight;
        
        showNotification('Donations loaded successfully', 'success');
        
    } catch (error) {
        showNotification(`Failed to load donations: ${error.message}`, 'error');
    }
}

function clearLogs() {
    document.getElementById('api-log').innerHTML = '<div class="text-gray-500">Logs cleared...</div>';
    document.getElementById('test-results').innerHTML = '<div class="text-gray-500">Test results cleared...</div>';
}

// Keyboard shortcuts
document.addEventListener('keydown', function(e) {
    // Ctrl/Cmd + Enter to run full flow test
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
        e.preventDefault();
        fullFlowTest();
    }
    
    // Ctrl/Cmd + H to check health
    if ((e.ctrlKey || e.metaKey) && e.key === 'h') {
        e.preventDefault();
        checkHealth();
    }
    
    // Ctrl/Cmd + L to clear logs
    if ((e.ctrlKey || e.metaKey) && e.key === 'l') {
        e.preventDefault();
        clearLogs();
    }
});

// Add keyboard shortcut info to page
document.addEventListener('DOMContentLoaded', function() {
    const shortcuts = document.createElement('div');
    shortcuts.className = 'fixed bottom-4 left-4 bg-gray-800 text-white p-3 rounded-lg text-xs';
    shortcuts.innerHTML = `
        <div class="font-semibold mb-1">Keyboard Shortcuts:</div>
        <div>Ctrl+Enter: Full Flow Test</div>
        <div>Ctrl+H: Health Check</div>
        <div>Ctrl+L: Clear Logs</div>
    `;
    document.body.appendChild(shortcuts);
}); 