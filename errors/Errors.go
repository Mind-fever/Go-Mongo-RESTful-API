package errors

var (
	// Errores Generales
	ErrValidationFailed = New("ERR_1", "Validation failed")
	ErrDatabasePut      = New("ERR_2", "Failed to update the database")
	ErrDatabaseDelete   = New("ERR_3", "Failed to delete from the database")
	ErrDatabasePost     = New("ERR_4", "Failed to insert into the database")
	ErrDatabaseGet      = New("ERR_5", "Failed to retrieve from the database")

	// Errores de Alimentos
	ErrFoodNotFound     = New("ERR_6", "Food item not found")
	ErrFoodInUse        = New("ERR_7", "Food item is in use")
	ErrMismatchMealTime = New("ERR_8", "Meal time does not match recipe usage")
	ErrNoFoodsFound     = New("ERR_28", "No foods found")

	// Errores de Validación de Alimentos
	ErrInvalidFoodType        = New("ERR_9", "Invalid food type")
	ErrInvalidMealTime        = New("ERR_10", "Invalid meal time")
	ErrCurrentQuantityInvalid = New("ERR_11", "Current quantity cannot be less than zero")
	ErrMinQuantityInvalid     = New("ERR_12", "Minimum quantity cannot be less than one")
	ErrPricePerUnitInvalid    = New("ERR_13", "Price per unit cannot be less than one")

	// Errores de Validación de Campos
	ErrInvalidName = New("ERR_14", "Invalid name provided")
	ErrInvalidID   = New("ERR_15", "The provided item could not be identified")

	// Errores de Recetas
	ErrRecipeInsufficientIngredients = New("ERR_16", "Not enough ingredients for this recipe")
	ErrRecipeNotFound                = New("ERR_17", "Recipe not found")
	ErrRecipeUpdatingIngredients     = New("ERR_18", "Error updating recipe ingredients")
	ErrRecipeIngredientNotFound      = New("ERR_19", "Recipe ingredient not found")
	ErrRecipeIngredientsMissing      = New("ERR_20", "No ingredients provided for the recipe")
	ErrNoRecipesFound                = New("ERR_27", "No recipes found")

	// Errores de Compras
	ErrPurchaseFoodNotFound = New("ERR_21", "Food item not found")

	// Errores de Usuarios
	ErrUserNotFound = New("ERR_22", "User not found")
	ErrUnauthorized = New("ERR_23", "Unauthorized access")

	// Errores Internos
	ErrCursorDecodeFailed    = New("ERR_24", "Failed to decode cursor result")
	ErrCursorIterationFailed = New("ERR_25", "Failed to iterate cursor result")
	ErrBindingFailed         = New("ERR_26", "Failed to bind data")
)
