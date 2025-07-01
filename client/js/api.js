const API_BASE_URL = "http://localhost:3100/api/v1";
const API_BASE_URL_WS = "localhost:3100/api/v1/ws";

async function apiPost(endpoint, data) {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include',  // Important for cookie-based auth
    body: JSON.stringify(data)
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.message || `POST ${endpoint} failed`);
  }

  return await response.json();
}

async function apiRequest(endpoint, method = 'GET', data = null, requiresAuth = false) {
    const headers = { 'Content-Type': 'application/json' };
    
    const options = {
      method,
      headers,
      credentials: 'include'  // 🔑 This is critical for sending cookies!
    };
  
    if (data) {
      options.body = JSON.stringify(data);
    }
  
    const response = await fetch(`${API_BASE_URL}${endpoint}`, options);  // Replace URL as needed
  
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText || 'API request failed');
    }
  
    return response.json();
  }