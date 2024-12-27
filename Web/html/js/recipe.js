let currentPage = 1;
const itemsPerPage = 4;
let recipes = [];
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
            fetchRecipes(); // Fetch recipes after foods are loaded
        },
        (status, error) => {
            console.error('Error:', error);
            const errorMessage = error.error || 'An unknown error occurred';
            alert(`An error occurred: ${errorMessage}`);
        }
    );
}

function fetchRecipes(filterParams = {}) {
    const queryString = new URLSearchParams(filterParams).toString();
    makeRequest(`http://localhost:8080/recipes/filter?${queryString}`, 'GET', null, 'application/json', true,
        (data) => {
            recipes = data;
            renderRecipes();
            renderPagination();
        },
        (status, error) => {
            console.error('Error:', error);
            recipes = []; // Clear the recipes list
            renderRecipes();
            renderPagination();
            const errorMessage = error.error || 'An unknown error occurred';
            alert(`An error occurred: ${errorMessage}`);
        }
    );
}

function renderRecipes() {
    const recipeList = document.getElementById('recipeList');
    recipeList.innerHTML = ''; // Clear the list before rendering

    if (recipes.length === 0) {
        // If no recipes, clear the list and return
        const pagination = document.getElementById('pagination');
        pagination.innerHTML = '';
        return;
    }

    const start = (currentPage - 1) * itemsPerPage;
    const end = start + itemsPerPage;
    const paginatedRecipes = recipes.slice(start, end);

    paginatedRecipes.forEach(recipe => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${recipe.name}</td>
            <td>${recipe.meal_time}</td>
            <td>
                <ul>
                    ${recipe.ingredients.map(ingredient => {
            const food = foods.find(f => f.id === ingredient.food_id);
            return `<li data-id="${ingredient.food_id}">${food ? food.name : ingredient.food_id}: ${ingredient.quantity}</li>`;
        }).join('')}
                </ul>
            </td>
            <td></td>
        `;
        const actionsCell = row.querySelector('td:last-child');
        actionsCell.appendChild(createEditButton(recipe));
        actionsCell.appendChild(createDeleteButton(recipe.id));
        recipeList.appendChild(row);
    });
}

function renderPagination() {
    const pagination = document.getElementById('pagination');
    pagination.innerHTML = ''; // Clear pagination before rendering

    if (recipes.length === 0) {
        // If no recipes, clear the pagination and return
        return;
    }

    const pageCount = Math.ceil(recipes.length / itemsPerPage);

    const prevItem = createPaginationItem('Previous', () => {
        if (currentPage > 1) {
            currentPage--;
            renderRecipes();
            renderPagination();
        }
    });
    pagination.appendChild(prevItem);

    for (let i = 1; i <= pageCount; i++) {
        const pageItem = createPaginationItem(i, () => {
            currentPage = i;
            renderRecipes();
            renderPagination();
        });
        pagination.appendChild(pageItem);
    }

    const nextItem = createPaginationItem('Next', () => {
        if (currentPage < pageCount) {
            currentPage++;
            renderRecipes();
            renderPagination();
        }
    });
    pagination.appendChild(nextItem);
}

function createEditButton(recipe) {
    const editButton = document.createElement('button');
    editButton.className = 'btn btn-primary';
    editButton.textContent = 'Edit';
    editButton.setAttribute('data-id', recipe.id);
    editButton.addEventListener('click', function() {
        const recipeId = this.getAttribute('data-id');
        editRecipe(recipeId);
    });
    return editButton;
}

function createDeleteButton(recipeId) {
    const deleteButton = document.createElement('button');
    deleteButton.className = 'btn btn-danger';
    deleteButton.textContent = 'Delete';
    deleteButton.setAttribute('data-id', recipeId); // Set the recipe ID as a data attribute
    deleteButton.addEventListener('click', function() {
        deleteRecipe(recipeId);
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

function deleteRecipe(id) {
    makeRequest(`http://localhost:8080/recipes/${id}`, 'DELETE', null, 'application/json', true,
        (data) => {
            if (data && data.success) {
                console.log('Recipe deleted successfully');
            } else {
                console.error('Error deleting recipe:', data ? data.error : 'No response data');
            }
            filterRecipes(); // Refresh the table after deletion
        },
        (status, error) => handleError(error)
    );
}

function filterRecipes() {
    const mealTime = document.getElementById('filterMealTime').value;
    const productType = document.getElementById('type').value;
    const productName = document.getElementById('filterProductName').value;

    const filterParams = {};

    if (mealTime && mealTime !== 'all') {
        filterParams.meal_time = mealTime;
    }
    if (productType) {
        filterParams.product_type = productType;
    }
    if (productName) {
        filterParams.product_name = productName;
    }

    fetchRecipes(filterParams);
}

function editRecipe(id) {
    window.location.href = `UpdateRecipe.html?id=${id}`;
}

function handleError(error, reject) {
    console.error('Error:', error);
    const errorMessage = error.error || 'An unknown error occurred';
    alert(`An error occurred: ${errorMessage}`);
    if (reject) reject(error);
}

// Expose functions to global scope for HTML to call
window.filterRecipes = filterRecipes;
window.deleteRecipe = deleteRecipe;