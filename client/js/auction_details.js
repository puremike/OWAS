let imagePath;

// async function checkAuthAndLoadAuctionsSec() {
//     try {
//         await apiRequest('/me', 'GET', null, true);
//         loadAuctions();
//     } catch (error) {
//         console.error('Auth check failed:', error);
//         window.location.href = 'auth.html';
//     }
// }

document.addEventListener('DOMContentLoaded', async () => {
    try {
       await apiRequest('/me', 'GET', null, true);  // Auth check passed
        document.getElementById('page-body').classList.remove('hidden');
       
    } catch (error) {
        console.error('JWT expired or user not authenticated.');
        window.location.href = 'auth.html';
        return; // Prevent rest of the page from loading
    }

    
    renderNav(document.getElementById('nav-buttons'));

    const urlParams = new URLSearchParams(window.location.search);
    const auctionId = urlParams.get('id');

    if (!auctionId) {
        document.getElementById('error-message').textContent = 'Invalid auction ID.';
        return;
    }

    await loadAuctionDetails(auctionId);

    const bidForm = document.getElementById('bid-form');
    bidForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        placeBid(auctionId);
    });
});

async function loadAuctionDetails(id) {
    try {
        const auction = await apiRequest(`/auctions/${id}`);
        const detailsDiv = document.getElementById('auction-details');

        imagePath = auction.image_path;

        detailsDiv.innerHTML = `
            ${auction.image_path ? `<img src="${auction.image_path}" alt="Auction Image" class="w-full h-48 object-cover rounded mb-2">` : ''}
            <h2 class="text-2xl font-bold mb-2">${auction.title}</h2>
            <p>${auction.description}</p>
            <p><strong>Current Price:</strong> $${auction.current_price}</p>
            <p><strong>Status:</strong> ${auction.status}</p>
            <p><strong>Start Time:</strong> ${auction.start_time}</p>
            <p><strong>End Time:</strong> ${auction.end_time}</p>
            <p><strong>Seller_ID:</strong> ${auction.seller_id}</p>
        `;

        // Show Close Auction button only for the seller and if auction is still OPEN
        const actionDiv = document.getElementById('action-buttons');
        actionDiv.innerHTML = '';

        try {
            const currentUser = await apiRequest('/me');
            if (currentUser.id === auction.seller_id) {
                const closeBtn = document.createElement('button');
                closeBtn.textContent = 'Close Auction';
                closeBtn.className = 'bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600';
                closeBtn.addEventListener('click', async () => {
                    await closeAuction(id);
                });
                actionDiv.appendChild(closeBtn);

                // NEW: Update Auction Button
                const updateBtn = document.createElement('button');
                updateBtn.textContent = 'Update Auction';
                updateBtn.className = 'bg-yellow-500 text-white px-4 py-2 rounded hover:bg-yellow-600 mr-2 ml-2';
                updateBtn.addEventListener('click', () => {
                    showUpdateForm(auction);
                });
                actionDiv.appendChild(updateBtn);


                const deleteBtn = document.createElement('button');
deleteBtn.textContent = 'Delete Auction';
deleteBtn.className = 'bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600';

// When user clicks delete button, show modal
deleteBtn.addEventListener('click', () => {
    document.getElementById('delete-modal').classList.remove('hidden');
});

// Confirm delete
document.getElementById('confirm-delete').addEventListener('click', async () => {
    await deleteAuction(id);
    document.getElementById('delete-modal').classList.add('hidden');

    const msg = document.getElementById('success-message');
    msg.textContent = 'Auction deleted successfully!';
    msg.classList.remove('hidden');

    setTimeout(() => {
        msg.classList.add('hidden');
        window.location.href = 'auctions.html';  // Redirect back to auction list
    }, 300);
});

// Cancel delete
document.getElementById('cancel-delete').addEventListener('click', () => {
    document.getElementById('delete-modal').classList.add('hidden');
});

// Append button to action div
actionDiv.appendChild(deleteBtn);
                //   // NEW: Delete Auction Button
                //   const deleteBtn = document.createElement('button');
                //   deleteBtn.textContent = 'Delete Auction';
                //   deleteBtn.className = 'bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600';
                //   deleteBtn.addEventListener('click', async () => {
                //       if (confirm('Are you sure you want to delete this auction?')) {
                //           await deleteAuction(id);
                //       }
                //   });
                // actionDiv.appendChild(deleteBtn);
            }
        } catch (userError) {
            actionDiv.innerHTML = 'User not logged in';
            window.location.href = 'auth.html';
            // console.error('Failed to fetch current user:', userError);
        }

    } catch (error) {
        console.error(error);
    }
}

async function placeBid(auctionId) {
    const bidAmount = document.getElementById('bid-amount').value;
    const messageDiv = document.getElementById('bid-message');

    try {
        await apiRequest(`/auctions/${auctionId}/bids`, 'POST', { bidAmount: parseFloat(bidAmount) }, true);
        messageDiv.textContent = '✅ Bid placed successfully!';
        messageDiv.className = 'text-green-600';
        loadAuctionDetails(auctionId);
    } catch (error) {
        messageDiv.textContent = `❌ Failed to place bid: ${error.message}`;
        messageDiv.className = 'text-red-600';
    }
}

async function closeAuction(auctionId) {
    const messageDiv = document.getElementById('close-message');
    try {
        const result = await apiRequest(`/auctions/${auctionId}/close`, 'POST', null, true);
        messageDiv.innerHTML = `
            ✅ Auction Closed Successfully!<br>
            Winning Bid: ₦${result.winning_bid} <br>
            Status: ${result.status}
        `;
        messageDiv.className = 'text-green-600 mt-4';
        await loadAuctionDetails(auctionId);
    } catch (error) {
        console.error('Failed to close auction:', error);
        messageDiv.textContent = `❌ Failed to close auction: ${error.message}`;
        messageDiv.className = 'text-red-600 mt-4';
    }
}

function formatDateForInput(dateString) {
    const date = new Date(dateString);
    if (isNaN(date)) {
      return ''; // handle invalid date if needed
    }
    return date.toISOString().split('T')[0]; // gives "2025-06-29"
  }

// 🔴 NEW: Show Update Auction Form
function showUpdateForm(auction) {
    const formDiv = document.getElementById('update-form');

    formDiv.innerHTML = `
        <h3 class="text-xl font-bold mb-2">Update Auction</h3>
        <form id="auction-update-form" class="space-y-2">
            <input type="text" id="update-title" value="${auction.title}" placeholder="Title" class="border p-2 w-full">
            <textarea id="update-description" placeholder="Description" class="border p-2 w-full">${auction.description}</textarea>
            <input type="number" id="update-starting-price" value="${auction.starting_price}" placeholder="Starting Price" class="border p-2 w-full">
            <input type="date" id="update-start-time" value="${formatDateForInput(auction.start_time)}" class="border p-2 w-full">
  <input type="date" id="update-end-time" value="${formatDateForInput(auction.end_time)}" class="border p-2 w-full">
            
            <select id="update-auction-type" class="border p-2 w-full">
                <option value="english" ${auction.type === 'english' ? 'selected' : ''}>English Auction</option>
                <option value="dutch" ${auction.type === 'dutch' ? 'selected' : ''}>Dutch Auction</option>
                <option value="sealed" ${auction.type === 'sealed' ? 'selected' : ''}>Sealed Bid Auction</option>
            </select>

            <button type="submit" class="bg-green-500 mt-3 text-white px-4 py-2 rounded hover:bg-green-600">Submit Update</button>
        </form>
    `;

    document.getElementById('auction-update-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        await submitAuctionUpdate(auction.id);
    });
}

// 🔴 NEW: Submit Auction Update
async function submitAuctionUpdate(auctionId) {
    const messageDiv = document.getElementById('update-message');

    const payload = {
        title: document.getElementById('update-title').value,
        description: document.getElementById('update-description').value,
        starting_price: parseFloat(document.getElementById('update-starting-price').value),
        type: document.getElementById('update-auction-type').value,
        start_time: document.getElementById('update-start-time').value,
        end_time: document.getElementById('update-end-time').value,
        image_path: imagePath
    };

    try {
        await apiRequest(`/auctions/${auctionId}`, 'PUT', payload, true);
        messageDiv.textContent = '✅ Auction updated successfully!';
        messageDiv.className = 'text-green-600';
        await loadAuctionDetails(auctionId);
        document.getElementById('update-form').innerHTML = '';  // Hide form after update
    } catch (error) {
        console.error('Failed to update auction:', error);
        messageDiv.textContent = `❌ Failed to update auction: ${error.message}`;
        messageDiv.className = 'text-red-600';
    }
}

// 🔴 NEW: Delete Auction
async function deleteAuction(auctionId) {
    try {
        await apiRequest(`/auctions/${auctionId}`, 'DELETE', null, true);
        window.location.href = 'auctions.html';  // Redirect to homepage or auction list
    } catch (error) {
        console.error('Failed to delete auction:', error);
    }
}






















// document.addEventListener('DOMContentLoaded', async () => {
//     const urlParams = new URLSearchParams(window.location.search);
//     const auctionId = urlParams.get('id');
  
//     if (!auctionId) {
//       document.getElementById('error-message').textContent = 'Invalid auction ID.';
//       return;
//     }

//     loadAuctionDetails(auctionId);

//     const bidForm = document.getElementById('bid-form');
//   bidForm.addEventListener('submit', async (e) => {
//     e.preventDefault();
//     placeBid(auctionId);
//   });
  
//   });

//   async function loadAuctionDetails(id) {
//     try {
//         const auction = await apiRequest(`/auctions/${id}`);
//         const detailsDiv = document.getElementById('auction-details');

//         detailsDiv.innerHTML = `
//         ${auction.image_path ? `<img src="${auction.image_path}" alt="Auction Image" class="w-full h-48 object-cover rounded mb-2">` : ''}
//       <h2 class="text-2xl font-bold mb-2">${auction.title}</h2>
//       <p>${auction.description}</p>
//       <p><strong>Current Price:</strong> $${auction.current_price}</p>
//       <p><strong>Status:</strong> ${auction.status}</p>
//       <p><strong>Start Time:</strong> ${auction.start_time}</p>
//       <p><strong>End Time:</strong> ${auction.end_time}</p>
//       <p><strong>Seller_ID:</strong> ${auction.seller_id}</p>
//     `;

//     card.innerHTML = `
//   <h3 class="text-xl font-semibold mb-2">${auction.title}</h3>
//   <p class="text-gray-700 mb-2">${auction.description}</p>
//   <p class="text-green-600 font-bold">Current Bid: $${auction.current_price || auction.starting_price || '0.00'}</p>
//   <a href="auction_detail.html?id=${auction.id}" class="text-blue-500 hover:underline">View Details / Place Bid</a>
// `;

//         // document.getElementById('auction-title').textContent = auction.title;
//         // document.getElementById('auction-description').textContent = auction.description;
//         // document.getElementById('auction-type').textContent = `Type: ${auction.type}`;
//         // document.getElementById('auction-starting-price').textContent = `Starting Bid: $${auction.starting_price}`;
//         // document.getElementById('auction-current-price').textContent = `Current Bid: $${auction.current_price}`;
//         // document.getElementById('auction-status').textContent = `Status: ${auction.status}`;
//         // document.getElementById('auction-seller-id').textContent = `Seller ID: ${auction.seller_id}`;
//         // document.getElementById('auction-end-date').textContent = `End Date: ${new Date(auction.end_time).toLocaleString()}`;
//         // document.getElementById('auction-created-at').textContent = `Created At: ${new Date(auction.created_at).toLocaleString()}`;
  
//       } catch (error) {
//         // document.getElementById('error-message').textContent = 'Failed to load auction details.';
//         console.error(error);
//       }
//   }
  
//   async function placeBid(auctionId) {
//     const bidAmount = document.getElementById('bid-amount').value;
//     const messageDiv = document.getElementById('bid-message');
  
//     try {
//       await apiRequest(`/auctions/${auctionId}/bids`, 'POST', { bidAmount: parseFloat(bidAmount) }, true);
//       messageDiv.textContent = '✅ Bid placed successfully!';
//       messageDiv.className = 'text-green-600';
//       loadAuctionDetails(auctionId);  // Refresh auction price
//     } catch (error) {
//       messageDiv.textContent = `❌ Failed to place bid: ${error.message}`;
//       messageDiv.className = 'text-red-600';
//     }
//   }