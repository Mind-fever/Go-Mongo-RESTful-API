package handlers

import (
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/errors"
	"log"
	"net/http"

	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/dtos"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/services"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/utils"
	"github.com/gin-gonic/gin"
)

type FoodHandler struct {
	foodService services.FoodServiceInterface
}

func NewFoodHandler(foodService services.FoodServiceInterface) *FoodHandler {
	return &FoodHandler{
		foodService: foodService,
	}
}

func (handler *FoodHandler) GetFoods(c *gin.Context) {
	log.Println("Handler: GetFoods")
	userInfo := utils.GetUserInfoFromContext(c)

	foods, err := handler.foodService.GetFoods(userInfo.UserId)
	if err != nil {
		log.Printf("Error fetching foods: %v", err)
		c.Error(err)
		return
	}

	log.Printf("Foods retrieved: %d", len(foods))
	c.JSON(http.StatusOK, foods)
}

func (handler *FoodHandler) GetFoodByID(c *gin.Context) {
	log.Println("Handler: GetFoodByID")
	foodID := c.Param("id")
	userInfo := utils.GetUserInfoFromContext(c)

	food, err := handler.foodService.GetFoodByID(userInfo.UserId, foodID)
	if err != nil {
		log.Printf("Error fetching food: %v", err)
		c.Error(err)
		return
	}
	log.Printf("Food retrieved: %v", food)
	c.JSON(http.StatusOK, food)
}

func (handler *FoodHandler) GetLowStockFoods(c *gin.Context) {
	log.Println("Handler: GetLowStockFoods")
	var filter dtos.LowStockFilterDTO
	userInfo := utils.GetUserInfoFromContext(c)

	if err := c.ShouldBindQuery(&filter); err != nil {
		log.Printf("Query binding error: %v", err)
		c.Error(errors.ErrBindingFailed)
		return
	}

	filter.UserID = userInfo.UserId

	foods, err := handler.foodService.GetLowStockFoods(&filter)
	if err != nil {
		log.Printf("Error retrieving low stock foods: %v", err)
		c.Error(err)
		return
	}

	log.Printf("Fetched %d low stock foods", len(foods))
	c.JSON(http.StatusOK, foods)
}
func (handler *FoodHandler) PutFood(c *gin.Context) {
	log.Println("Handler: PutFood")
	var food dtos.FoodUpdateDTO
	userInfo := utils.GetUserInfoFromContext(c)

	if err := c.ShouldBindJSON(&food); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.Error(errors.ErrBindingFailed)
		return
	}

	result, err := handler.foodService.PutFood(&food, userInfo.UserId)

	if err != nil {
		log.Printf("Error updating food: %v", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (handler *FoodHandler) PostFood(c *gin.Context) {
	var food dtos.FoodDTO

	userInfo := utils.GetUserInfoFromContext(c)

	if err := c.ShouldBindJSON(&food); err != nil {
		log.Printf("Error binding food DTO: %v", err)
		c.Error(errors.ErrBindingFailed)
		return
	}

	createdFoodDTO, err := handler.foodService.PostFood(&food, userInfo.UserId)
	if err != nil {
		log.Printf("Error creating food: %v", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, createdFoodDTO)
}

func (handler *FoodHandler) DeleteFood(c *gin.Context) {
	foodID := c.Param("id")
	userInfo := utils.GetUserInfoFromContext(c)

	result, err := handler.foodService.DeleteFood(userInfo.UserId, foodID)
	if err != nil {
		log.Printf("Error deleting food: %v", err)
		c.Error(err)
		return
	}

	log.Println("[handler:FoodHandler][method:DeleteFood][user:", userInfo, "]")
	c.JSON(http.StatusOK, result)
}
