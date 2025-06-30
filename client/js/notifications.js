document.addEventListener('DOMContentLoaded', async () => {
    try {
        const notifications = await apiRequest('/notifications');
        const list = document.getElementById('notification-list');

        if (notifications.length === 0) {
            list.innerHTML = '<li>No notifications.</li>';
        } else {
            notifications.forEach(n => {
                const li = document.createElement('li');
                li.textContent = `${new Date(n.created_at).toLocaleString()}: ${n.message}`;
                list.appendChild(li);
            });
        }
    } catch (error) {
        console.error('Failed to load notifications:', error);
    }
});
