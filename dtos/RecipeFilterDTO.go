package dtos

import "github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/common"

type RecipeFilterDTO struct {
	UserID      string `json:"user_id" form:"user_id" validate:"required"`
	Use         string `json:"meal_time" form:"meal_time" validate:"omitempty,validMealTime"`
	ProductType string `json:"product_type" form:"product_type" validate:"omitempty,validFoodType"`
	ProductName string `json:"product_name" form:"product_name" validate:"omitempty"`
}

func (dto *RecipeFilterDTO) Validate() error {
	return common.Validate.Struct(dto)
}
