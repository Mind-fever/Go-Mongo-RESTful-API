package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Recipe struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"user_id,omitempty" `
	Name        string             `bson:"name"`
	MealTime    string             `bson:"meal_time"`
	Ingredients []RecipeIngredient `bson:"ingredients" `
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" `
}

type RecipeIngredient struct {
	FoodID   primitive.ObjectID `bson:"food_id"`
	Quantity float64            `bson:"quantity"`
}
