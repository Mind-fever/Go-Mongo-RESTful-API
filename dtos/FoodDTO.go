package dtos

import (
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/common"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/models"
)

type FoodDTO struct {
	Name            string   `json:"name" validate:"required"`
	Type            string   `json:"type" validate:"required,validFoodType"`
	PricePerUnit    float64  `json:"price_per_unit" validate:"required,gt=0"`
	CurrentQuantity float64  `json:"current_quantity" validate:"required,gte=0"`
	MinQuantity     float64  `json:"min_quantity" validate:"required,gt=0"`
	MealTimes       []string `json:"meal_times" validate:"required,min=1,dive,validMealTime"`
}

func (dto *FoodDTO) Validate() error {
	return common.Validate.Struct(dto)
}

func (dto *FoodDTO) GetModel() *models.Food {
	return &models.Food{
		Name:            dto.Name,
		Type:            dto.Type,
		PricePerUnit:    dto.PricePerUnit,
		CurrentQuantity: dto.CurrentQuantity,
		MinQuantity:     dto.MinQuantity,
		MealTimes:       dto.MealTimes,
	}
}
