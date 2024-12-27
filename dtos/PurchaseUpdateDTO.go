package dtos

import (
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/common"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PurchaseUpdateDTO struct {
	ID        string         `json:"id" validate:"required"`
	Date      time.Time      `json:"date" validate:"required"`
	TotalCost float64        `json:"total_cost" validate:"required,gte=0"`
	Items     []PurchaseItem `json:"items" validate:"required,min=1,dive"`
}

func (dto *PurchaseUpdateDTO) Validate() error {
	if err := common.Validate.Struct(dto); err != nil {
		return err
	}
	for _, item := range dto.Items {
		if err := item.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func NewPurchaseUpdateDTO(purchase *models.Purchase) *PurchaseUpdateDTO {
	var items []PurchaseItem
	for _, item := range purchase.Items {
		items = append(items, PurchaseItem{
			FoodID:   item.FoodID.Hex(),
			Quantity: item.Quantity,
			UnitCost: item.UnitCost,
		})
	}

	return &PurchaseUpdateDTO{
		Date:      purchase.Date,
		TotalCost: purchase.TotalCost,
		Items:     items,
	}
}

func (dto *PurchaseUpdateDTO) GetModel() *models.Purchase {
	var items []models.PurchaseItem
	for _, item := range dto.Items {
		foodID, _ := primitive.ObjectIDFromHex(item.FoodID)
		items = append(items, models.PurchaseItem{
			FoodID:   foodID,
			Quantity: item.Quantity,
			UnitCost: item.UnitCost,
		})
	}

	return &models.Purchase{
		Date:      dto.Date,
		TotalCost: dto.TotalCost,
		Items:     items,
	}
}
