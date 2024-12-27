package dtos

import "github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/common"

type LowStockFilterDTO struct {
	UserID string `form:"user_id" validate:"required"`
	Name   string `form:"name" validate:"omitempty"`
	Type   string `form:"type" validate:"omitempty,validFoodType"`
}

func (dto *LowStockFilterDTO) Validate() error {
	return common.Validate.Struct(dto)
}
