document.getElementById('login-form').addEventListener('submit', async (e) => {
    e.preventDefault();
  
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value.trim();
    const errorMsg = document.getElementById('login-error');
    errorMsg.textContent = '';
  
    try {
      const result = await apiPost('/login', { email, password });
      console.log('Login success:', result);
  
      // Redirect to homepage after successful login
      window.location.href = 'index.html';
    } catch (error) {
      console.error('Login error:', error);
      errorMsg.textContent = error.message;
    }
  });
  