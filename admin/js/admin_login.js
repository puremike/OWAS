// const SET_BASE_URL = 'http://localhost:3100/api/v1';

document.getElementById('login-btn').addEventListener('click', async () => {
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    try {
        const result = await apiPost('/admin/login', { email, password });
        console.log('Login success:', result);
        window.location.href = 'admin_dashboard.html';

    } catch (error) {
        document.getElementById('error-message').classList.remove('hidden');
    }

});
