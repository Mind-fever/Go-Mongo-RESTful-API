package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Food struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" `
	UserID          string             `bson:"user_id,omitempty"`
	Name            string             `bson:"name"`
	Type            string             `bson:"type"`
	PricePerUnit    float64            `bson:"price_per_unit"`
	CurrentQuantity float64            `bson:"current_quantity"`
	MinQuantity     float64            `bson:"min_quantity"`
	MealTimes       []string           `bson:"meal_times"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
}
