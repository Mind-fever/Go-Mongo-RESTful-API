package services

import (
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/errors"
	"github.com/go-playground/validator/v10"
	"time"

	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/utils"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/dtos"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipeServiceInterface interface {
	GetRecipesByFilter(filter dtos.RecipeFilterDTO) ([]*dtos.RecipeUpdateDTO, error)
	GetRecipeByID(userID string, recipeID string) (*dtos.RecipeUpdateDTO, error)
	PostRecipe(userID string, recipe *dtos.RecipeDTO) (*dtos.RecipeUpdateDTO, error)
	DeleteRecipe(userID string, recipeID string) (*dtos.RecipeUpdateDTO, error)
	PutRecipe(userID string, recipe *dtos.RecipeUpdateDTO) (*dtos.RecipeUpdateDTO, error)
}

type RecipeService struct {
	recipeRepository repositories.RecipeRepositoryInterface
	foodRepository   repositories.FoodRepositoryInterface
}

func NewRecipeService(recipeRepository repositories.RecipeRepositoryInterface, foodRepository repositories.FoodRepositoryInterface) *RecipeService {
	return &RecipeService{
		recipeRepository: recipeRepository,
		foodRepository:   foodRepository,
	}
}

func (service *RecipeService) GetRecipesByFilter(filter dtos.RecipeFilterDTO) ([]*dtos.RecipeUpdateDTO, error) {
	err := filter.Validate()

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "UserID":
				if err.Tag() == "required" {
					return nil, errors.ErrUserNotFound
				}
			case "Use":
				if err.Tag() == "validMealTime" {
					return nil, errors.ErrInvalidMealTime
				}
			case "ProductType":
				if err.Tag() == "validProductType" {
					return nil, errors.ErrInvalidFoodType
				}
			}
		}
		return nil, errors.ErrValidationFailed
	}
	filters := bson.M{"user_id": filter.UserID}
	if filter.Use != "" {
		filters["meal_time"] = filter.Use
	}
	if filter.ProductType != "" || filter.ProductName != "" {
		foodFilters := bson.M{}
		if filter.ProductType != "" {
			foodFilters["type"] = filter.ProductType
		}
		if filter.ProductName != "" {
			foodFilters["name"] = bson.M{"$regex": filter.ProductName, "$options": "i"}
		}
		foods, err := service.foodRepository.GetFoodsByFilter(foodFilters)
		if err != nil {
			return nil, err
		}
		var foodIDs []primitive.ObjectID
		for _, food := range foods {
			foodIDs = append(foodIDs, food.ID)
		}
		filters["ingredients.food_id"] = bson.M{"$in": foodIDs}
	}
	recipesDB, err := service.recipeRepository.GetRecipesByFilter(filters)
	if err != nil {
		return nil, err
	}

	var recipes []*dtos.RecipeUpdateDTO
	for _, recipeDB := range recipesDB {
		recipe := dtos.NewRecipeUpdateDTO(recipeDB)
		recipes = append(recipes, recipe)
	}

	if len(recipes) == 0 {
		return nil, errors.ErrNoRecipesFound
	}

	return recipes, nil
}
func (service *RecipeService) GetRecipeByID(userID string, recipeID string) (*dtos.RecipeUpdateDTO, error) {
	if recipeID == "" {
		return nil, errors.ErrInvalidID
	}

	objectID := utils.GetObjectIDFromStringID(recipeID)

	filters := bson.M{"user_id": userID, "_id": objectID}

	recipe, err := service.recipeRepository.GetRecipeByID(filters)
	if err != nil {
		return nil, err
	}

	recipeDTO := dtos.NewRecipeUpdateDTO(recipe)
	return recipeDTO, nil
}

func (service *RecipeService) PostRecipe(userID string, recipe *dtos.RecipeDTO) (*dtos.RecipeUpdateDTO, error) {
	err := recipe.Validate()

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Name":
				if err.Tag() == "required" {
					return nil, errors.ErrInvalidName
				}
			case "MealTime":
				if err.Tag() == "required" {
					return nil, errors.ErrInvalidMealTime
				}
				if err.Tag() == "validMealTime" {
					return nil, errors.ErrInvalidMealTime
				}
			case "Ingredients":
				if err.Tag() == "required" {
					return nil, errors.ErrRecipeIngredientsMissing
				}
				if err.Tag() == "min" {
					return nil, errors.ErrRecipeIngredientsMissing
				}
			}
		}
		return nil, errors.ErrValidationFailed
	}

	for _, ingredient := range recipe.Ingredients {
		filter := bson.M{"_id": utils.GetObjectIDFromStringID(ingredient.FoodID), "user_id": userID}
		food, err := service.foodRepository.GetFoodByID(filter)
		if err != nil {
			return nil, errors.ErrFoodNotFound
		}

		if food.CurrentQuantity < ingredient.Quantity {
			return nil, errors.ErrRecipeInsufficientIngredients
		}

		validMealTime := false
		for _, mealTime := range food.MealTimes {
			if mealTime == recipe.MealTime {
				validMealTime = true
				break
			}
		}
		if !validMealTime {
			return nil, errors.ErrMismatchMealTime
		}

		food.CurrentQuantity -= ingredient.Quantity
		update := bson.M{"$set": bson.M{"current_quantity": food.CurrentQuantity}}
		_, err = service.foodRepository.UpdateFoodQuantity(filter, update)
		if err != nil {
			return nil, errors.ErrRecipeUpdatingIngredients
		}
	}

	recipeModel := recipe.GetModel()
	recipeModel.UserID = userID
	recipeModel.CreatedAt = time.Now()
	recipeModel.UpdatedAt = time.Now()

	result, err := service.recipeRepository.PostRecipe(recipeModel)
	if err != nil {
		return nil, err
	}

	insertedID := result.InsertedID.(primitive.ObjectID)
	responseDTO := &dtos.RecipeUpdateDTO{
		ID:          insertedID.Hex(),
		Name:        recipe.Name,
		MealTime:    recipe.MealTime,
		Ingredients: recipe.Ingredients,
	}
	return responseDTO, nil
}

func (service *RecipeService) DeleteRecipe(userID string, recipeID string) (*dtos.RecipeUpdateDTO, error) {
	if recipeID == "" {
		return nil, errors.ErrInvalidID
	}
	objectID := utils.GetObjectIDFromStringID(recipeID)

	filter := bson.M{"_id": objectID, "user_id": userID}
	recipe, err := service.recipeRepository.GetRecipeByID(filter)
	if err != nil {
		return nil, errors.ErrRecipeNotFound
	}

	for _, ingredient := range recipe.Ingredients {
		filter := bson.M{"_id": ingredient.FoodID, "user_id": userID}

		food, err := service.foodRepository.GetFoodByID(filter)
		if err != nil {
			return nil, errors.ErrRecipeIngredientNotFound
		}

		food.CurrentQuantity += ingredient.Quantity
		update := bson.M{"$set": bson.M{"current_quantity": food.CurrentQuantity}}
		_, err = service.foodRepository.UpdateFoodQuantity(filter, update)
		if err != nil {
			return nil, errors.ErrRecipeUpdatingIngredients
		}
	}

	result, err := service.recipeRepository.DeleteRecipe(filter)
	if err != nil {
		return nil, err
	}

	if result.DeletedCount == 0 {
		return nil, errors.ErrRecipeNotFound
	}

	deletedRecipeDTO := dtos.NewRecipeUpdateDTO(recipe)
	return deletedRecipeDTO, nil
}

func (service *RecipeService) PutRecipe(userID string, recipe *dtos.RecipeUpdateDTO) (*dtos.RecipeUpdateDTO, error) {
	err := recipe.Validate()

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "ID":
				if err.Tag() == "required" {
					return nil, errors.ErrInvalidID
				}
			case "Name":
				if err.Tag() == "required" {
					return nil, errors.ErrInvalidName
				}
			case "MealTime":
				if err.Tag() == "required" {
					return nil, errors.ErrInvalidMealTime
				}
				if err.Tag() == "validMealTime" {
					return nil, errors.ErrInvalidMealTime
				}
			case "Ingredients":
				if err.Tag() == "min" {
					return nil, errors.ErrRecipeIngredientsMissing
				}
			}
		}
		return nil, errors.ErrValidationFailed
	}

	objectID := utils.GetObjectIDFromStringID(recipe.ID)
	filter := bson.M{"_id": objectID, "user_id": userID}

	// Obtener la receta existente
	existingRecipe, err := service.recipeRepository.GetRecipeByID(filter)
	if err != nil {
		return nil, errors.ErrRecipeNotFound
	}

	// Primera iteracion: Verificacion de ingredientes de la receta actualizada
	for _, newIngredient := range recipe.Ingredients {
		filter := bson.M{"_id": utils.GetObjectIDFromStringID(newIngredient.FoodID), "user_id": userID}
		food, err := service.foodRepository.GetFoodByID(filter)
		if err != nil {
			return nil, errors.ErrRecipeIngredientNotFound
		}

		// Obtener el stock por cada uno en la receta existente
		existingQty := float64(0)
		for _, existingIngredient := range existingRecipe.Ingredients {
			if existingIngredient.FoodID == utils.GetObjectIDFromStringID(newIngredient.FoodID) {
				existingQty = existingIngredient.Quantity
				break
			}
		}

		// Calcular la diferencia de cantidad
		quantityDiff := newIngredient.Quantity - existingQty

		// Si se necesita mAs cantidad de la existente, verifica que estE disponible
		if quantityDiff > 0 && food.CurrentQuantity < quantityDiff {
			return nil, errors.ErrRecipeInsufficientIngredients
		}

		// Verificar MealTimes
		validMealTime := false
		for _, mealTime := range food.MealTimes {
			if mealTime == recipe.MealTime {
				validMealTime = true
				break
			}
		}
		if !validMealTime {
			return nil, errors.ErrMismatchMealTime
		}
	}

	// Segunda Iteracion: Aplicación de los cambios
	// Devolver ingredientes de la receta existente al inventario al restarlos
	for _, ingredient := range existingRecipe.Ingredients {
		filter := bson.M{"_id": ingredient.FoodID, "user_id": userID}

		food, err := service.foodRepository.GetFoodByID(filter)
		if err != nil {
			return nil, errors.ErrRecipeIngredientNotFound
		}

		food.CurrentQuantity += ingredient.Quantity
		update := bson.M{"$set": bson.M{"current_quantity": food.CurrentQuantity}}
		_, err = service.foodRepository.UpdateFoodQuantity(filter, update)
		if err != nil {
			return nil, errors.ErrRecipeUpdatingIngredients
		}
	}

	// Reducir cantidades segun los ingredientes de la nueva receta al sumarlos
	for _, newIngredient := range recipe.Ingredients {
		filter := bson.M{"_id": utils.GetObjectIDFromStringID(newIngredient.FoodID), "user_id": userID}

		food, err := service.foodRepository.GetFoodByID(filter)
		if err != nil {
			return nil, errors.ErrRecipeIngredientNotFound
		}

		food.CurrentQuantity -= newIngredient.Quantity
		update := bson.M{"$set": bson.M{"current_quantity": food.CurrentQuantity}}
		_, err = service.foodRepository.UpdateFoodQuantity(filter, update)
		if err != nil {
			return nil, errors.ErrRecipeUpdatingIngredients
		}
	}

	//Actualizar la receta en la base de datos
	recipeModel := recipe.GetModel()
	recipeModel.UpdatedAt = time.Now()
	recipeModel.CreatedAt = existingRecipe.CreatedAt

	update := bson.M{
		"$set": recipeModel,
	}

	result, err := service.recipeRepository.PutRecipe(filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.ErrRecipeNotFound
	}

	return recipe, nil
}
