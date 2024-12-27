document.addEventListener('DOMContentLoaded', function() {
    checkTokenAndInitialize(initializeNewPurchaseForm);
});

function checkTokenAndInitialize(callback) {
    const token = localStorage.getItem('authToken');
    const loginStatus = document.getElementById('loginStatus');
    const logoutButton = document.getElementById('logoutButton');

    if (!token) {
        alert('You are not logged in. You will be redirected to the login page.');
        window.location.href = 'login.html';
    } else {
        loginStatus.textContent = 'You are logged in.';
        logoutButton.style.display = 'inline-block';
        callback();
    }

    logoutButton.addEventListener('click', function() {
        localStorage.removeItem('authToken');
        loginStatus.textContent = 'You are not logged in.';
        window.location.href = 'login.html';
    });
}

function initializeNewPurchaseForm() {
    const filterForm = document.getElementById('filterForm');
    const createPurchaseBtn = document.getElementById('createPurchaseBtn');

    filterForm.addEventListener('submit', function(event) {
        event.preventDefault();
        fetchLowStockFoods();
    });

    createPurchaseBtn.addEventListener('click', function() {
        createPurchase();
    });

    fetchLowStockFoods();
}

function fetchLowStockFoods() {
    const filterForm = document.getElementById('filterForm');
    const formData = new FormData(filterForm);
    const name = formData.get('name');
    const type = formData.get('type');
    let url = 'http://localhost:8080/stock/';

    if (name || (type && type !== '')) {
        url += '?';
        if (name) url += `name=${encodeURIComponent(name)}&`;
        if (type && type !== '') url += `type=${encodeURIComponent(type)}`;
    }

    makeRequest(url, 'GET', null, 'application/json', true,
        (foods) => {
            renderLowStockFoods(foods);
        },
        (status, error) => {
            console.error('Error:', error);
            renderLowStockFoods([]);
            const errorMessage = error.error || 'An unknown error occurred';
            alert(`An error occurred: ${errorMessage}`);
        }
    );
}

function renderLowStockFoods(foods) {
    const productList = document.getElementById('productList');
    productList.innerHTML = '';
    if (!Array.isArray(foods)) {
        console.error("Expected an array of foods, but got:", foods);
        return; // Exit if data is not valid
    }
    foods.forEach(food => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${food.name}</td>
            <td>${food.type}</td>
            <td>${food.current_quantity}</td>
            <td>${food.min_quantity}</td>
            <td>${food.price_per_unit}</td>
        `;
        productList.appendChild(row);
    });
}

function createPurchase() {
    makeRequest('http://localhost:8080/purchases/', 'POST', {}, 'application/json', true,
        (data) => {
            alert('Purchase created successfully');
            window.location.href = 'purchases.html'; // Redirect to purchases.html
        },
        (status, error) => {
            console.error('Error:', error);
            const errorMessage = error.error || 'An unknown error occurred';
            alert(`An error occurred: ${errorMessage}`);
        }
    );
}