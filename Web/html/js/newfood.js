document.addEventListener('DOMContentLoaded', function() {
    checkTokenAndInitialize(initializeNewFoodForm);
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

function initializeNewFoodForm() {
    const newFoodForm = document.getElementById('newFoodForm');
    newFoodForm.addEventListener('submit', handleFormSubmit);
}

function handleFormSubmit(event) {
    event.preventDefault();

    const newFoodForm = event.target;
    const formData = new FormData(newFoodForm);
    const foodData = {
        name: formData.get('foodName'),
        type: formData.get('foodType'),
        price_per_unit: parseFloat(formData.get('pricePerUnit')),
        current_quantity: parseFloat(formData.get('currentQuantity')),
        min_quantity: parseFloat(formData.get('minQuantity')),
        meal_times: Array.from(formData.getAll('mealTimes'))
    };

    makeRequest('http://localhost:8080/foods/', 'POST', foodData, 'application/json', true,
        (data) => {
            console.log('Success:', data);
            alert('Food added successfully!');
            window.location.href = 'foods.html';
        },
        (status, error) => handleError(error)
    );
    function handleError(error, reject) {
        console.error('Error:', error);
        const errorMessage = error.error || 'An unknown error occurred';
        alert(`An error occurred: ${errorMessage}`);
        if (reject) reject(error);
    }
}