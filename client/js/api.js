const API_BASE_URL = "https://owas-server.onrender.com/api/v1";
const API_BASE_URL_WS = "owas-server.onrender.com/api/v1/ws";

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

async function apiRequest(path, method = 'GET', body = null, withAuth = false, retry = true) {
  const headers = {
      'Content-Type': 'application/json',
  };

  const options = {
      method,
      headers,
      credentials: 'include', // Ensure cookies (jwt + refresh_token) are sent
  };

  if (body) {
      options.body = JSON.stringify(body);
  }

  const response = await fetch(`${API_BASE_URL}${path}`, options);

  if (response.status === 401 && retry) {
      // Try refreshing the token
      const refreshResponse = await fetch(`${API_BASE_URL}/refresh`, {
          method: 'POST',
          credentials: 'include'
      });

      if (refreshResponse.ok) {
          // Retry original request once
          return apiRequest(path, method, body, withAuth, false);
      } else {
          throw new Error('Authentication expired. Please log in again.');
      }
  }

  if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Request failed');
  }

  return response.json();
}
