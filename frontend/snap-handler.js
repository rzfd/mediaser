/**
 * Snap.js Handler for MediaShar
 * Handles Midtrans Snap payment integration with proper error handling
 */

class SnapHandler {
    constructor() {
        this.isInitialized = false;
        this.snapToken = null;
        this.callbacks = {};
        this.messageListeners = [];
        
        this.init();
    }
    
    init() {
        console.log('🔧 Initializing SnapHandler...');
        
        // Wait for Snap.js to be loaded
        this.waitForSnap().then(() => {
            this.isInitialized = true;
            console.log('✅ SnapHandler initialized successfully');
        }).catch((error) => {
            console.error('❌ Failed to initialize SnapHandler:', error);
        });
        
        // Setup message listeners
        this.setupMessageListeners();
        
        // Setup cleanup on page unload
        window.addEventListener('beforeunload', () => this.cleanup());
    }
    
    waitForSnap(timeout = 10000) {
        return new Promise((resolve, reject) => {
            const startTime = Date.now();
            
            const checkSnap = () => {
                if (typeof window.snap !== 'undefined') {
                    resolve();
                } else if (Date.now() - startTime > timeout) {
                    reject(new Error('Snap.js failed to load within timeout'));
                } else {
                    setTimeout(checkSnap, 100);
                }
            };
            
            checkSnap();
        });
    }
    
    setupMessageListeners() {
        const messageHandler = (event) => {
            try {
                // Only handle messages from Midtrans domains
                const allowedOrigins = [
                    'https://app.sandbox.midtrans.com',
                    'https://app.midtrans.com',
                    'https://simulator.sandbox.midtrans.com'
                ];
                
                if (!allowedOrigins.includes(event.origin)) {
                    return;
                }
                
                console.log('📨 Received message from Midtrans:', event.data);
                
                // Handle the message data
                if (event.data && typeof event.data === 'object') {
                    this.handleMidtransMessage(event.data);
                }
            } catch (error) {
                console.warn('⚠️ Error handling Midtrans message:', error);
            }
        };
        
        window.addEventListener('message', messageHandler);
        this.messageListeners.push(messageHandler);
    }
    
    handleMidtransMessage(data) {
        // Handle different types of Midtrans messages
        if (data.status_code) {
            console.log('💳 Transaction status update:', data.status_code);
        }
        
        if (data.transaction_status) {
            console.log('📊 Transaction status:', data.transaction_status);
        }
        
        // Trigger custom callbacks if registered
        if (this.callbacks.onMessage) {
            this.callbacks.onMessage(data);
        }
    }
    
    async pay(snapToken, options = {}) {
        if (!this.isInitialized) {
            throw new Error('SnapHandler not initialized. Please wait for initialization.');
        }
        
        if (!snapToken) {
            throw new Error('Snap token is required');
        }
        
        console.log('💰 Starting Snap payment with token:', snapToken.substring(0, 20) + '...');
        
        return new Promise((resolve, reject) => {
            try {
                const defaultOptions = {
                    onSuccess: (result) => {
                        console.log('✅ Payment successful:', result);
                        resolve({ type: 'success', data: result });
                    },
                    onPending: (result) => {
                        console.log('⏳ Payment pending:', result);
                        resolve({ type: 'pending', data: result });
                    },
                    onError: (result) => {
                        console.log('❌ Payment error:', result);
                        reject(new Error('Payment failed: ' + JSON.stringify(result)));
                    },
                    onClose: () => {
                        console.log('🔒 Payment popup closed');
                        resolve({ type: 'closed', data: null });
                    }
                };
                
                // Merge user options with defaults
                const finalOptions = { ...defaultOptions, ...options };
                
                // Use setTimeout to ensure proper cleanup of promises
                setTimeout(() => {
                    window.snap.pay(snapToken, finalOptions);
                }, 0);
                
            } catch (error) {
                console.error('💥 Error in snap.pay:', error);
                reject(error);
            }
        });
    }
    
    setCallbacks(callbacks) {
        this.callbacks = { ...this.callbacks, ...callbacks };
    }
    
    cleanup() {
        console.log('🧹 Cleaning up SnapHandler...');
        
        // Remove message listeners
        this.messageListeners.forEach(listener => {
            window.removeEventListener('message', listener);
        });
        this.messageListeners = [];
        
        // Clear callbacks
        this.callbacks = {};
        
        console.log('✅ SnapHandler cleanup completed');
    }
    
    isReady() {
        return this.isInitialized && typeof window.snap !== 'undefined';
    }
}

// Create global instance
window.snapHandler = new SnapHandler();

// Export for module systems if needed
if (typeof module !== 'undefined' && module.exports) {
    module.exports = SnapHandler;
} 