package services

import (
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/dtos"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/errors"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/models"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/repositories"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PurchaseServiceInterface interface {
	GetPurchases(userID string) ([]*dtos.PurchaseUpdateDTO, error)
	PostPurchase(userID string) (*dtos.PurchaseUpdateDTO, error)
}

type PurchaseService struct {
	foodRepository     repositories.FoodRepositoryInterface
	purchaseRepository repositories.PurchaseRepositoryInterface
	foodService        *FoodService
}

func NewPurchaseService(purchaseRepository repositories.PurchaseRepositoryInterface, foodRepository repositories.FoodRepositoryInterface, foodService *FoodService) *PurchaseService {
	return &PurchaseService{
		purchaseRepository: purchaseRepository,
		foodRepository:     foodRepository,
		foodService:        foodService,
	}
}

func (service *PurchaseService) GetPurchases(userID string) ([]*dtos.PurchaseUpdateDTO, error) {
	filter := bson.M{"user_id": userID}

	purchasesDB, err := service.purchaseRepository.GetPurchases(filter)
	if err != nil {
		return nil, err
	}

	var purchases []*dtos.PurchaseUpdateDTO
	for _, purchaseDB := range purchasesDB {
		purchase := dtos.NewPurchaseUpdateDTO(purchaseDB)
		purchases = append(purchases, purchase)
	}

	return purchases, nil
}
func (service *PurchaseService) PostPurchase(userID string) (*dtos.PurchaseUpdateDTO, error) {
	filter := dtos.LowStockFilterDTO{UserID: userID}

	foods, err := service.foodService.GetLowStockFoods(&filter)
	if err != nil {
		return nil, errors.ErrPurchaseFoodNotFound
	}

	var totalCost float64
	var purchaseItems []models.PurchaseItem
	for _, food := range foods {
		quantityToPurchase := food.MinQuantity - food.CurrentQuantity
		totalCost += food.PricePerUnit * quantityToPurchase

		purchaseItem := models.PurchaseItem{
			FoodID:   utils.GetObjectIDFromStringID(food.ID),
			Quantity: quantityToPurchase,
			UnitCost: food.PricePerUnit,
		}
		purchaseItems = append(purchaseItems, purchaseItem)

		food.CurrentQuantity = food.MinQuantity

		foodModel := food.GetModel()
		filterPut := bson.M{"user_id": userID, "_id": foodModel.ID}
		update := bson.M{
			"$set": foodModel,
		}
		_, err := service.foodRepository.PutFood(filterPut, update)
		if err != nil {
			return nil, err
		}
	}

	purchase := models.Purchase{
		UserID:    userID,
		Date:      time.Now(),
		TotalCost: totalCost,
		Items:     purchaseItems,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := service.purchaseRepository.PostPurchase(&purchase)
	if err != nil {
		return nil, err
	}

	purchase.ID = result.InsertedID.(primitive.ObjectID)

	purchaseDTO := dtos.NewPurchaseUpdateDTO(&purchase)

	return purchaseDTO, nil
}
