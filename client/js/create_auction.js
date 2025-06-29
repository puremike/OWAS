document.getElementById('image-upload-form').addEventListener('submit', async (e) => {
    e.preventDefault();
  
    const formData = new FormData();
    const imageFile = document.getElementById('image').files[0];
    formData.append('image', imageFile);
  
    try {
      const response = await fetch(`${API_BASE_URL}/auctions/image_upload`, {
        method: 'POST',
        body: formData,
        credentials: 'include'
      });
  
      if (!response.ok) {
        throw new Error('Image upload failed');
      }
  
      const result = await response.json();
      document.getElementById('image-upload-message').textContent = '✅ Image uploaded successfully!';
      document.getElementById('image-upload-message').className = 'text-green-600 mt-2';
  
      // Save the uploaded image filename/path to hidden field
      document.getElementById('image_path').value = result.image_path;
  
    } catch (error) {
      document.getElementById('image-upload-message').textContent = '❌ ' + error.message;
      document.getElementById('image-upload-message').className = 'text-red-600 mt-2';
    }
  });
  

  document.getElementById('create-auction-form').addEventListener('submit', async (e) => {
    e.preventDefault();
  
    const auctionData = {
      title: document.getElementById('title').value,
      description: document.getElementById('description').value,
      starting_price: parseFloat(document.getElementById('starting_price').value),
      type: document.getElementById('auction_type').value,
      start_time: document.getElementById('start_time').value,
      end_time: document.getElementById('end_time').value,
      image_path: document.getElementById('image_path').value
    };
  
    try {
      const response = await fetch(`${API_BASE_URL}/auctions`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify(auctionData)
      });
  
      if (!response.ok) {
        throw new Error('Auction creation failed');
      }
  
      document.getElementById('create-message').textContent = '✅ Auction created successfully!';
      document.getElementById('create-message').className = 'text-green-600 mt-2';
      document.getElementById('create-auction-form').reset();
  
    } catch (error) {
      document.getElementById('create-message').textContent = '❌ ' + error.message;
      document.getElementById('create-message').className = 'text-red-600 mt-2';
    }
  });
  
















