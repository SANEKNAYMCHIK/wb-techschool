document.addEventListener('DOMContentLoaded', () => {
    const searchForm = document.getElementById('searchForm');
    const searchBtn = document.getElementById('searchBtn');
    const orderUidInput = document.getElementById('orderUid');
    const resultDiv = document.getElementById('result');
    
    searchForm.addEventListener('submit', (e) => e.preventDefault());
    
    searchBtn.addEventListener('click', async () => {
        const orderUid = orderUidInput.value.trim();

        if (!orderUid) {
            resultDiv.textContent = 'Please enter Order UID';
            return;
        }
        
        try {
            resultDiv.innerHTML = '<div class="loading">Loading...</div>';
            const encodedUid = encodeURIComponent(orderUid);
            const response = await fetch(`/order/${encodedUid}`);
            
            if (response.status === 404) {
                const errorData = await response.json();
                throw new Error(errorData.error);
            }
            
            if (!response.ok) throw new Error(`Server error: ${response.status}`);
            const orderData = await response.json();
            displayOrder(orderData);
        } catch (error) {
            resultDiv.textContent = `Error: ${error.message}`;
        }
    });
    
    orderUidInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') searchBtn.click();
    });
    
    function displayOrder(order) {
        const formatDate = (dateString) => new Date(dateString).toLocaleString('en-US', {
            year: 'numeric', month: 'short', day: 'numeric',
            hour: '2-digit', minute: '2-digit'
        });
        
        const formatPrice = (price) => `${price.toFixed(2)} руб.`;
        
        let html = `<div class="order-header">
            <h2>Order Details</h2>
            <div><strong>Customer:</strong> ${order.delivery.name}</div>
            <div><strong>Phone:</strong> ${order.delivery.phone}</div>
            <div><strong>Address:</strong> ${order.delivery.city}, ${order.delivery.address}</div>
            <div><strong>Order UID:</strong> ${order.order_uid}</div>
            <div><strong>Track Number:</strong> ${order.track_number}</div>
            <div><strong>Date Created:</strong> ${formatDate(order.date_created)}</div>
            <div><strong>Total Amount:</strong> ${formatPrice(order.payment.amount)}</div>
        </div>`;
        
        html += `<div class="order-section">
            <h3>Order Items (${order.items.length})</h3>`;
        
        order.items.forEach(item => {
            html += `<div class="item">
                <div><strong>${item.name}</strong></div>
                <div>Brand: ${item.brand}</div>
                <div>Price: ${formatPrice(item.price)}</div>
                <div>Status: ${getStatusText(item.status)}</div>
            </div>`;
        });

        resultDiv.innerHTML = html;
    }
    
    function getStatusText(statusCode) {
        const statusMap = {202: 'Processing', 203: 'Shipped', 204: 'Delivered', 205: 'Cancelled'};
        return statusMap[statusCode] || `Status: ${statusCode}`;
    }
});
