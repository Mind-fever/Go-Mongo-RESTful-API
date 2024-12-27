package dtos

import (
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/common"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/models"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/utils"
)

type FoodUpdateDTO struct {
	ID              string   `json:"id" validate:"required"`
	Name            string   `json:"name" validate:"required"`
	Type            string   `json:"type" validate:"required,validFoodType"`
	PricePerUnit    float64  `json:"price_per_unit" validate:"required,gt=0"`
	CurrentQuantity float64  `json:"current_quantity" validate:"required,gte=0"`
	MinQuantity     float64  `json:"min_quantity" validate:"required,gt=0"`
	MealTimes       []string `json:"meal_times" validate:"required,min=1,dive,validMealTime"`
}

func (dto *FoodUpdateDTO) Validate() error {
	return common.Validate.Struct(dto)
}

func (dto *FoodUpdateDTO) GetModel() *models.Food {
	return &models.Food{
		ID:              utils.GetObjectIDFromStringID(dto.ID),
		Name:            dto.Name,
		Type:            dto.Type,
		PricePerUnit:    dto.PricePerUnit,
		CurrentQuantity: dto.CurrentQuantity,
		MinQuantity:     dto.MinQuantity,
		MealTimes:       dto.MealTimes,
	}
}

func NewFoodUpdateDTO(food *models.Food) *FoodUpdateDTO {
	return &FoodUpdateDTO{
		ID:              food.ID.Hex(),
		Name:            food.Name,
		Type:            food.Type,
		PricePerUnit:    food.PricePerUnit,
		CurrentQuantity: food.CurrentQuantity,
		MinQuantity:     food.MinQuantity,
		MealTimes:       food.MealTimes,
	}
}
