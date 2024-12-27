package repositories

import (
	"context"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/errors"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FoodRepositoryInterface interface {
	GetFoods(filters bson.M) ([]*models.Food, error)
	GetLowStockFoods(filters bson.M) ([]*models.Food, error)
	PostFood(food *models.Food) (*mongo.InsertOneResult, error)
	PutFood(filter bson.M, update bson.M) (*mongo.UpdateResult, error)
	UpdateFoodQuantity(filter bson.M, update bson.M) (*mongo.UpdateResult, error)
	DeleteFood(filter bson.M) (*mongo.DeleteResult, error)
	GetFoodByID(filter bson.M) (*models.Food, error)
	GetFoodsByFilter(filters bson.M) ([]*models.Food, error)
}

type FoodRepository struct {
	db DB
}

func (repository *FoodRepository) getFoodCollection() *mongo.Collection {
	return getFoodCollection(repository.db)
}

func NewFoodRepository(db DB) *FoodRepository {
	return &FoodRepository{
		db: db,
	}
}

func (repository FoodRepository) GetFoods(filters bson.M) ([]*models.Food, error) {
	collection := repository.getFoodCollection()

	cursor, err := collection.Find(context.TODO(), filters)
	if err != nil {
		return nil, errors.ErrDatabaseGet
	}

	defer cursor.Close(context.Background())

	var foods []*models.Food
	for cursor.Next(context.Background()) {
		var food models.Food
		if err := cursor.Decode(&food); err != nil {
			return nil, errors.ErrCursorDecodeFailed
		}
		foods = append(foods, &food)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.ErrCursorIterationFailed
	}

	return foods, nil
}

func (repository FoodRepository) GetLowStockFoods(filters bson.M) ([]*models.Food, error) {
	collection := repository.getFoodCollection()

	cursor, err := collection.Find(context.Background(), filters)
	if err != nil {
		return nil, errors.ErrDatabaseGet
	}

	defer cursor.Close(context.Background())

	var foods []*models.Food
	for cursor.Next(context.Background()) {
		var food models.Food
		if err := cursor.Decode(&food); err != nil {
			return nil, errors.ErrCursorDecodeFailed
		}
		foods = append(foods, &food)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.ErrCursorIterationFailed
	}

	return foods, nil
}

func (repository FoodRepository) PostFood(food *models.Food) (*mongo.InsertOneResult, error) {
	collection := repository.getFoodCollection()

	result, err := collection.InsertOne(context.TODO(), food)
	if err != nil {
		return nil, errors.ErrDatabasePost
	}
	return result, nil
}

func (repository FoodRepository) PutFood(filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	collection := repository.getFoodCollection()

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, errors.ErrDatabasePut
	}

	return result, nil
}

func (repository FoodRepository) UpdateFoodQuantity(filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	collection := repository.getFoodCollection()

	result, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return nil, errors.ErrDatabasePut
	}

	return result, nil
}

func (repository FoodRepository) DeleteFood(filter bson.M) (*mongo.DeleteResult, error) {
	collection := repository.getFoodCollection()

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, errors.ErrDatabaseDelete
	}

	return result, nil
}

func (repository FoodRepository) GetFoodByID(filter bson.M) (*models.Food, error) {
	collection := repository.db.GetClient().Database("SuperCook").Collection("foods")

	var food models.Food
	err := collection.FindOne(context.TODO(), filter).Decode(&food)
	if err != nil {
		return nil, errors.ErrDatabaseGet
	}

	return &food, nil
}

func (repository FoodRepository) GetFoodsByFilter(filters bson.M) ([]*models.Food, error) {
	collection := repository.getFoodCollection()

	cursor, err := collection.Find(context.Background(), filters)
	if err != nil {
		return nil, errors.ErrDatabaseGet
	}

	defer cursor.Close(context.Background())

	var foods []*models.Food
	for cursor.Next(context.Background()) {
		var food models.Food
		if err := cursor.Decode(&food); err != nil {
			return nil, errors.ErrCursorDecodeFailed
		}
		foods = append(foods, &food)
	}
	return foods, nil
}
