document.addEventListener('DOMContentLoaded', async () => {
    renderNav(document.getElementById('nav-buttons'));

    try {
        await apiRequest('/me', 'GET', null, true); // check auth
    } catch (error) {
        console.error('JWT expired or user not authenticated.');
        window.location.href = 'auth.html';
        return;
    }

    await loadCreatedAuctions();
});

async function loadCreatedAuctions() {
    try {
        const auctions = await apiRequest('/auctions/created-auctions', 'GET', null, true);
        const gridContainer = document.getElementById('auctions-grid-container'); 

        if (!auctions.length) {
            gridContainer.innerHTML = '<p class="col-span-full text-center text-gray-600">You have not created any auctions yet.</p>';
            return;
        }

        gridContainer.innerHTML = auctions.map(auction => `
            <div id="auction-${auction.id}" class="border border-gray-200 bg-white p-5 rounded-lg shadow-sm flex flex-col h-full relative"> ${auction.image_path ? `<img src="${auction.image_path}" alt="Auction Image" class="w-full h-48 object-cover rounded-md mb-3 flex-shrink-0">` : ''}
                <h2 class="text-xl font-bold text-gray-800 mb-2">${auction.title}</h2>
                <p class="text-gray-600 mb-3 flex-grow overflow-hidden text-ellipsis">${auction.description}</p>
                <p class="text-gray-700 mb-1"><strong>Current Price:</strong> ₦${auction.current_price ? auction.current_price.toLocaleString() : '0.00'}</p>
                <p class="text-gray-700 mb-1"><strong>Status:</strong> <span class="${auction.status === 'active' ? 'text-green-600' : (auction.status === 'closed' ? 'text-red-600' : 'text-yellow-600')} font-semibold">${auction.status.toUpperCase()}</span></p>
                <p class="text-gray-700 mb-1"><strong>Type:</strong> ${auction.type}</p>
                <p class="text-gray-700 text-sm mb-4"><strong>Start:</strong> ${new Date(auction.start_time).toLocaleString()} | <strong>End:</strong> ${new Date(auction.end_time).toLocaleString()}</p>
                
                <div class="mt-auto flex flex-wrap justify-center gap-2"> 
                    <button class="w-full sm:w-auto flex-grow bg-red-500 text-white px-3 py-2 rounded hover:bg-red-600 text-sm md:text-base" 
                            onclick="handleCloseAuction('${auction.id}')"
                            ${auction.status === 'closed' ? 'disabled' : ''}>
                        ${auction.status === 'closed' ? 'Closed' : 'Close Auction'}
                    </button>
                    <button class="w-full sm:w-auto flex-grow bg-yellow-500 text-white px-3 py-2 rounded hover:bg-yellow-600 text-sm md:text-base" 
                            onclick="handleUpdateAuction('${auction.id}', '${encodeURIComponent(JSON.stringify(auction))}')">
                        Update Auction
                    </button>
                    <button class="w-full sm:w-auto flex-grow bg-gray-500 text-white px-3 py-2 rounded hover:bg-gray-600 text-sm md:text-base" 
                            onclick="handleDeleteAuction('${auction.id}')">
                        Delete Auction
                    </button>
                </div>
                
                <div id="update-form-${auction.id}" 
                     class="absolute inset-0 bg-white p-5 rounded-lg shadow-lg overflow-y-auto transform transition-transform duration-300 translate-y-full opacity-0 z-10">
                </div>
                <div id="message-${auction.id}" class="mt-3 text-center text-sm"></div>
            </div>
        `).join('');

    } catch (error) {
        console.error('Failed to load created auctions:', error);
        document.getElementById('auctions-grid-container').innerHTML = `<p class="col-span-full text-red-600 text-center">Error loading your created auctions.</p>`;
    }
}

// Helper function to format date for input
function formatDateForInput(dateString) {
    const date = new Date(dateString);
    if (isNaN(date)) {
        return ''; // handle invalid date if needed
    }
    return date.toISOString().split('T')[0];
}

// Function to handle Close Auction button click
async function handleCloseAuction(auctionId) {
    const messageDiv = document.getElementById(`message-${auctionId}`); // Target specific message div
    messageDiv.textContent = 'Closing auction...';
    messageDiv.className = 'text-blue-600 mt-3 text-sm';
    try {
        const result = await apiRequest(`/auctions/${auctionId}/close`, 'POST', null, true);
        messageDiv.innerHTML = `
            ✅ Auction Closed!<br>
            Winning Bid: ₦${result.winning_bid ? result.winning_bid.toLocaleString() : 'N/A'} <br>
            Status: ${result.status}
        `;
        messageDiv.className = 'text-green-600 mt-3 text-sm';
        await loadCreatedAuctions(); // Refresh the list to update status and button state
    } catch (error) {
        console.error('Failed to close auction:', error);
        messageDiv.textContent = `❌ Failed to close auction: ${error.message}`;
        messageDiv.className = 'text-red-600 mt-3 text-sm';
    }
}

// Function to handle Update Auction button click
function handleUpdateAuction(auctionId, auctionString) {
    const auction = JSON.parse(decodeURIComponent(auctionString));
    const formDiv = document.getElementById(`update-form-${auctionId}`); // Target specific update form div
    const messageDiv = document.getElementById(`message-${auctionId}`); // Target specific message div

    // Hide any other active update forms AND reset their transform/opacity
    document.querySelectorAll('[id^="update-form-"]').forEach(form => {
        if (form.id !== `update-form-${auctionId}`) {
            form.classList.remove('translate-y-0', 'opacity-100');
            form.classList.add('translate-y-full', 'opacity-0');
            // Give it a moment to transition out before clearing content
            setTimeout(() => form.innerHTML = '', 300); 
        }
    });
    // Clear messages from other auctions
    document.querySelectorAll('[id^="message-"]').forEach(msg => {
        if (msg.id !== `message-${auctionId}`) {
            msg.textContent = '';
        }
    });

    // Populate the form content
    formDiv.innerHTML = `
        <h3 class="text-lg font-bold text-gray-700 mb-3">Update Auction: ${auction.title}</h3>
        <form id="auction-update-form-${auction.id}" class="space-y-3 text-gray-700">
            <div>
                <label for="update-title-${auction.id}" class="block text-sm font-medium">Title</label>
                <input type="text" id="update-title-${auction.id}" value="${auction.title}" 
                       placeholder="Title" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500">
            </div>
            <div>
                <label for="update-description-${auction.id}" class="block text-sm font-medium">Description</label>
                <textarea id="update-description-${auction.id}" placeholder="Description" 
                          class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500">${auction.description}</textarea>
            </div>
            <div>
                <label for="update-starting-price-${auction.id}" class="block text-sm font-medium">Starting Price</label>
                <input type="number" id="update-starting-price-${auction.id}" value="${auction.starting_price}" 
                       placeholder="Starting Price" step="0.01"
                       class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500">
            </div>
            <div>
                <label for="update-start-time-${auction.id}" class="block text-sm font-medium">Start Date</label>
                <input type="date" id="update-start-time-${auction.id}" value="${formatDateForInput(auction.start_time)}" 
                       class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500">
            </div>
            <div>
                <label for="update-end-time-${auction.id}" class="block text-sm font-medium">End Date</label>
                <input type="date" id="update-end-time-${auction.id}" value="${formatDateForInput(auction.end_time)}" 
                       class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500">
            </div>
            <div>
                <label for="update-auction-type-${auction.id}" class="block text-sm font-medium">Auction Type</label>
                <select id="update-auction-type-${auction.id}" 
                        class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500">
                    <option value="english" ${auction.type === 'english' ? 'selected' : ''}>English Auction</option>
                    <option value="dutch" ${auction.type === 'dutch' ? 'selected' : ''}>Dutch Auction</option>
                    <option value="sealed" ${auction.type === 'sealed' ? 'selected' : ''}>Sealed Bid Auction</option>
                </select>
            </div>
            <div class="flex justify-end space-x-2 pt-2">
                <button type="submit" class="bg-green-600 text-white px-5 py-2 rounded-md hover:bg-green-700 text-sm font-semibold">Submit Update</button>
                <button type="button" class="bg-gray-400 text-white px-5 py-2 rounded-md hover:bg-gray-500 text-sm font-semibold" 
                        onclick="handleCancelUpdate('${auction.id}')">Cancel</button>
            </div>
        </form>
    `;

    // Show the form with a slight delay to allow rendering
    setTimeout(() => {
        formDiv.classList.remove('translate-y-full', 'opacity-0');
        formDiv.classList.add('translate-y-0', 'opacity-100');
    }, 10); // Small delay

    document.getElementById(`auction-update-form-${auction.id}`).addEventListener('submit', async (e) => {
        e.preventDefault();
        await submitAuctionUpdate(auction.id, auction.image_path);
    });
}

// New function to handle Cancel Update
function handleCancelUpdate(auctionId) {
    const formDiv = document.getElementById(`update-form-${auctionId}`);
    const messageDiv = document.getElementById(`message-${auctionId}`);
    
    formDiv.classList.remove('translate-y-0', 'opacity-100');
    formDiv.classList.add('translate-y-full', 'opacity-0');
    
    // Clear content after transition
    setTimeout(() => {
        formDiv.innerHTML = ''; 
        messageDiv.textContent = '';
    }, 300); // Match transition duration
}

// Function to submit Auction Update
async function submitAuctionUpdate(auctionId, imagePath) {
    const messageDiv = document.getElementById(`message-${auctionId}`);
    messageDiv.textContent = 'Updating auction...';
    messageDiv.className = 'text-blue-600 mt-3 text-sm';

    const payload = {
        title: document.getElementById(`update-title-${auctionId}`).value,
        description: document.getElementById(`update-description-${auctionId}`).value,
        starting_price: parseFloat(document.getElementById(`update-starting-price-${auctionId}`).value),
        type: document.getElementById(`update-auction-type-${auctionId}`).value,
        start_time: document.getElementById(`update-start-time-${auctionId}`).value,
        end_time: document.getElementById(`update-end-time-${auctionId}`).value,
        image_path: imagePath // Use the existing image path
    };

    try {
        await apiRequest(`/auctions/${auctionId}`, 'PUT', payload, true);
        messageDiv.textContent = '✅ Auction updated successfully!';
        messageDiv.className = 'text-green-600 mt-3 text-sm';
        
        // Hide and clear the form after successful update
        handleCancelUpdate(auctionId); // Reuse the cancel logic
        await loadCreatedAuctions(); // Refresh the list
    } catch (error) {
        console.error('Failed to update auction:', error);
        messageDiv.textContent = `❌ Failed to update auction: ${error.message}`;
        messageDiv.className = 'text-red-600 mt-3 text-sm';
    }
}

// Function to handle Delete Auction button click (improved)
async function handleDeleteAuction(auctionId) {
    const confirmed = await showCustomConfirm('Are you sure you want to delete this auction? This action cannot be undone.');
    if (!confirmed) {
        return; // User cancelled
    }

    const messageDiv = document.getElementById(`message-${auctionId}`); // Target specific message div
    messageDiv.textContent = 'Deleting auction...';
    messageDiv.className = 'text-blue-600 mt-3 text-sm';

    try {
        await apiRequest(`/auctions/${auctionId}`, 'DELETE', null, true);
        messageDiv.textContent = '✅ Auction deleted successfully!';
        messageDiv.className = 'text-green-600 mt-3 text-sm';
        // Remove the deleted auction's card from the DOM without a full refresh
        document.getElementById(`auction-${auctionId}`).remove();
    } catch (error) {
        console.error('Failed to delete auction:', error);
        messageDiv.textContent = `❌ Failed to delete auction: ${error.message}`;
        messageDiv.className = 'text-red-600 mt-3 text-sm';
    }
}

// Custom confirmation dialog function
function showCustomConfirm(message) {
    return new Promise((resolve) => {
        const dialog = document.createElement('div');
        dialog.className = 'fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50';
        dialog.innerHTML = `
            <div class="bg-white p-6 rounded-lg shadow-xl max-w-sm w-full">
                <p class="text-lg font-semibold mb-4 text-gray-800">${message}</p>
                <div class="flex justify-end space-x-3">
                    <button id="confirm-cancel" class="bg-gray-300 text-gray-800 px-4 py-2 rounded hover:bg-gray-400">Cancel</button>
                    <button id="confirm-ok" class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600">Delete</button>
                </div>
            </div>
        `;
        document.body.appendChild(dialog);

        document.getElementById('confirm-ok').addEventListener('click', () => {
            dialog.remove();
            resolve(true);
        });

        document.getElementById('confirm-cancel').addEventListener('click', () => {
            dialog.remove();
            resolve(false);
        });
    });
}