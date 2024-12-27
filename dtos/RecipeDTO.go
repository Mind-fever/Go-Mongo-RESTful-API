package dtos

import (
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/common"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipeDTO struct {
	Name        string                `json:"name" validate:"required"`
	MealTime    string                `json:"meal_time" validate:"required,validMealTime"`
	Ingredients []RecipeIngredientDTO `json:"ingredients" validate:"required,min=1,dive"`
}

type RecipeIngredientDTO struct {
	FoodID   string  `json:"food_id"`
	Quantity float64 `json:"quantity"`
}

func (dto *RecipeDTO) Validate() error {
	return common.Validate.Struct(dto)
}

func (dto *RecipeDTO) GetModel() *models.Recipe {
	var ingredients []models.RecipeIngredient
	for _, ingredient := range dto.Ingredients {
		foodID, _ := primitive.ObjectIDFromHex(ingredient.FoodID)
		ingredients = append(ingredients, models.RecipeIngredient{
			FoodID:   foodID,
			Quantity: ingredient.Quantity,
		})
	}

	return &models.Recipe{
		Name:        dto.Name,
		MealTime:    dto.MealTime,
		Ingredients: ingredients,
	}
}
