document.addEventListener('DOMContentLoaded', async () => {
    const params = new URLSearchParams(window.location.search);
    const sessionId = params.get('session_id');

    if (sessionId) {
        try {
            const res = await fetch(`${API_BASE_URL}/stripe/session/${sessionId}`, { credentials: 'include' });
            const session = await res.json();

            document.getElementById('payment-info').innerHTML = `
                <p><strong>Amount Paid:</strong> $${(session.amount_total / 100).toFixed(2)}</p>
                <p><strong>Status:</strong> ${session.payment_status}</p>
                <p><strong>Amount Paid:</strong> ${(session.currency)}</p>
            `;
        } catch (error) {
            document.getElementById('payment-info').textContent = 'Could not load payment details.';
        }
    }

    setTimeout(() => {
        window.location.href ='my_auctions.html';
    }, 3000);
});

