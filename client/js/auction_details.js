document.addEventListener('DOMContentLoaded', async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const auctionId = urlParams.get('id');
  
    if (!auctionId) {
      document.getElementById('error-message').textContent = 'Invalid auction ID.';
      return;
    }

    loadAuctionDetails(auctionId);

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

    card.innerHTML = `
  <h3 class="text-xl font-semibold mb-2">${auction.title}</h3>
  <p class="text-gray-700 mb-2">${auction.description}</p>
  <p class="text-green-600 font-bold">Current Bid: $${auction.current_price || auction.starting_price || '0.00'}</p>
  <a href="auction_detail.html?id=${auction.id}" class="text-blue-500 hover:underline">View Details / Place Bid</a>
`;

        // document.getElementById('auction-title').textContent = auction.title;
        // document.getElementById('auction-description').textContent = auction.description;
        // document.getElementById('auction-type').textContent = `Type: ${auction.type}`;
        // document.getElementById('auction-starting-price').textContent = `Starting Bid: $${auction.starting_price}`;
        // document.getElementById('auction-current-price').textContent = `Current Bid: $${auction.current_price}`;
        // document.getElementById('auction-status').textContent = `Status: ${auction.status}`;
        // document.getElementById('auction-seller-id').textContent = `Seller ID: ${auction.seller_id}`;
        // document.getElementById('auction-end-date').textContent = `End Date: ${new Date(auction.end_time).toLocaleString()}`;
        // document.getElementById('auction-created-at').textContent = `Created At: ${new Date(auction.created_at).toLocaleString()}`;
  
      } catch (error) {
        // document.getElementById('error-message').textContent = 'Failed to load auction details.';
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
      loadAuctionDetails(auctionId);  // Refresh auction price
    } catch (error) {
      messageDiv.textContent = `❌ Failed to place bid: ${error.message}`;
      messageDiv.className = 'text-red-600';
    }
  }