package services

import (
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/errors"
	"github.com/go-playground/validator/v10"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/dtos"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/repositories"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FoodServiceInterface interface {
	GetFoods(userID string) ([]*dtos.FoodUpdateDTO, error)
	PostFood(food *dtos.FoodDTO, userID string) (*dtos.FoodUpdateDTO, error)
	GetLowStockFoods(filter *dtos.LowStockFilterDTO) ([]*dtos.FoodUpdateDTO, error)
	PutFood(food *dtos.FoodUpdateDTO, userID string) (*dtos.FoodUpdateDTO, error)
	DeleteFood(userID string, foodID string) (*dtos.FoodUpdateDTO, error)
	GetFoodByID(userID string, foodID string) (*dtos.FoodUpdateDTO, error)
}

type FoodService struct {
	foodRepository   repositories.FoodRepositoryInterface
	recipeRepository repositories.RecipeRepositoryInterface
}

func NewFoodService(foodRepository repositories.FoodRepositoryInterface, recipeRepository repositories.RecipeRepositoryInterface) *FoodService {
	return &FoodService{
		foodRepository:   foodRepository,
		recipeRepository: recipeRepository,
	}
}

func (service *FoodService) GetFoods(userID string) ([]*dtos.FoodUpdateDTO, error) {
	filters := bson.M{"user_id": userID}

	foodsDB, err := service.foodRepository.GetFoods(filters)
	if err != nil {
		return nil, err
	}

	var foods []*dtos.FoodUpdateDTO
	for _, foodDB := range foodsDB {
		food := dtos.NewFoodUpdateDTO(foodDB)
		foods = append(foods, food)
	}

	return foods, nil
}

func (service *FoodService) GetLowStockFoods(filter *dtos.LowStockFilterDTO) ([]*dtos.FoodUpdateDTO, error) {
	err := filter.Validate()

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "UserID":
				if err.Tag() == "required" {
					return nil, errors.ErrUserNotFound
				}
			case "Type":
				if err.Tag() == "validFoodTime" {
					return nil, errors.ErrInvalidMealTime
				}
			}
		}
		return nil, errors.ErrValidationFailed
	}

	filters := bson.M{
		"$expr": bson.M{
			"$lt": []interface{}{"$current_quantity", "$min_quantity"},
		},
	}
	if filter.Type != "" {
		filters["type"] = filter.Type
	}
	if filter.Name != "" {
		filters["name"] = bson.M{"$regex": filter.Name, "$options": "i"}
	}

	foodsDB, err := service.foodRepository.GetLowStockFoods(filters)
	if err != nil {
		return nil, err
	}

	var foods []*dtos.FoodUpdateDTO
	for _, foodDB := range foodsDB {
		food := dtos.NewFoodUpdateDTO(foodDB)
		foods = append(foods, food)
	}

	if len(foods) == 0 {
		return nil, errors.ErrNoFoodsFound

	}

	return foods, nil
}

func (service *FoodService) GetFoodByID(userID string, foodID string) (*dtos.FoodUpdateDTO, error) {
	if foodID == "" {
		return nil, errors.ErrInvalidID
	}

	objectID := utils.GetObjectIDFromStringID(foodID)

	filter := bson.M{"_id": objectID, "user_id": userID}

	food, err := service.foodRepository.GetFoodByID(filter)
	if err != nil {
		return nil, err
	}

	if food == nil {
		return nil, errors.ErrFoodNotFound
	}

	foodDTO := dtos.NewFoodUpdateDTO(food)
	return foodDTO, nil
}

func (service *FoodService) PostFood(food *dtos.FoodDTO, userID string) (*dtos.FoodUpdateDTO, error) {
	err := food.Validate()

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Type":
				if err.Tag() == "validFoodType" {
					return nil, errors.ErrInvalidFoodType
				}
			case "MealTimes":
				if err.Tag() == "validMealTime" {
					return nil, errors.ErrInvalidMealTime
				}
			case "CurrentQuantity":
				if err.Tag() == "gte" {
					return nil, errors.ErrCurrentQuantityInvalid
				}
			case "MinQuantity":
				if err.Tag() == "gt" {
					return nil, errors.ErrMinQuantityInvalid
				}
			case "PricePerUnit":
				if err.Tag() == "gt" {
					return nil, errors.ErrPricePerUnitInvalid
				}
			}
		}
		return nil, errors.ErrValidationFailed
	}

	foodModel := food.GetModel()
	foodModel.CreatedAt = time.Now()
	foodModel.UpdatedAt = time.Now()
	foodModel.UserID = userID

	result, err := service.foodRepository.PostFood(foodModel)
	if err != nil {
		return nil, err
	}

	insertedID := result.InsertedID.(primitive.ObjectID).Hex()

	responseDTO := &dtos.FoodUpdateDTO{
		ID:              insertedID,
		Name:            food.Name,
		Type:            food.Type,
		PricePerUnit:    food.PricePerUnit,
		CurrentQuantity: food.CurrentQuantity,
		MinQuantity:     food.MinQuantity,
		MealTimes:       food.MealTimes,
	}

	return responseDTO, nil
}

func (service *FoodService) PutFood(food *dtos.FoodUpdateDTO, userID string) (*dtos.FoodUpdateDTO, error) {
	err := food.Validate()

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "ID":
				if err.Tag() == "required" {
					return nil, errors.ErrInvalidID
				}
			case "Type":
				if err.Tag() == "validFoodType" {
					return nil, errors.ErrInvalidFoodType
				}
			case "MealTimes":
				if err.Tag() == "validMealTime" {
					return nil, errors.ErrInvalidMealTime
				}
			case "CurrentQuantity":
				if err.Tag() == "gte" {
					return nil, errors.ErrCurrentQuantityInvalid
				}
			case "MinQuantity":
				if err.Tag() == "gt" {
					return nil, errors.ErrMinQuantityInvalid
				}
			case "PricePerUnit":
				if err.Tag() == "gt" {
					return nil, errors.ErrPricePerUnitInvalid
				}
			}
		}
		return nil, errors.ErrValidationFailed
	}

	foodModel := food.GetModel()
	foodModel.UpdatedAt = time.Now()

	filter := bson.M{"user_id": userID, "_id": foodModel.ID}
	existingFood, err := service.foodRepository.GetFoodByID(filter)
	if err != nil {
		return nil, errors.ErrFoodNotFound
	}
	foodModel.CreatedAt = existingFood.CreatedAt

	update := bson.M{
		"$set": foodModel,
	}
	result, err := service.foodRepository.PutFood(filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.ErrFoodNotFound
	}

	return food, nil
}

func (service *FoodService) DeleteFood(userID string, foodID string) (*dtos.FoodUpdateDTO, error) {
	if foodID == "" {
		return nil, errors.ErrInvalidID
	}

	objectID := utils.GetObjectIDFromStringID(foodID)
	filter := bson.M{"user_id": userID, "_id": objectID}

	food, err := service.foodRepository.GetFoodByID(filter)
	if err != nil {
		return nil, errors.ErrFoodNotFound
	}

	recipeFilter := bson.M{"ingredients.food_id": objectID, "user_id": userID}
	recipes, err := service.recipeRepository.GetRecipesByFilter(recipeFilter)
	if err != nil {
		return nil, errors.ErrDatabaseGet
	}
	if len(recipes) > 0 {
		return nil, errors.ErrFoodInUse
	}

	result, err := service.foodRepository.DeleteFood(filter)
	if err != nil {
		return nil, err
	}

	if result.DeletedCount == 0 {
		return nil, errors.ErrFoodNotFound
	}

	deletedFoodDTO := dtos.NewFoodUpdateDTO(food)
	return deletedFoodDTO, nil
}
