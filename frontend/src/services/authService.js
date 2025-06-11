const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

class AuthService {
  async register(userData) {
    try {
      const response = await fetch(`${API_BASE_URL}/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Registration failed');
      }

      return data;
    } catch (error) {
      console.error('Registration error:', error);
      throw error;
    }
  }

  async login(credentials) {
    try {
      const response = await fetch(`${API_BASE_URL}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Login failed');
      }

      console.log('Login response data:', data);

      // Store token in localStorage - backend returns { success, message, data: { user, token } }
      const token = data.data?.token || data.token;
      const user = data.data?.user || data.user;
      
      console.log('Token from response:', token ? token.substring(0, 20) + '...' : 'No token found');
      console.log('User from response:', user);

      if (token) {
        localStorage.setItem('authToken', token);
        console.log('Token stored in localStorage');
      } else {
        console.error('No token found in login response');
      }

      if (user) {
        localStorage.setItem('user', JSON.stringify(user));
        console.log('User stored in localStorage');
      }

      return { ...data, token, user };
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  }

  async logout() {
    localStorage.removeItem('authToken');
    localStorage.removeItem('user');
  }

  getCurrentUser() {
    const user = localStorage.getItem('user');
    return user ? JSON.parse(user) : null;
  }

  getToken() {
    return localStorage.getItem('authToken');
  }

  isAuthenticated() {
    return !!this.getToken();
  }

  async loginWithGoogle(credential) {
    try {
      const response = await fetch(`${API_BASE_URL}/auth/google`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ credential }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Google login failed');
      }

      console.log('Google login response data:', data);

      // Store token in localStorage - backend returns { success, message, data: { user, token } }
      const token = data.data?.token || data.token;
      const user = data.data?.user || data.user;
      
      console.log('Token from Google response:', token ? token.substring(0, 20) + '...' : 'No token found');
      console.log('User from Google response:', user);

      if (token) {
        localStorage.setItem('authToken', token);
        console.log('Google token stored in localStorage');
      } else {
        console.error('No token found in Google login response');
      }

      if (user) {
        localStorage.setItem('user', JSON.stringify(user));
        console.log('Google user stored in localStorage');
      }

      return { ...data, token, user };
    } catch (error) {
      console.error('Google login error:', error);
      throw error;
    }
  }
}

export default new AuthService(); 