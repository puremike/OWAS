console.log("✅ auctions.js is running");

let currentPage = 1;
const pageSize = 10;

let filters = { type: '', status: ''};

document.addEventListener('DOMContentLoaded', () => {
    checkAuthAndLoadAuctions();

    document.getElementById('logout-btn').addEventListener('click', async () => {
        try {
            await apiRequest('/logout', 'POST', null, true);
            window.location.href = 'index.html';
        } catch (error) {
            showMessage('Error logging out.', 'error');
        }
    });

    document.getElementById('apply-filters').addEventListener('click', () => {
        filters.type = document.getElementById('type-filter').value;
        filters.status = document.getElementById('status-filter').value;

        currentPage = 1;
        loadAuctions();
    });

    document.getElementById('prev-page').addEventListener('click', () => {
        if (currentPage > 1) {
            currentPage--;
            loadAuctions();
        }
    });

    document.getElementById('next-page').addEventListener('click', () => {
        currentPage++;
        loadAuctions();
    });
});

async function checkAuthAndLoadAuctions() {
    try {
        await apiRequest('/me', 'GET', null, true);
        loadAuctions();
    } catch (error) {
        console.error('Auth check failed:', error);
        window.location.href = 'auth.html';
    }
}

async function loadAuctions() {
    try {
        const query = `limit=${pageSize}&offset=${(currentPage - 1) * pageSize}` +
                      `&type=${encodeURIComponent(filters.type)}` +
                      `&status=${encodeURIComponent(filters.status)}`;

        const auctions = await apiRequest(`/auctions?${query}`, 'GET', null, true);
        const container = document.getElementById('auction-list');
        container.innerHTML = '';

        if (auctions.length === 0) {
            container.innerHTML = '<p class="text-gray-500">No auctions available.</p>';
            document.getElementById('next-page').disabled = true;
            document.getElementById('page-info').textContent = `Page ${currentPage}`;
            return;
        }

        auctions.forEach(auction => {
            const card = document.createElement('div');
            card.className = 'bg-white p-4 rounded shadow hover:shadow-lg transition';

            const endTime = new Date(auction.end_time).getTime();
            const countdownId = `countdown-${auction.id}`;
            
                card.innerHTML = `
         ${auction.image_path ? `<img src="${auction.image_path}" alt="Auction Image" class="w-full h-48 object-cover rounded mb-2">` : ''}


          <h3 class="text-xl font-semibold">${auction.title}</h3>
          <p class="text-green-600 font-bold">Starting Price: $${auction.starting_price || '0.00'}</p>
          <p class="text-green-600 font-bold">Current Bid: $${auction.current_price || '0.00'}</p>
          <p class="text-gray-600">Start Time: ${new Date(auction.start_time).toLocaleString()}</p>
          <p class="text-gray-600">End Time: ${new Date(auction.end_time).toLocaleString()}</p>
          <p class="text-gray-600">Status: ${auction.status}</p>
          <a href="auction_detail.html?id=${auction.id}" class="inline-block bg-blue-500 text-white px-3 py-1 rounded hover:bg-blue-600 transition">View Details</a>

           <p id="${countdownId}" class="text-red-600 font-bold mt-2"></p>
        `;
            container.appendChild(card);

            // Start countdown timer for this auction
            startCountdown(endTime, countdownId);
        });

        document.getElementById('prev-page').disabled = currentPage === 1;
        document.getElementById('next-page').disabled = auctions.length < pageSize;
        document.getElementById('page-info').textContent = `Page ${currentPage}`;
    } catch (error) {
      console.error('API Error Details:', error);
        showMessage('Failed to load auctions.', 'error');
    }
}

function showMessage(message, type) {
    const msgDiv = document.getElementById('message');
    msgDiv.textContent = message;
    msgDiv.className = `mt-4 text-center ${type === 'error' ? 'text-red-600' : 'text-green-600'}`;
    msgDiv.classList.remove('hidden');
}

function startCountdown(endTime, countdownId) {
  function updateCountdown() {
      const now = new Date().getTime();
      const distance = endTime - now;
      const countdownElement = document.getElementById(countdownId);

      if (!countdownElement) return;

      if (distance <= 0) {
          countdownElement.textContent = "Auction ended";
          clearInterval(interval);
          return;
      }

      const hours = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
      const minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
      const seconds = Math.floor((distance % (1000 * 60)) / 1000);

      countdownElement.textContent = `Time Left: ${hours}h ${minutes}m ${seconds}s`;
  }

  updateCountdown();
  const interval = setInterval(updateCountdown, 1000);
}




// console.log("✅ auctions.js is running");

// document.addEventListener('DOMContentLoaded', () => {
//     checkAuthAndLoadAuctions();

//     document.getElementById('logout-btn').addEventListener('click', async () => {
//         try {
//             await apiRequest('/logout', 'POST', null, true);
//             window.location.href = 'index.html';
//         } catch (error) {
//             showMessage('Error logging out.', 'error');
//         }
//     });
// });

// async function checkAuthAndLoadAuctions() {
//     try {
//         const user = await apiRequest('/me', 'GET', null, true);
//         console.log('User data:', user);
//         loadAuctions();
//     } catch (error) {
//         console.error('Auth check failed:', error);
//         window.location.href = 'auth.html';
//     }
// }
  
//   async function loadAuctions() {
//     try {
//       const auctions = await apiRequest('/auctions', 'GET', null, true);
//       const container = document.getElementById('auction-list');
//       container.innerHTML = '';
  
//       if (auctions.length === 0) {
//         container.innerHTML = '<p class="text-gray-500">No auctions available at the moment.</p>';
//         return;
//       }

//       auctions.forEach(auction => {
//         const card = document.createElement('div');
//         card.className = 'bg-white p-4 rounded shadow hover:shadow-lg transition space-y-2';
  
//         card.innerHTML = `
//          ${auction.image_path ? `<img src="${auction.image_path}" alt="Auction Image" class="w-full h-48 object-cover rounded mb-2">` : ''}


//           <h3 class="text-xl font-semibold">${auction.title}</h3>
//           <p class="text-green-600 font-bold">Starting Price: $${auction.starting_price || '0.00'}</p>
//           <p class="text-green-600 font-bold">Current Bid: $${auction.current_price || '0.00'}</p>
//           <p class="text-gray-600">Start Time: ${new Date(auction.start_time).toLocaleString()}</p>
//           <p class="text-gray-600">End Time: ${new Date(auction.end_time).toLocaleString()}</p>
//           <p class="text-gray-600">Status: ${auction.status}</p>
//           <a href="auction_detail.html?id=${auction.id}" class="inline-block bg-blue-500 text-white px-3 py-1 rounded hover:bg-blue-600 transition">View Details</a>
//         `;
//         container.appendChild(card);
//       });

  
//     } catch (error) {
//       showMessage('Failed to load auctions.', 'error');
//     }
//   }
  
//   function showMessage(message, type) {
//     const msgDiv = document.getElementById('message');
//     msgDiv.textContent = message;
//     msgDiv.className = `mt-4 text-center ${type === 'error' ? 'text-red-600' : 'text-green-600'}`;
//     msgDiv.classList.remove('hidden');
//   }