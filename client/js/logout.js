  
  document.addEventListener('DOMContentLoaded', async function () {
    const logoutBtn = document.getElementById('logout-btn');
  
    // Check if user is authenticated
    try {
      const user = await apiRequest('/me', 'GET', null, true);  // Replace `/me` with your backend's auth check endpoint
      if (user && logoutBtn) {
        logoutBtn.classList.remove('hidden');
      }
    } catch (error) {
      console.warn('User not authenticated:', error);
      // Optional: Redirect to login if user is not authenticated
      window.location.href = 'login.html';
    }
  
    // Handle Logout Button Click
    if (logoutBtn) {
      logoutBtn.addEventListener('click', async function () {
        try {
          await apiRequest('/logout', 'POST', {}, true);
          window.location.href = 'index.html';
        } catch (error) {
          console.error('Logout failed:', error);
          alert('Logout failed. Please try again.');
        }
      });
    }
  });
  