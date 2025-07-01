document.addEventListener('DOMContentLoaded', async () => {
    renderNav(document.getElementById('nav-buttons'));
    await loadWonAuctions();
});

async function loadWonAuctions() {
    try {
        const res = await fetch(`${API_BASE_URL}/auctions/won`, {
            credentials: 'include'
        });

        if (!res.ok) throw new Error('Failed to load won auctions');

        const auctions = await res.json();
        const container = document.getElementById('auction-list');
        container.innerHTML = '';

        if (auctions.length === 0) {
            container.innerHTML = '<p>You have not won any auctions yet.</p>';
            return;
        }

        auctions.forEach(auction => {
            const div = document.createElement('div');
            div.className = 'border p-4 rounded bg-white shadow';

            div.innerHTML = `
                <h2 class="font-bold text-lg">${auction.title}</h2>
                <p><strong>Auction ID:</strong> ${auction.id}</p>
                <p><strong>Status:</strong> ${auction.status}</p>
                <p><strong>Final Price:</strong> $${auction.current_price}</p>
                <button class="pay-btn bg-green-500 text-white px-3 py-1 rounded mt-2" data-id="${auction.id}">Make Payment</button>
            `;

            container.appendChild(div);
        });

        document.querySelectorAll('.pay-btn').forEach(btn => {
            btn.addEventListener('click', async (e) => {
                const auctionId = e.target.getAttribute('data-id');
                const response = await fetch(`${API_BASE_URL}/auctions/${auctionId}/stripe/create-checkout-session`, {
                    method: 'POST',
                    credentials: 'include'
                });
        
                const data = await response.json();
                window.location.href = data.checkout_url;  // Stripe Checkout URL
            });
        });

    } catch (error) {
        console.error(error);
        document.getElementById('auction-list').innerHTML = '<p class="text-red-500">Failed to load auctions.</p>';
    }
}
