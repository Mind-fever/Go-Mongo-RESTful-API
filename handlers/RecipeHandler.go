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

type RecipeHandler struct {
	recipeService services.RecipeService
}

func NewRecipeHandler(recipeService services.RecipeService) *RecipeHandler {
	return &RecipeHandler{
		recipeService: recipeService,
	}
}

func (handler *RecipeHandler) GetRecipesByFilter(c *gin.Context) {
	log.Println("Handler: GetRecipesByFilter")
	var filter dtos.RecipeFilterDTO

	userInfo := utils.GetUserInfoFromContext(c)

	filter = dtos.RecipeFilterDTO{
		UserID:      userInfo.UserId,
		Use:         c.Query("meal_time"),
		ProductType: c.Query("product_type"),
		ProductName: c.Query("product_name"),
	}

	log.Printf("Filter parameters: %+v", filter)

	recipes, err := handler.recipeService.GetRecipesByFilter(filter)
	if err != nil {
		log.Printf("Error retrieving recipes: %v", err)
		c.Error(err)
		return
	}

	log.Printf("Recipes retrieved: %d", len(recipes))
	c.JSON(http.StatusOK, recipes)
}

func (handler *RecipeHandler) GetRecipeByID(c *gin.Context) {
	recipeID := c.Param("id")

	userInfo := utils.GetUserInfoFromContext(c)

	recipe, err := handler.recipeService.GetRecipeByID(userInfo.UserId, recipeID)
	if err != nil {
		log.Printf("Error retrieving recipe: %v", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func (handler *RecipeHandler) PostRecipe(c *gin.Context) {
	var recipe dtos.RecipeDTO

	userInfo := utils.GetUserInfoFromContext(c)

	if err := c.ShouldBindJSON(&recipe); err != nil {
		log.Printf("Error binding recipe DTO: %v", err)
		c.Error(errors.ErrBindingFailed)
		return
	}

	createdRecipeDTO, err := handler.recipeService.PostRecipe(userInfo.UserId, &recipe)
	if err != nil {
		log.Printf("Error creating recipe: %v", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, createdRecipeDTO)
}

func (handler *RecipeHandler) DeleteRecipe(c *gin.Context) {
	log.Println("Handler: DeleteRecipe")
	recipeID := c.Param("id")

	userInfo := utils.GetUserInfoFromContext(c)

	result, err := handler.recipeService.DeleteRecipe(userInfo.UserId, recipeID)
	if err != nil {
		log.Printf("Error deleting recipe: %v", err)
		c.Error(err)
		return
	}
	log.Println("Recipe deleted successfully")
	c.JSON(http.StatusOK, result)
}

func (handler *RecipeHandler) PutRecipe(c *gin.Context) {
	log.Println("Handler: PutRecipe")
	var recipe dtos.RecipeUpdateDTO
	userInfo := utils.GetUserInfoFromContext(c)

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.Error(errors.ErrBindingFailed)
		return
	}

	result, err := handler.recipeService.PutRecipe(userInfo.UserId, &recipe)
	if err != nil {
		log.Printf("Error updating recipe: %v", err)
		c.Error(err)
		return
	}
	log.Println("Recipe updated successfully")
	c.JSON(http.StatusCreated, result)
}
