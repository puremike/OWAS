document.addEventListener('DOMContentLoaded', async () => {
    renderNav(document.getElementById('nav-buttons'));

    try {
        await apiRequest('/me', 'GET', null, true); // check auth
    } catch (error) {
        console.error('JWT expired or user not authenticated.');
        window.location.href = 'auth.html';
        return;
    }

    await loadBiddedAuctionDetails();
});

async function loadBiddedAuctionDetails() {
    try {
        const auctions = await apiRequest('/auctions/bidded', 'GET', null, true);
        const detailsDiv = document.getElementById('auction-details');

        if (!auctions.length) {
            detailsDiv.innerHTML = '<p>You have not placed any bids yet.</p>';
            return;
        }

        detailsDiv.innerHTML = auctions.map(auction => `
            <div class="border p-4 mb-4 rounded shadow">
                ${auction.image_path ? `<img src="${auction.image_path}" alt="Auction Image" class="w-full h-48 object-cover rounded mb-2">` : ''}
                <h2 class="text-xl font-bold mb-2">${auction.title}</h2>
                <p>${auction.description}</p>
                <p><strong>Current Price:</strong> $${auction.current_price}</p>
                <p><strong>Status:</strong> ${auction.status}</p>
                <p><strong>Start:</strong> ${new Date(auction.start_time).toLocaleString()}</p>
                <p><strong>End:</strong> ${new Date(auction.end_time).toLocaleString()}</p>

                <form class="bid-form mt-4" data-auction-id="${auction.id}">
                    <input type="number" name="bid-amount" placeholder="Enter new bid" required class="border p-1 rounded w-1/2" step="0.01">
                    <button type="submit" class="bg-blue-500 text-white px-4 py-1 rounded">Bid</button>
                    <div class="bid-message text-sm mt-1"></div>
                </form>
            </div>
        `).join('');

        // Attach bid form submit handlers
        const bidForms = document.querySelectorAll('.bid-form');
        bidForms.forEach(form => {
            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                const auctionId = form.dataset.auctionId;
                const bidAmount = parseFloat(form.querySelector('input[name="bid-amount"]').value);
                const messageDiv = form.querySelector('.bid-message');

                try {
                    await apiRequest(`/auctions/${auctionId}/bids`, 'POST', { bidAmount }, true);
                    messageDiv.textContent = '✅ Bid placed successfully!';
                    messageDiv.className = 'bid-message text-green-600 mt-1';
                    // Optionally refresh the page or update the current_price
                    loadBiddedAuctionDetails(); 
                } catch (error) {
                    messageDiv.textContent = `❌ Failed: ${error.message}`;
                    messageDiv.className = 'bid-message text-red-600 mt-1';
                }
            });
        });

    } catch (error) {
        console.error('Failed to load bidded auctions:', error);
        document.getElementById('auction-details').innerHTML = `<p class="text-red-600">Error loading your bidded auctions.</p>`;
    }
}