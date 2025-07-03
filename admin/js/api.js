// const API_BASE_URL = "https://owas-server.onrender.com/api/v1";

const isLocalhost = window.location.hostname === "localhost";

const API_BASE_URL = isLocalhost
  ? "http://localhost:3100/api/v1"
  : "https://owas-server.onrender.com/api/v1";

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