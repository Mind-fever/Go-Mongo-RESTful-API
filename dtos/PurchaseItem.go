package dtos

import (
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/common"
)

type PurchaseItem struct {
	FoodID   string  `json:"food_id" validate:"required"`
	Quantity float64 `json:"quantity" validate:"required,gte=0"`
	UnitCost float64 `json:"unit_cost" validate:"required,gte=0"`
}

func (item *PurchaseItem) Validate() error {
	return common.Validate.Struct(item)
}
