
// const SET_BASE_URL = 'http://localhost:3100/api/v1';

let auctionToDelete = null;

document.addEventListener('DOMContentLoaded', async () => {
    renderNav(document.getElementById('nav-buttons'));
    await loadAuctions();
});

async function loadAuctions() {
    const res = await fetch(`${API_BASE_URL}/auctions`, { credentials: 'include' });
    const auctions = await res.json();

    const container = document.getElementById('auction-list');
    container.innerHTML = '';

    let payMsg = auctions.is_paid;

        if (payMsg === true) {
            payMsg = '✅ Paid';
        }
        else {
            payMsg = '❌ Not Paid';
        }

    auctions.forEach(auction => {
        const div = document.createElement('div');
        div.className = 'border p-4 rounded flex justify-between items-center';

        div.innerHTML = `
          <div>
            <h2 class="font-semibold">${auction.title}</h2>
            <p>ID: ${auction.id}</p>
            <p class="font-semibold">Payment Status: ${payMsg}</p>
          </div>
          <button data-id="${auction.id}" class="delete-btn bg-red-500 text-white px-3 py-1 rounded">Delete</button>
        `;

        container.appendChild(div);
    });

    document.querySelectorAll('.delete-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            auctionToDelete = e.target.getAttribute('data-id');
            document.getElementById('delete-modal').classList.remove('hidden');
        });
    });
}

// Confirm Delete
document.getElementById('confirm-delete').addEventListener('click', async () => {
    if (auctionToDelete) {
        await fetch(`${API_BASE_URL}/admin/auctions/${auctionToDelete}`, {
            method: 'DELETE',
            credentials: 'include'
        });

        document.getElementById('delete-modal').classList.add('hidden');
        showSuccessMessage('Auction deleted successfully!');
        await loadAuctions();
    }
});

// Cancel Delete
document.getElementById('cancel-delete').addEventListener('click', () => {
    document.getElementById('delete-modal').classList.add('hidden');
    auctionToDelete = null;
});

// Success message function
function showSuccessMessage(text) {
    const msg = document.getElementById('success-message');
    msg.textContent = text;
    msg.classList.remove('hidden');
    setTimeout(() => msg.classList.add('hidden'), 3000);
}
