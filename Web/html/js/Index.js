document.addEventListener('DOMContentLoaded', function() {
    const token = localStorage.getItem('authToken');
    if (!token) {
        window.location.href = 'login.html';
        alert('No estás logueado. Serás redirigido a la página de login.');
    } else {
        const loginStatus = document.getElementById('loginStatus');
        const logoutButton = document.getElementById('logoutButton');
        loginStatus.textContent = 'Logged in';
        logoutButton.style.display = 'inline-block';

        logoutButton.addEventListener('click', function() {
            localStorage.removeItem('authToken');
            window.location.href = 'login.html';
            loginStatus.textContent = 'Not logged in';
            logoutButton.style.display = 'none';
        });
    }
});