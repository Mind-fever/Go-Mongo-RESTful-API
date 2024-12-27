document.addEventListener('DOMContentLoaded', function() {
    checkTokenAndInitialize();
});

function checkTokenAndInitialize() {
    const token = localStorage.getItem('authToken');
    const loginStatus = document.getElementById('loginStatus');
    const logoutButton = document.getElementById('logoutButton');

    if (!token) {
        alert('No estás logueado. Serás redirigido a la página de login.');
        window.location.href = 'login.html';
    } else {
        loginStatus.textContent = 'Logged in';
        logoutButton.style.display = 'inline-block';
        initializeForm();
    }

    logoutButton.addEventListener('click', function() {
        localStorage.removeItem('authToken');
        window.location.href = 'login.html';
        loginStatus.textContent = 'Not logged in';
        logoutButton.style.display = 'none';
    });
}

function initializeForm() {
    const urlParams = new URLSearchParams(window.location.search);
    const foodId = urlParams.get('id');
    const foodName = urlParams.get('name');
    const foodType = urlParams.get('type');
    const pricePerUnit = urlParams.get('price_per_unit');
    const currentQuantity = urlParams.get('current_quantity');
    const minQuantity = urlParams.get('min_quantity');
    const mealTimes = urlParams.get('meal_times').split(',');

    populateFoodForm(foodName, foodType, pricePerUnit, currentQuantity, minQuantity, mealTimes);
    handleFormSubmit(foodId);
}

function populateFoodForm(foodName, foodType, pricePerUnit, currentQuantity, minQuantity, mealTimes) {
    document.getElementById('foodName').value = foodName;
    document.getElementById('foodType').value = foodType;
    document.getElementById('pricePerUnit').value = pricePerUnit;
    document.getElementById('currentQuantity').value = currentQuantity;
    document.getElementById('minQuantity').value = minQuantity;
    mealTimes.forEach(mealTime => {
        document.getElementById(`mealTime${mealTime.charAt(0).toUpperCase() + mealTime.slice(1)}`).checked = true;
    });
}

function handleFormSubmit(foodId) {
    const updateFoodForm = document.getElementById('updateFoodForm');
    updateFoodForm.addEventListener('submit', function(event) {
        event.preventDefault();

        const formData = new FormData(updateFoodForm);
        const foodData = {
            id: foodId,
            name: formData.get('foodName'),
            type: formData.get('foodType'),
            price_per_unit: parseFloat(formData.get('pricePerUnit')),
            current_quantity: parseFloat(formData.get('currentQuantity')),
            min_quantity: parseFloat(formData.get('minQuantity')),
            meal_times: Array.from(formData.getAll('mealTimes'))
        };

        makeRequest(`http://localhost:8080/foods/${foodId}`, 'PUT', foodData, 'application/json', true,
            (data) => {
                console.log('Success:', data);
                alert('Food updated successfully!');
                window.location.href = 'foods.html';
            },
            (status, error) => {
                console.error('Error:', error);
                const errorMessage = error.error || 'An unknown error occurred';
                alert(`An error occurred: ${errorMessage}`);
            }
        );
    });
}