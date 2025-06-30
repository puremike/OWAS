// const SET_BASE_URL = 'http://localhost:3100/api/v1';
let userToDelete = null;

document.addEventListener('DOMContentLoaded', async () => {
  renderNav(document.getElementById('nav-buttons'));
  await loadUsers();

  document.getElementById('confirm-delete').addEventListener('click', deleteUser);
  document.getElementById('cancel-delete').addEventListener('click', () => {
    userToDelete = null;
    document.getElementById('delete-modal').classList.add('hidden');
  });
});

async function loadUsers() {
  const res = await fetch(`${API_BASE_URL}/admin/users`, { method: 'GET', credentials: 'include' });
  const users = await res.json();

  const container = document.getElementById('user-list');
  container.innerHTML = '';

  users.forEach(user => {
    const div = document.createElement('div');
    div.className = 'border p-4 rounded bg-white shadow';

    div.innerHTML = `
      <p><strong>ID:</strong> ${user.id}</p>
      <p><strong>Username:</strong> ${user.username}</p>
      <p><strong>Email:</strong> ${user.email}</p>
      <p><strong>Name:</strong> ${user.full_name}</p>
      <p><strong>Location:</strong> ${user.location}</p>
      <button data-id="${user.id}" class="delete-btn bg-red-500 text-white px-3 py-1 rounded hover:bg-red-600 mt-2">Delete</button>
    `;

    container.appendChild(div);
  });

  document.querySelectorAll('.delete-btn').forEach(btn => {
    btn.addEventListener('click', () => {
      userToDelete = btn.getAttribute('data-id');
      document.getElementById('delete-modal').classList.remove('hidden');
    });
  });
}

async function deleteUser() {
  if (!userToDelete) return;

  const res = await fetch(`${API_BASE_URL}/admin/users/${userToDelete}`, {
    method: 'DELETE',
    credentials: 'include'
  });

  if (res.ok) {
    document.getElementById('success-message').classList.remove('hidden');
    setTimeout(() => {
      document.getElementById('success-message').classList.add('hidden');
    }, 3000);
    await loadUsers();
  } else {
    document.getElementById('error-message').classList.remove('hidden');
    setTimeout(() => {
      document.getElementById('error-message').classList.add('hidden');
    }, 3000);
  }

  userToDelete = null;
  document.getElementById('delete-modal').classList.add('hidden');
}
