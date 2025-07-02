// document.addEventListener('DOMContentLoaded', async function () {
//     const navButtons = document.getElementById('nav-buttons');
  
//     try {
//       // ✅ Check if user is authenticated (backend `/me` endpoint)
//       const user = await apiRequest('/me', 'GET', null, true);
  
//       // ✅ If logged in → Show Logout button
//       navButtons.innerHTML = `
//         <span class="text-gray-700">Hello, ${user.username}</span>
//         <button id="logout-btn" class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600">Logout</button>
//       `;
  
//       document.getElementById('logout-btn').addEventListener('click', async () => {
//         try {
//           await apiRequest('/logout', 'POST', null, true);
//           window.location.reload();
//         } catch (err) {
//           alert('Logout failed: ' + err.message);
//         }
//       });
  
//     } catch (err) {
//       // ✅ Not logged in → Show Login/Signup button
//       navButtons.innerHTML = `
//         <a href="auth.html" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">Login / Signup</a>
//       `;
//     }
//   });
  

document.addEventListener('DOMContentLoaded', async () => {
    const navButtons = document.getElementById('nav-buttons');
    const startBiddingBtn = document.getElementById('start-bidding-btn');  // Assuming you gave the "Start Bidding Today" button this ID
    
    try {
      // Check if user is logged in (by calling /me)
      await apiRequest('/me');
  
      // ✅ User is logged in
      navButtons.innerHTML = `
  <div class="relative">
    <button id="menu-btn" class="text-gray-800 focus:outline-none">
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16m-7 6h7"></path>
      </svg>
    </button>
    <div id="menu-dropdown" class="hidden absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-2 z-10">
      <a href="create_auction.html" class="block px-4 py-2 text-gray-800 hover:bg-blue-500 hover:text-white">Create Auctions</a>
      <a href="auctions.html" class="block px-4 py-2 text-gray-800 hover:bg-blue-500 hover:text-white">View Auctions</a>
       <a href="auction_bidded.html" class="block px-4 py-2 text-gray-800 hover:bg-blue-500 hover:text-white">My Bidded Auctions</a>
       <a href="my_auctions.html" class="block px-4 py-2 text-gray-800 hover:bg-blue-500 hover:text-white">My Won Auctions</a>
      <a href="profile.html" class="block px-4 py-2 text-gray-800 hover:bg-blue-500 hover:text-white">Profile</a>
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
  
      // ✅ Change Start Bidding button too (if you have it on the page)
      if (startBiddingBtn) {
        startBiddingBtn.textContent = 'View Auctions';
        startBiddingBtn.onclick = () => {
          window.location.href = 'auctions.html';
        };
      }
  
    } catch (error) {
      // Not logged in → Show Login / Signup buttons
      navButtons.innerHTML = `
        <a href="auth.html" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">Login / Signup</a>
      `;
  
      // Update start bidding button too (if needed)
      if (startBiddingBtn) {
        startBiddingBtn.textContent = 'Login / Signup to Start Bidding';
        startBiddingBtn.onclick = () => {
          window.location.href = 'auth.html';
        };
      }
    }
  
    // Optional: Handle Logout click
    document.addEventListener('click', async (event) => {
      if (event.target.id === 'logout-btn') {
        try {
          await apiRequest('/logout', 'POST');
          window.location.reload();
        } catch (error) {
          console.error('Logout failed', error);
        }
      }
    });
  });
  