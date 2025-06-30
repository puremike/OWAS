document.addEventListener('DOMContentLoaded', () => {
    renderNav(document.getElementById('nav-buttons'));
    loadUserInfo();

    document.getElementById('password-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        await changePassword();
    });

});

async function loadUserInfo() {
    try {
        const user = await apiRequest('/me', 'GET', null, true);
        const userInfoDiv = document.getElementById('user-info');
        userInfoDiv.innerHTML = `
            <p><strong>Username:</strong> ${user.username}</p>
            <p><strong>Email:</strong> ${user.email}</p>
            <p><strong>Name:</strong> ${user.full_name}</p>
            <p><strong>Address:</strong> ${user.location}</p>
        `;

        // Set avatar image
        const avatarImg = document.getElementById('user-avatar');
        avatarImg.src = `https://ui-avatars.com/api/?name=${encodeURIComponent(user.username)}&background=random&size=128`;

    } catch (error) {
        window.location.href = 'auth.html';
        // showMessage('Failed to load user info.', 'error');
    }
}

async function changePassword() {
    const current = document.getElementById('current-password').value;
    const newPass = document.getElementById('new-password').value;
    const confirm = document.getElementById('confirm-password').value;

    if (newPass !== confirm) {
        showMessage('New passwords do not match.', 'error');
        return;
    }

    try {
        await apiRequest('/change-password', 'POST', {
            current_password: current,
            new_password: newPass
        }, true);
        showMessage('Password updated successfully.', 'success');
        document.getElementById('password-form').reset();
    } catch (error) {
        showMessage('Failed to update password.', 'error');
    }
}

function showMessage(message, type) {
    const msgDiv = document.getElementById('message');
    msgDiv.textContent = message;
    msgDiv.className = `mt-4 text-center ${type === 'error' ? 'text-red-600' : 'text-green-600'}`;
    msgDiv.classList.remove('hidden');
}
