document.addEventListener('DOMContentLoaded', function() {
    checkTokenAndInitialize(initializeNewRecipeForm);
});

function checkTokenAndInitialize(callback) {
    const token = localStorage.getItem('authToken');
    const loginStatus = document.getElementById('loginStatus');
    const logoutButton = document.getElementById('logoutButton');

    if (!token) {
        alert('No estás logueado. Serás redirigido a la página de login.');
        window.location.href = 'login.html';
    } else {
        loginStatus.textContent = 'Estás logueado.';
        logoutButton.style.display = 'inline-block';
        callback();
    }

    logoutButton.addEventListener('click', function() {
        localStorage.removeItem('authToken');
        window.location.href = 'login.html';
        loginStatus.textContent = 'No estás logueado.';
        logoutButton.style.display = 'none';
    });
}

function initializeNewRecipeForm() {
    const recipeForm = document.getElementById('recipeForm');
    const availableFoods = document.getElementById('availableFoods');
    const ingredientList = document.getElementById('ingredientList');
    const addFoodButton = document.getElementById('addFood');
    const removeFoodButton = document.getElementById('removeFood');
    const pagination = document.getElementById('pagination');
    const itemsPerPage = 5;
    let currentPage = 1;
    let foods = [];

    fetchFoods();

    addFoodButton.addEventListener('click', handleAddFood);
    removeFoodButton.addEventListener('click', handleRemoveFood);

    if (recipeForm) {
        recipeForm.addEventListener('submit', handleSubmit);
    }

    function handleAddFood() {
        const selectedItems = Array.from(availableFoods.querySelectorAll('.list-group-item.active'));
        selectedItems.forEach(item => {
            const row = ingredientList.insertRow();
            const foodIdCell = row.insertCell(0);
            const hiddenFoodIdInput = document.createElement('input');
            hiddenFoodIdInput.type = 'hidden';
            hiddenFoodIdInput.value = item.dataset.foodId;
            foodIdCell.appendChild(hiddenFoodIdInput);
            foodIdCell.appendChild(document.createTextNode(item.textContent.split(' (')[0]));
            const quantityCell = row.insertCell(1);
            const quantityInput = document.createElement('input');
            quantityInput.type = 'number';
            quantityInput.className = 'form-control';
            quantityInput.value = 1;
            quantityCell.appendChild(quantityInput);
            const actionsCell = row.insertCell(2);
            const removeButton = createRemoveButton(hiddenFoodIdInput.value);
            actionsCell.appendChild(removeButton);
            item.remove();
        });
    }

    function handleRemoveFood() {
        const selectedRows = Array.from(ingredientList.querySelectorAll('tr'));
        selectedRows.forEach(row => {
            const foodId = row.querySelector('input[type="hidden"]').value;
            const food = foods.find(f => f.id === foodId);
            if (food) {
                const listItem = createFoodListItem(food);
                availableFoods.appendChild(listItem);
            }
            row.remove();
        });
    }

    function handleSubmit(event) {
        event.preventDefault();
        const formData = new FormData(recipeForm);
        const ingredients = Array.from(ingredientList.querySelectorAll('tr')).map(row => {
            const hiddenInput = row.querySelector('input[type="hidden"]');
            if (hiddenInput) {
                return {
                    food_id: hiddenInput.value,
                    quantity: parseFloat(row.cells[1].querySelector('input').value)
                };
            }
            return null;
        }).filter(ingredient => ingredient !== null);

        const recipeData = {
            name: formData.get('recipeName'),
            meal_time: formData.get('mealTime'),
            ingredients: ingredients
        };

        makeRequest('http://localhost:8080/recipes/', 'POST', recipeData, 'application/json', true,
            (data) => {
                console.log('Success:', data);
                alert('Recipe added successfully!');
                window.location.href = 'recipes.html';
            },
            (status, error) => handleError(error)
        );
    }

    function fetchFoods() {
        makeRequest('http://localhost:8080/foods/', 'GET', null, 'application/json', true,
            (data) => {
                foods = data;
                renderFoods();
                renderPagination();
            },
            (status, error) => handleError(error)
        );
    }

    function renderFoods() {
        availableFoods.innerHTML = '';
        const start = (currentPage - 1) * itemsPerPage;
        const end = start + itemsPerPage;
        const paginatedFoods = foods.slice(start, end);

        paginatedFoods.forEach(food => {
            const listItem = createFoodListItem(food);
            availableFoods.appendChild(listItem);
        });
    }

    function renderPagination() {
        pagination.innerHTML = '';
        const pageCount = Math.ceil(foods.length / itemsPerPage);

        for (let i = 1; i <= pageCount; i++) {
            const pageItem = createPaginationItem(i);
            pagination.appendChild(pageItem);
        }
    }

    function createRemoveButton(foodId) {
        const removeButton = document.createElement('button');
        removeButton.className = 'btn btn-danger';
        removeButton.textContent = 'Remove';
        removeButton.addEventListener('click', function() {
            const food = foods.find(f => f.id === foodId);
            if (food) {
                const listItem = createFoodListItem(food);
                availableFoods.appendChild(listItem);
            }
            this.closest('tr').remove();
        });
        return removeButton;
    }

    function createFoodListItem(food) {
        const listItem = document.createElement('li');
        listItem.className = 'list-group-item';
        listItem.dataset.foodId = food.id;
        listItem.textContent = `${food.name} (Stock: ${food.current_quantity}, Meal Times: ${food.meal_times.join(', ')})`;
        listItem.addEventListener('click', function() {
            listItem.classList.toggle('active');
        });
        return listItem;
    }

    function createPaginationItem(pageNumber) {
        const pageItem = document.createElement('li');
        pageItem.className = 'page-item';
        const pageLink = document.createElement('a');
        pageLink.className = 'page-link';
        pageLink.textContent = pageNumber;
        pageLink.addEventListener('click', function() {
            currentPage = pageNumber;
            renderFoods();
        });
        pageItem.appendChild(pageLink);
        return pageItem;
    }

    function handleError(error, reject) {
        console.error('Error:', error);
        const errorMessage = error.error || 'An unknown error occurred';
        alert(`An error occurred: ${errorMessage}`);
        if (reject) reject(error);
    }
}