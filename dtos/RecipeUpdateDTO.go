package dtos

import (
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/models"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/common"
)

type RecipeUpdateDTO struct {
	ID          string                `json:"id" validate:"required"`
	Name        string                `json:"name" validate:"required"`
	MealTime    string                `json:"meal_time" validate:"required,validMealTime"`
	Ingredients []RecipeIngredientDTO `json:"ingredients" validate:"required,min=1,dive"`
}

func (dto *RecipeUpdateDTO) Validate() error {
	return common.Validate.Struct(dto)
}

func NewRecipeUpdateDTO(recipe *models.Recipe) *RecipeUpdateDTO {
	var ingredients []RecipeIngredientDTO
	for _, ingredient := range recipe.Ingredients {
		ingredients = append(ingredients, RecipeIngredientDTO{
			FoodID:   ingredient.FoodID.Hex(),
			Quantity: ingredient.Quantity,
		})
	}

	return &RecipeUpdateDTO{
		ID:          recipe.ID.Hex(),
		Name:        recipe.Name,
		MealTime:    recipe.MealTime,
		Ingredients: ingredients,
	}
}

func (dto *RecipeUpdateDTO) GetModel() *models.Recipe {
	var ingredients []models.RecipeIngredient
	for _, ingredient := range dto.Ingredients {
		foodID, _ := primitive.ObjectIDFromHex(ingredient.FoodID)
		ingredients = append(ingredients, models.RecipeIngredient{
			FoodID:   foodID,
			Quantity: ingredient.Quantity,
		})
	}

	return &models.Recipe{
		ID:          utils.GetObjectIDFromStringID(dto.ID),
		Name:        dto.Name,
		MealTime:    dto.MealTime,
		Ingredients: ingredients,
	}
}
