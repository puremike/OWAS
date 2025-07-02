// document.addEventListener('DOMContentLoaded', async () => {
//     renderNav(document.getElementById('nav-buttons'));
//     // await loadWonAuctions();

//     console.log('Script loaded');
// document.addEventListener('DOMContentLoaded', async () => {
//     console.log('DOM fully loaded');
//     renderNav(document.getElementById('nav-buttons'));
//     await loadWonAuctions();
// });
// });

// async function loadWonAuctions() {
//     try {
//         const res = await fetch(`${API_BASE_URL}/auctions/won`, {
//             credentials: 'include'
//         });

//         console.log('Status:', res.status);
//         if (!res.ok) throw new Error('Failed to load won auctions');

//         const auctions = await res.json();
//         const container = document.getElementById('auction-list');
//         container.innerHTML = '';

//         if (!Array.isArray(auctions) || auctions.length === 0) {
//             container.innerHTML = '<p>You have not won any auctions yet.</p>';
//             return;
//         }

//         console.log(auctions.is_paid);

//         auctions.forEach(auction => {
//             console.log(auctions.is_paid);
//             const div = document.createElement('div');
//             div.className = 'border p-4 rounded bg-white shadow';

//             let paymentStatus = auction.is_paid ? '✅ Paid' : '❌ Not Paid';
//             let paymentButton = '';

//             if (!auction.is_paid) {
//                 paymentButton = `
//                     <button class="pay-btn bg-green-500 text-white px-3 py-1 rounded mt-2" data-id="${auction.id}">
//                         Make Payment
//                     </button>
//                 `;
//             }

//             div.setAttribute('data-auction-id', auction.id);  // Add data attribute for easy DOM targeting

//             div.innerHTML = `
//                 <h2 class="font-bold text-lg">${auction.title}</h2>
//                 <p><strong>Auction ID:</strong> ${auction.id}</p>
//                 <p><strong>Status:</strong> ${auction.status}</p>
//                 <p><strong>Final Price:</strong> $${auction.current_price}</p>
//                 <p><strong>Payment Status:</strong> ${paymentStatus}</p>
//                 ${paymentButton}
//             `;

//             container.appendChild(div);
//         });

//         // Attach event listeners to all "Make Payment" buttons
//         document.querySelectorAll('.pay-btn').forEach(btn => {
//             btn.addEventListener('click', async (e) => {
//                 const auctionId = e.target.getAttribute('data-id');
//                 try {
//                     const response = await fetch(`${API_BASE_URL}/auctions/${auctionId}/stripe/create-checkout-session`, {
//                         method: 'POST',
//                         credentials: 'include'
//                     });

//                     const data = await response.json();
//                     if (data.checkout_url) {
//                         window.location.href = data.checkout_url;  // Redirect to Stripe Checkout
//                     } else {
//                         alert('Failed to create checkout session.');
//                     }
//                 } catch (err) {
//                     console.error(err);
//                     alert('Error creating checkout session.');
//                 }
//             });
//         });

//     } catch (error) {
//         console.error(error);
//         window.location.href = 'auth.html';
//     }
// }


document.addEventListener('DOMContentLoaded', async () => {
    console.log('DOM fully loaded');
    renderNav(document.getElementById('nav-buttons'));
    await loadWonAuctions();
});

async function loadWonAuctions() {
    try {
        console.log('Starting loadWonAuctions()...');
        const res = await fetch(`${API_BASE_URL}/auctions/won`, {
            credentials: 'include'
        });

        if (res.status === 401) {
            window.location.href = 'auth.html';
            return;
        }
        
        if (!res.ok) {
            const errText = await res.text();  // Optional: for logging
            console.error(`Error ${res.status}: ${errText}`);
            document.getElementById('auction-list').innerHTML = '<p>Failed to load your auctions. Please try again later.</p>';
            return;
        }

        const auctions = await res.json();
        console.log('Auctions:', auctions);

        const container = document.getElementById('auction-list');
        container.innerHTML = '';

        if (!Array.isArray(auctions) || auctions.length === 0) {
            container.innerHTML = '<p>You have not won any auctions yet.</p>';
            return;
        }

        auctions.forEach(auction => {
            const div = document.createElement('div');
            div.className = 'border p-4 rounded bg-white shadow';

            const paymentStatus = auction.is_paid ? '✅ Paid' : '❌ Not Paid';
            const paymentButton = auction.is_paid ? '' : `
                <button class="pay-btn bg-green-500 text-white px-3 py-1 rounded mt-2" data-id="${auction.id}">
                    Make Payment
                </button>`;

            div.setAttribute('data-auction-id', auction.id);
            div.innerHTML = `
                <h2 class="font-bold text-lg">${auction.title}</h2>
                <p><strong>Auction ID:</strong> ${auction.id}</p>
                <p><strong>Status:</strong> ${auction.status}</p>
                <p><strong>Final Price:</strong> $${auction.current_price}</p>
                <p><strong>Payment Status:</strong> ${paymentStatus}</p>
                ${paymentButton}
            `;

            container.appendChild(div);
        });

        document.querySelectorAll('.pay-btn').forEach(btn => {
            btn.addEventListener('click', async (e) => {
                const auctionId = e.target.getAttribute('data-id');
                try {
                    const response = await fetch(`${API_BASE_URL}/auctions/${auctionId}/stripe/create-checkout-session`, {
                        method: 'POST',
                        credentials: 'include'
                    });

                    const data = await response.json();
                    if (data.checkout_url) {
                        window.location.href = data.checkout_url;
                    } else {
                        alert('Failed to create checkout session.');
                    }
                } catch (err) {
                    console.error(err);
                    alert('Error creating checkout session.');
                }
            });
        });

    } catch (error) {
        console.error('Caught error in loadWonAuctions:', error);
    }
}
