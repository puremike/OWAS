function showMessage(message, type = 'error') {
    const box = document.getElementById('message-box');
    box.textContent = message;
    box.classList.remove('hidden');
  
    if (type === 'error') {
      box.className = 'mb-4 text-center text-sm px-4 py-2 rounded bg-red-100 text-red-700 border border-red-300';
    } else if (type === 'success') {
      box.className = 'mb-4 text-center text-sm px-4 py-2 rounded bg-green-100 text-green-700 border border-green-300';
    }
  }

document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('login-form');
    const signupForm = document.getElementById('signup-form');
    const formTitle = document.getElementById('form-title');
    const toggleButton = document.getElementById('toggle-button');
    const toggleMessage = document.getElementById('toggle-message');
  
    let showingLogin = true;
  
    // Toggle between Login and Signup forms
    toggleButton.addEventListener('click', () => {
      showingLogin = !showingLogin;
  
      if (showingLogin) {
        signupForm.classList.add('hidden');
        loginForm.classList.remove('hidden');
        formTitle.textContent = 'Login';
        toggleMessage.textContent = "Don't have an account?";
        toggleButton.textContent = "Signup";
      } else {
        loginForm.classList.add('hidden');
        signupForm.classList.remove('hidden');
        formTitle.textContent = 'Signup';
        toggleMessage.textContent = "Already have an account?";
        toggleButton.textContent = "Login";
      }
    });
  
    // Handle Login Form Submission
    loginForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const email = document.getElementById('login-email').value;
      const password = document.getElementById('login-password').value;
  
      try {
        await apiRequest('/login', 'POST', { email, password });
        window.location.href = 'index.html';
      } catch (err) {
        showMessage('Login failed: ' + err.message, 'error');
      }
    });
  
    // Handle Signup Form Submission
    signupForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const full_name = document.getElementById('signup-fullname').value;
      const username = document.getElementById('signup-username').value;
      const email = document.getElementById('signup-email').value;
      const location = document.getElementById('signup-location').value;
      const password = document.getElementById('signup-password').value;
      const confirm_password = document.getElementById('signup-confirm-password').value;
  
      if (password !== confirm_password) {
        showMessage('Passwords do not match.', 'error');
        return;
      }

      const passwordRequirements = /^(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*()_\-+=~`[\]{}|:;"'<>,.?/]).{6,}$/;
      
      if (!passwordRequirements.test(password)) {showMessage('Password must be at least 6 characters long and include at least one uppercase letter, one number, and one special character.', 'error');
        return;
    }
  
      try {
        await apiRequest('/signup', 'POST', {
          full_name,
          username,
          email,
          location,
          password,
          confirm_password
        });
        showMessage('Signup successful! You can now log in.', 'success');
        showingLogin = true;
        toggleButton.click(); // Automatically switch to Login form
      } catch (err) {
        showMessage('Signup failed: ' + err.message, 'error');
      }
    });
  });
  