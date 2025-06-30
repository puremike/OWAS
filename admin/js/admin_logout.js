document.getElementById('logout-btn').addEventListener('click', async () => {
    
    try {
        const result = await apiPost('/admin/logout');
        console.log('Logout success:', result);
        window.location.href = 'admin_login.html';

    } catch (error) {
        document.getElementById('error-message').classList.remove('hidden');
    }
});
