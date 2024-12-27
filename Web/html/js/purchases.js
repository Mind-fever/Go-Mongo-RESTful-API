let currentPage = 1;
const itemsPerPage = 5;
let purchases = [];
document.addEventListener('DOMContentLoaded', function() {
    checkTokenAndInitialize(fetchPurchases);
});

function checkTokenAndInitialize(callback) {
    const token = localStorage.getItem('authToken');
    const loginStatus = document.getElementById('loginStatus');
    const logoutButton = document.getElementById('logoutButton');


    if (!token) {
        alert('You are not logged in. Redirecting to the login page.');
        window.location.href = 'login.html';
    } else {
        loginStatus.textContent = 'Logged in';
        logoutButton.style.display = 'inline-block';
        callback();
    }

    logoutButton.addEventListener('click', function() {
        localStorage.removeItem('authToken');
        window.location.href = 'login.html';
        loginStatus.textContent = 'Not logged in';
        logoutButton.style.display = 'none';
    });
}

function fetchPurchases() {
    const url = 'http://localhost:8080/purchases/';
    makeRequest(url, 'GET', null, 'application/json', true,
        (response) => {
            purchases = Array.isArray(response) ? response : response.purchases || [];
            renderPurchases();
            renderPagination();
        },
        (status, error) => {
            console.error('Error:', error);
            const errorMessage = error.error || 'An unknown error occurred';
            alert(`An error occurred: ${errorMessage}`);
        }
    );
}

function renderPurchases() {
    const purchaseList = document.getElementById('purchaseList');
    purchaseList.innerHTML = '';
    const start = (currentPage - 1) * itemsPerPage;
    const end = start + itemsPerPage;
    const paginatedPurchases = purchases.slice(start, end);

    paginatedPurchases.forEach(purchase => {
        const row = document.createElement('tr');
        const formattedDate = new Date(purchase.date).toLocaleString();
        row.innerHTML = `
            <td>${formattedDate}</td>
            <td>${purchase.total_cost}</td>
        `;
        purchaseList.appendChild(row);
    });
}

function renderPagination() {
    const pagination = document.getElementById('pagination');
    pagination.innerHTML = '';
    const pageCount = Math.ceil(purchases.length / itemsPerPage);

    const prevItem = createPaginationItem('Previous', () => {
        if (currentPage > 1) {
            currentPage--;
            renderPurchases();
            renderPagination();
        }
    });
    pagination.appendChild(prevItem);

    for (let i = 1; i <= pageCount; i++) {
        const pageItem = createPaginationItem(i, () => {
            currentPage = i;
            renderPurchases();
            renderPagination();
        });
        pagination.appendChild(pageItem);
    }

    const nextItem = createPaginationItem('Next', () => {
        if (currentPage < pageCount) {
            currentPage++;
            renderPurchases();
            renderPagination();
        }
    });
    pagination.appendChild(nextItem);
}

function createPaginationItem(text, onClick) {
    const pageItem = document.createElement('li');
    pageItem.className = 'page-item';
    const pageLink = document.createElement('a');
    pageLink.className = 'page-link';
    pageLink.textContent = text;
    pageLink.addEventListener('click', onClick);
    pageItem.appendChild(pageLink);
    return pageItem;
}