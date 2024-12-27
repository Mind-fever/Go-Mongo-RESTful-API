document.addEventListener('DOMContentLoaded', function() {
    const recipeForm = document.getElementById('recipeForm');
    const recipeNameInput = document.getElementById('recipeName');
    const mealTimeSelect = document.getElementById('mealTime');
    const recipeFoodsList = document.getElementById('recipeFoods');
    const userFoodsList = document.getElementById('userFoods');

    const urlParams = new URLSearchParams(window.location.search);
    const recipeId = urlParams.get('id');

    let foods = [];
    let recipeFoods = [];

    checkTokenAndInitialize(initializePage);

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

    function initializePage() {
        fetchFoods()
            .then(() => fetchRecipeById(recipeId))
            .then(recipe => {
                recipeNameInput.value = recipe.name;
                mealTimeSelect.value = recipe.meal_time;
                recipeFoods = recipe.ingredients;
                renderFoods();
            })
            .catch(error => console.error('Error fetching data:', error));

        recipeForm.addEventListener('submit', handleSubmit);
        userFoodsList.addEventListener('click', handleUserFoodClick);
        recipeFoodsList.addEventListener('click', handleRecipeFoodClick);
    }

    function fetchRecipeById(id) {
        return new Promise((resolve, reject) => {
            makeRequest(`http://localhost:8080/recipes/${id}`, Method.GET, null, ContentType.JSON, CallType.PRIVATE,
                (data) => resolve(data),
                (status, error) => handleError(error, reject)
            );
        });
    }

    function fetchFoods() {
        return new Promise((resolve, reject) => {
            makeRequest('http://localhost:8080/foods/', Method.GET, null, ContentType.JSON, CallType.PRIVATE,
                (data) => {
                    foods = data;
                    resolve(data);
                },
                (status, error) => handleError(error, reject)
            );
        });
    }

    function renderFoods() {
        userFoodsList.innerHTML = '';
        recipeFoodsList.innerHTML = '';

        foods.forEach(food => {
            const li = createFoodListItem(food);
            userFoodsList.appendChild(li);
        });

        recipeFoods.forEach(ingredient => {
            const food = foods.find(f => f.id === ingredient.food_id);
            const li = createRecipeFoodListItem(food, ingredient);
            recipeFoodsList.appendChild(li);
        });
    }

    function createFoodListItem(food) {
        const li = document.createElement('li');
        li.className = 'list-group-item';
        li.textContent = `${food.name} (${food.current_quantity})`;
        li.dataset.foodId = food.id;
        return li;
    }

    function createRecipeFoodListItem(food, ingredient) {
        const li = document.createElement('li');
        li.className = 'list-group-item d-flex justify-content-between align-items-center';
        li.innerHTML = `${food ? food.name : ingredient.food_id}:
            <input type="number" value="${ingredient.quantity}" class="form-control ingredient-quantity" data-food-id="${ingredient.food_id}" style="width: 80px; margin-left: 10px;">`;
        li.dataset.foodId = ingredient.food_id;
        return li;
    }

    function handleUserFoodClick(event) {
        if (event.target && event.target.matches('li.list-group-item')) {
            userFoodsList.querySelectorAll('li').forEach(li => li.classList.remove('active'));
            event.target.classList.add('active');
        }
    }

    function handleRecipeFoodClick(event) {
        if (event.target && event.target.matches('li.list-group-item')) {
            recipeFoodsList.querySelectorAll('li').forEach(li => li.classList.remove('active'));
            event.target.classList.add('active');
        }
    }

    function addFood() {
        const selectedFood = userFoodsList.querySelector('.list-group-item.active');
        if (selectedFood) {
            const foodId = selectedFood.dataset.foodId;
            const food = foods.find(f => f.id === foodId);
            const quantity = 1;
            recipeFoods.push({ food_id: foodId, quantity: parseFloat(quantity) });
            renderFoods();
        }
    }

    function removeFood() {
        const selectedFood = recipeFoodsList.querySelector('.list-group-item.active');
        if (selectedFood) {
            const foodId = selectedFood.dataset.foodId;
            recipeFoods = recipeFoods.filter(ingredient => ingredient.food_id !== foodId);
            renderFoods();
        }
    }

    function handleSubmit(event) {
        event.preventDefault();
        const updatedRecipe = {
            id: recipeId,
            name: recipeNameInput.value,
            meal_time: mealTimeSelect.value,
            ingredients: recipeFoods.map(ingredient => {
                const quantityInput = document.querySelector(`.ingredient-quantity[data-food-id="${ingredient.food_id}"]`);
                return {
                    food_id: ingredient.food_id,
                    quantity: parseFloat(quantityInput.value)
                };
            })
        };

        makeRequest(`http://localhost:8080/recipes/${recipeId}`, Method.PUT, updatedRecipe, ContentType.JSON, CallType.PRIVATE,
            (data) => {
                alert('Recipe updated successfully');
                window.location.href = 'recipes.html';
            },
            (status, error) => handleError(error)
        );
    }

    function handleError(error, reject) {
        console.error('Error:', error);
        const errorMessage = error.error || 'An unknown error occurred';
        alert(`An error occurred: ${errorMessage}`);
        if (reject) reject(error);
    }

    // Expose functions to global scope for HTML to call
    window.addFood = addFood;
    window.removeFood = removeFood;
});