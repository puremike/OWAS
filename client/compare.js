document.addEventListener('DOMContentLoaded', async () => {
    const navButtons = document.getElementById('nav-buttons');
    const startBiddingBtn = document.getElementById('start-bidding-btn');  // Assuming you gave the "Start Bidding Today" button this ID
    
    try {
      // Check if user is logged in (by calling /me)
      await apiRequest('/me');
  
      // ✅ User is logged in
      navButtons.innerHTML = `
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
        <a href="auth.html" class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600">Login / Signup</a>
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
  