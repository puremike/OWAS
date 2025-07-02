document.addEventListener('DOMContentLoaded', () => {
    console.log('Main JS loaded and DOM ready.');
    // Initialize WebSocket connection
    connectWebSocket();
    

    // Optional: Call other global initializations if needed
    // Example: loadWonAuctions(); or loadActiveAuctions();
});

function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
    const wsUrl = protocol + API_BASE_URL_WS;

    const socket = new WebSocket(wsUrl);

    socket.onopen = () => {
        console.log('WebSocket connection opened');
    };

    socket.onmessage = (event) => {
        console.log('WebSocket message received:', event.data);

        try {
            const msg = JSON.parse(event.data);

            // Handle auction price updates
            if (msg.id && msg.current_price) {
                console.log('Auction Update:', msg);
                updateAuctionPrice(msg.id, msg.current_price);
            }

            // Handle user notifications
            if (msg.user_id && msg.message) {
                console.log('Notification:', msg);
                showNotification(msg.message);
            }

        } catch (e) {
            console.error('Failed to parse WebSocket message', e);
        }
    };

    socket.onclose = (event) => {
        console.warn('WebSocket closed. Reconnecting in 5 seconds...', event.reason);
        setTimeout(connectWebSocket, 5000);
    };

    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        socket.close();
    };
}

function updateAuctionPrice(auctionId, newPrice) {
    const auctionDiv = document.querySelector(`[data-auction-id="${auctionId}"]`);
    if (auctionDiv) {
        const priceP = auctionDiv.querySelector('.auction-price');
        if (priceP) {
            priceP.textContent = `Final Price: $${newPrice}`;
        }
    }
}

function showNotification(message) {
    // Create wrapper if it doesn't exist
    let wrapper = document.getElementById('toast-wrapper');
    if (!wrapper) {
        wrapper = document.createElement('div');
        wrapper.id = 'toast-wrapper';
        wrapper.className = 'fixed bottom-5 right-5 space-y-3 z-50';
        document.body.appendChild(wrapper);
    }

    // bg-indigo-100 border-indigo-300 text-indigo-900

    // Create toast
    const toast = document.createElement('div');
    toast.className = `
        max-w-xs w-full bg-indigo-100 border border-indigo-300 shadow-lg rounded-lg px-4 py-3 text-indigo-900 
        flex items-center justify-between gap-3 
        opacity-0 transform translate-y-2 transition-all duration-300 ease-in-out
    `;

    toast.innerHTML = `
        <span class="text-md font-medium">${message}</span>
        <button class="text-gray-400 hover:text-gray-600 focus:outline-none" onclick="this.parentElement.remove()">
            &times;
        </button>
    `;

    wrapper.appendChild(toast);

    // Animate in
    requestAnimationFrame(() => {
        toast.classList.remove('opacity-0', 'translate-y-2');
        toast.classList.add('opacity-100', 'translate-y-0');
    });

    // Auto-remove after 5 seconds
    setTimeout(() => {
        toast.classList.remove('opacity-100', 'translate-y-0');
        toast.classList.add('opacity-0', 'translate-y-2');
        setTimeout(() => toast.remove(), 1000);
    }, 5000);
}
