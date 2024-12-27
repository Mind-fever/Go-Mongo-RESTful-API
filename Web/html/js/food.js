let currentPage = 1;
const itemsPerPage = 5;
let foods = [];
document.addEventListener('DOMContentLoaded', function() {
    checkTokenAndInitialize(fetchFoods);
});

function checkTokenAndInitialize(callback) {
    const token = localStorage.getItem('authToken');
    const loginStatus = document.getElementById('loginStatus');
    const logoutButton = document.getElementById('logoutButton');

    if (!token) {
        alert('No estás logueado. Serás redirigido a la página de login.');
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

function fetchFoods() {
    makeRequest('http://localhost:8080/foods/', 'GET', null, 'application/json', true,
        (data) => {
            foods = data;
            renderFoods();
            renderPagination();
        },
        (status, error) => {
            console.error('Error:', error);
            const errorMessage = error.error || 'An unknown error occurred';
            foods = [];
            renderFoods();
            renderPagination();
            alert(`An error occurred: ${errorMessage}`);
        }
    );
}

function renderFoods() {
    const foodTableBody = document.getElementById('foodTable');
    foodTableBody.innerHTML = ''; // Clear existing rows
    const start = (currentPage - 1) * itemsPerPage;
    const end = start + itemsPerPage;
    const paginatedFoods = foods.slice(start, end);

    paginatedFoods.forEach(food => {
        const row = foodTableBody.insertRow();
        row.insertCell(0).textContent = food.name;
        row.insertCell(1).textContent = food.type;
        row.insertCell(2).textContent = food.price_per_unit;
        row.insertCell(3).textContent = food.current_quantity;
        row.insertCell(4).textContent = food.min_quantity;
        row.insertCell(5).textContent = food.meal_times.join(', ');

        // Add actions cell with Edit and Delete buttons
        const actionsCell = row.insertCell(6);
        const editButton = createEditButton(food);
        const deleteButton = createDeleteButton(food.id);
        actionsCell.appendChild(editButton);
        actionsCell.appendChild(deleteButton);
    });
}

function renderPagination() {
    pagination.innerHTML = '';
    const pageCount = Math.ceil(foods.length / itemsPerPage);

    const prevItem = createPaginationItem('Previous', () => {
        if (currentPage > 1) {
            currentPage--;
            renderFoods();
            renderPagination();
        }
    });
    pagination.appendChild(prevItem);

    for (let i = 1; i <= pageCount; i++) {
        const pageItem = createPaginationItem(i, () => {
            currentPage = i;
            renderFoods();
            renderPagination();
        });
        pagination.appendChild(pageItem);
    }

    const nextItem = createPaginationItem('Next', () => {
        if (currentPage < pageCount) {
            currentPage++;
            renderFoods();
            renderPagination();
        }
    });
    pagination.appendChild(nextItem);
}

function createEditButton(food) {
    const editButton = document.createElement('button');
    editButton.className = 'btn btn-primary';
    editButton.textContent = 'Edit';
    editButton.addEventListener('click', function() {
        const queryParams = new URLSearchParams({
            id: food.id,
            name: food.name,
            type: food.type,
            price_per_unit: food.price_per_unit,
            current_quantity: food.current_quantity,
            min_quantity: food.min_quantity,
            meal_times: food.meal_times.join(',')
        }).toString();
        window.location.href = `UpdateFood.html?${queryParams}`;
    });
    return editButton;
}

function createDeleteButton(foodId) {
    const deleteButton = document.createElement('button');
    deleteButton.className = 'btn btn-danger';
    deleteButton.textContent = 'Delete';
    deleteButton.setAttribute('data-id', foodId); // Set the food ID as a data attribute
    deleteButton.addEventListener('click', function() {
        deleteFood(foodId);
    });
    return deleteButton;
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

function deleteFood(foodId) {
    const url = `http://localhost:8080/foods/${foodId}`;
    makeRequest(url, 'DELETE', null, 'application/json', true,
        (data) => {
            console.log('Food deleted:', data);
            fetchFoods(); // Refresh the food list
        },
        (status, error) => {
            console.error('Error:', error);
            const errorMessage = error.error || 'An unknown error occurred';
            alert(`An error occurred: ${errorMessage}`);
        }
    );
}