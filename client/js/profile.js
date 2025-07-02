document.addEventListener('DOMContentLoaded', () => {
    renderNav(document.getElementById('nav-buttons'));
    loadUserInfo();

    document.getElementById('password-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        await changePassword();
    });

    // Handle Delete Account button
    document.getElementById('delete-account-btn').addEventListener('click', () => {
        document.getElementById('delete-confirm-modal').classList.remove('hidden');
    });

    document.getElementById('cancel-delete-account').addEventListener('click', () => {
        document.getElementById('delete-confirm-modal').classList.add('hidden');
    });

    document.getElementById('confirm-delete-account').addEventListener('click', async () => {
        await deleteAccount();
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
    const old_password = document.getElementById('current-password').value;
    const new_password = document.getElementById('new-password').value;
    const confirm_password = document.getElementById('confirm-password').value;

    if (new_password !== confirm_password) {
        showMessage('New passwords do not match.', 'error');
        return;
    }

    if (new_password == old_password) {
        showMessage('New password cannot be the same as the old password.', 'error');
        return;
    }

    try {
        await apiRequest('/change-password', 'PUT', {
            old_password: old_password, new_password: new_password, confirm_password: confirm_password
        }, true);
        showMessage('Password updated successfully.', 'success');
        document.getElementById('password-form').reset();
        window.location.href = 'auth.html';
    } catch (error) {
        showMessage('Failed to update password.', 'error');
    }
}

async function deleteAccount() {
    try {
        await apiRequest('/users', 'DELETE', null, true);
        showDeleteMessage('Account deleted successfully. Redirecting to homepage...', 'success');
        setTimeout(() => window.location.href = 'index.html', 2000);
    } catch (error) {
        showDeleteMessage('Failed to delete account.', 'error');
    } finally {
        document.getElementById('delete-confirm-modal').classList.add('hidden');
    }
}

function showMessage(message, type) {
    const msgDiv = document.getElementById('message');
    msgDiv.textContent = message;
    msgDiv.className = `mt-4 text-center ${type === 'error' ? 'text-red-600' : 'text-green-600'}`;
    msgDiv.classList.remove('hidden');
}


function showDeleteMessage(message, type) {
    const msgDiv = document.getElementById('delete-message');
    msgDiv.textContent = message;
    msgDiv.className = `${type === 'error' ? 'text-red-600' : 'text-green-600'} mt-4`;
    msgDiv.classList.remove('hidden');
}