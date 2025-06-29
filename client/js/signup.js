document.getElementById('signup-form').addEventListener('submit', async function (e) {
    e.preventDefault();
  
    const full_name = document.getElementById('full_name').value.trim();
    const username = document.getElementById('username').value.trim();
    const email = document.getElementById('email').value.trim();
    const location = document.getElementById('location').value.trim();
    const password = document.getElementById('password').value;
    const confirm_password = document.getElementById('confirm_password').value;
  
    const errorElement = document.getElementById('signup-error');
    errorElement.textContent = '';
  
    if (password !== confirm_password) {
      errorElement.textContent = 'Passwords do not match.';
      return;
    }
  
    const payload = {
      full_name,
      username,
      email,
      location,
      password,
      confirm_password
    };
  
    try {
      const response = await apiRequest('/signup', 'POST', payload);  // Assuming /signup is your API endpoint path
  
      if (response && response.success) {
        window.location.href = 'login.html';
      } else {
        errorElement.textContent = response.message || 'Signup failed. Please try again.';
      }
    } catch (error) {
      console.error(error);
      errorElement.textContent = 'An error occurred. Please try again later.';
    }
  });
  