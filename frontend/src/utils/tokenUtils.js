// Token utility functions to handle authentication tokens safely

/**
 * Get a valid token from localStorage
 * Returns null if token is invalid or malformed
 */
export const getValidToken = () => {
  const token = localStorage.getItem('authToken');
  
  // Check if token exists and is not undefined/null string
  if (!token || token === 'undefined' || token === 'null' || token === '') {
    return null;
  }
  
  // Basic JWT format validation (should have 3 parts separated by dots)
  const parts = token.split('.');
  if (parts.length !== 3) {
    console.warn('Invalid token format detected, cleaning up...');
    cleanupInvalidToken();
    return null;
  }
  
  return token;
};

/**
 * Clean up invalid tokens from localStorage
 */
export const cleanupInvalidToken = () => {
  localStorage.removeItem('authToken');
  localStorage.removeItem('user');
};

/**
 * Check if user is properly authenticated with a valid token
 */
export const isAuthenticated = () => {
  return getValidToken() !== null;
};

/**
 * Get authorization headers for API requests
 * Only includes Authorization header if valid token exists
 */
export const getAuthHeaders = () => {
  const headers = {
    'Content-Type': 'application/json'
  };
  
  const token = getValidToken();
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  
  return headers;
};

/**
 * Safe API fetch with automatic token handling
 */
export const apiRequest = async (url, options = {}) => {
  const defaultOptions = {
    headers: getAuthHeaders(),
    ...options
  };
  
  // Merge headers if options already has headers
  if (options.headers) {
    defaultOptions.headers = {
      ...getAuthHeaders(),
      ...options.headers
    };
  }
  
  try {
    const response = await fetch(url, defaultOptions);
    
    // If we get 401 Unauthorized, clean up invalid token
    if (response.status === 401) {
      console.warn('Received 401 Unauthorized, cleaning up token...');
      cleanupInvalidToken();
      throw new Error('Authentication required');
    }
    
    return response;
  } catch (error) {
    // If it's a token-related error, clean up
    if (error.message.includes('token') || error.message.includes('malformed')) {
      cleanupInvalidToken();
    }
    throw error;
  }
}; 