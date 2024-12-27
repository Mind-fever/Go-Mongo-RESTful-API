package repositories

import (
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

var (
	recipeCollectionInstance   *mongo.Collection
	foodCollectionInstance     *mongo.Collection
	purchaseCollectionInstance *mongo.Collection
	onceRecipe                 sync.Once
	onceFood                   sync.Once
	oncePurchase               sync.Once
)

func getRecipeCollection(db DB) *mongo.Collection {
	onceRecipe.Do(func() {
		recipeCollectionInstance = db.GetClient().Database("SuperCook").Collection("recipes")
	})
	return recipeCollectionInstance
}

func getFoodCollection(db DB) *mongo.Collection {
	onceFood.Do(func() {
		foodCollectionInstance = db.GetClient().Database("SuperCook").Collection("foods")
	})
	return foodCollectionInstance
}

func getPurchaseCollection(db DB) *mongo.Collection {
	oncePurchase.Do(func() {
		purchaseCollectionInstance = db.GetClient().Database("SuperCook").Collection("purchases")
	})
	return purchaseCollectionInstance
}
