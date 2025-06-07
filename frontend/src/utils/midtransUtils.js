// Midtrans utility functions for handling Snap.js and payments

/**
 * Check if Midtrans Snap.js is loaded
 */
export const isSnapLoaded = () => {
  return typeof window !== 'undefined' && window.snap && typeof window.snap.pay === 'function';
};

/**
 * Wait for Snap.js to load with timeout
 */
export const waitForSnap = (timeout = 5000) => {
  return new Promise((resolve) => {
    if (isSnapLoaded()) {
      resolve(true);
      return;
    }

    const checkInterval = 100;
    let elapsed = 0;

    const interval = setInterval(() => {
      elapsed += checkInterval;
      
      if (isSnapLoaded()) {
        clearInterval(interval);
        resolve(true);
      } else if (elapsed >= timeout) {
        clearInterval(interval);
        resolve(false);
      }
    }, checkInterval);
  });
};

/**
 * Load Midtrans Snap.js dynamically if not already loaded
 */
export const loadSnapScript = (clientKey) => {
  return new Promise((resolve, reject) => {
    if (isSnapLoaded()) {
      resolve(true);
      return;
    }

    // Check if script already exists
    const existingScript = document.querySelector('script[src*="snap.js"]');
    if (existingScript) {
      // Script exists but snap not ready yet, wait for it
      waitForSnap().then(resolve);
      return;
    }

    // Create and load script
    const script = document.createElement('script');
    script.type = 'text/javascript';
    script.src = 'https://app.sandbox.midtrans.com/snap/snap.js';
    script.setAttribute('data-client-key', clientKey);

    script.onload = () => {
      waitForSnap().then(resolve);
    };

    script.onerror = () => {
      reject(new Error('Failed to load Midtrans Snap.js'));
    };

    document.head.appendChild(script);
  });
};

/**
 * Process payment with Snap
 */
export const processSnapPayment = async (snapToken, callbacks = {}) => {
  try {
    // Ensure Snap is loaded
    const snapLoaded = await waitForSnap();
    
    if (!snapLoaded) {
      throw new Error('Midtrans Snap.js not available');
    }

    console.log('Processing payment with snap token:', snapToken);

    // Default callbacks
    const defaultCallbacks = {
      onSuccess: (result) => {
        console.log('Payment success:', result);
        callbacks.onSuccess?.(result);
      },
      onPending: (result) => {
        console.log('Payment pending:', result);
        callbacks.onPending?.(result);
      },
      onError: (result) => {
        console.log('Payment error:', result);
        callbacks.onError?.(result);
      },
      onClose: () => {
        console.log('Payment popup closed');
        callbacks.onClose?.();
      }
    };

    // Open Snap payment popup
    window.snap.pay(snapToken, defaultCallbacks);
    
    return true;
  } catch (error) {
    console.error('Error processing Snap payment:', error);
    throw error;
  }
};

/**
 * Fallback to redirect if Snap.js fails
 */
export const fallbackToRedirect = (redirectUrl) => {
  if (redirectUrl) {
    console.log('Falling back to redirect payment');
    window.open(redirectUrl, '_blank');
    return true;
  }
  return false;
}; 