package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Purchase struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id,omitempty"`
	Date      time.Time          `bson:"date"`
	TotalCost float64            `bson:"total_cost"`
	Items     []PurchaseItem     `bson:"items"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type PurchaseItem struct {
	FoodID   primitive.ObjectID `bson:"food_id"`
	Quantity float64            `bson:"quantity"`
	UnitCost float64            `bson:"unit_cost"`
}
