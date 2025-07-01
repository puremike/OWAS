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
    alert(`New Notification: ${message}`);
}
