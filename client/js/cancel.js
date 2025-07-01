document.addEventListener('DOMContentLoaded', () => {
    const params = new URLSearchParams(window.location.search);
    const sessionId = params.get('session_id');
    console.log('Canceled session ID:', sessionId);
});
