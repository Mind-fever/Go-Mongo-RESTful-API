package common

// Food types
const (
	FoodTypeVegetables = "vegetables"
	FoodTypeFruits     = "fruits"
	FoodTypeCheeses    = "cheeses"
	FoodTypeDairy      = "dairy"
	FoodTypeMeats      = "meats"
)

// Valid food types
var ValidFoodTypes = []string{
	FoodTypeVegetables,
	FoodTypeFruits,
	FoodTypeCheeses,
	FoodTypeDairy,
	FoodTypeMeats,
}

// Meal Times
const (
	MealTimeBreakfast = "breakfast"
	MealTimeLunch     = "lunch"
	MealTimeSnack     = "snack"
	MealTimeDinner    = "dinner"
)

// Valid Meal Times
var ValidMealTimes = []string{
	MealTimeBreakfast,
	MealTimeLunch,
	MealTimeSnack,
	MealTimeDinner,
}

// Food Filter
const (
	FilterByName = "name"
	FilterByType = "type"
)

// Valid Food Filters
var FoodFilter = []string{
	FilterByName,
	FilterByType,
}
