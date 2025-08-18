document.addEventListener('DOMContentLoaded', () => {
    const searchBtn = document.getElementById('searchBtn');
    const orderUidInput = document.getElementById('orderUid');
    const resultDiv = document.getElementById('result');
    searchBtn.addEventListener('click', async () => {
        const orderUid = orderUidInput.value.trim();

        if (!orderUid) {
            resultDiv.textContent = 'Please enter Order UID';
            return;
        }
        
        try {
            resultDiv.textContent = 'Loading...';
            
            const encodedUid = encodeURIComponent(orderUid);
            const response = await fetch(`/order/${encodedUid}`);
            
            if (response.status === 404) {
                const errorData = await response.json();
                throw new Error(errorData.error);
            }
            
            if (!response.ok) {
                throw new Error(`Server returned status: ${response.status}`);
            }
            
            const orderData = await response.json();
            resultDiv.textContent = JSON.stringify(orderData, null, 2);
        } catch (error) {
            resultDiv.textContent = `Error: ${error.message}`;
        }
    });
    
    orderUidInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            searchBtn.click();
        }
    });
});