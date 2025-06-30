function renderNav(navButtons) {
    // const navButtons = document.getElementById('nav-buttons');
    // if (!navButtons) return;  // Safety check in case element is missing

    navButtons.innerHTML = `
      <div class="relative">
        <button id="menu-btn" class="text-gray-800 focus:outline-none">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16m-7 6h7"></path>
          </svg>
        </button>
        <div id="menu-dropdown" class="hidden absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-2 z-10">
          <a href="admin_auctions.html" class="block px-4 py-2 text-gray-800 hover:bg-blue-500 hover:text-white">Manage Auctions</a>
          <a href="admin_users.html" class="block px-4 py-2 text-gray-800 hover:bg-blue-500 hover:text-white">Manage Users</a>
          <a href="index.html" class="block px-4 py-2 text-gray-800 hover:bg-blue-500 hover:text-white">Home</a>
          <button id="logout-btn" class="block w-full text-left px-4 py-2 text-gray-800 hover:bg-red-500 hover:text-white">Logout</button>
        </div>
      </div>
    `;

    // Toggle dropdown on menu button click
    document.getElementById('menu-btn').addEventListener('click', () => {
        const dropdown = document.getElementById('menu-dropdown');
        dropdown.classList.toggle('hidden');
    });

    // Close dropdown when clicking outside
    document.addEventListener('click', (event) => {
        const menu = document.getElementById('menu-dropdown');
        const menuBtn = document.getElementById('menu-btn');
        if (!menu.contains(event.target) && !menuBtn.contains(event.target)) {
            menu.classList.add('hidden');
        }
    });

    // Logout button event listener
    document.getElementById('logout-btn').addEventListener('click', async () => {
    
      try {
          const result = await apiPost('/admin/logout');
          console.log('Logout success:', result);
          window.location.href = 'admin_login.html';
    
      } catch (error) {
          document.getElementById('error-message').classList.remove('hidden');
      }
    });
      
}



