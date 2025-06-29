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
        <a href="create_auction.html" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">Create Auctions</a>
        <a href="auctions.html" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">View Auctions</a>
        <button id="logout-btn" class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600">Logout</button>
      `;
  
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
  