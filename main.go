package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/clients"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/handlers"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/middlewares"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/repositories"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/services"
)

var (
	foodHandler     *handlers.FoodHandler
	recipeHandler   *handlers.RecipeHandler
	purchaseHandler *handlers.PurchaseHandler
	router          *gin.Engine
)

func main() {
	router = gin.Default()
	router.Use(middlewares.CORSMiddleware()) // Usar el middleware de la carpeta middlewares
	router.Use(middlewares.ErrorHandler)
	// Initialize dependencies
	dependencies()
	// Map routes
	mappingRoutes()

	log.Println("Server running on port 8080")
	router.Run(":8080")
}

func mappingRoutes() {
	// External API client
	var authClient clients.AuthClientInterface
	authClient = clients.NewAuthClient()
	// Authentication middleware
	authMiddleware := middlewares.NewAuthMiddleware(authClient)

	group := router.Group("/foods")
	group.Use(authMiddleware.ValidateToken)
	group.GET("/", foodHandler.GetFoods)
	group.POST("/", foodHandler.PostFood)
	group.PUT("/:id", foodHandler.PutFood)
	group.DELETE("/:id", foodHandler.DeleteFood)

	group = router.Group("/stock")
	group.Use(authMiddleware.ValidateToken)
	group.GET("/", foodHandler.GetLowStockFoods)

	group = router.Group("/purchases")
	group.Use(authMiddleware.ValidateToken)
	group.GET("/", purchaseHandler.GetPurchases)
	group.POST("/", purchaseHandler.PostPurchase)

	group = router.Group("/recipes")
	group.Use(authMiddleware.ValidateToken)
	group.GET("/filter", recipeHandler.GetRecipesByFilter)
	group.GET("/:id", recipeHandler.GetRecipeByID)
	group.POST("/", recipeHandler.PostRecipe)
	group.PUT("/:id", recipeHandler.PutRecipe)
	group.DELETE("/:id", recipeHandler.DeleteRecipe)

}

// Initialize dependencies
func dependencies() {
	// Initialize MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}
	log.Println("MongoDB client created")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB")
	db := repositories.NewMongoDB()
	// Initialize food repository
	foodRepo := repositories.NewFoodRepository(db)
	recipeRepo := repositories.NewRecipeRepository(db)
	purchaseRepo := repositories.NewPurchaseRepository(db)
	log.Println("Repositories initialized")

	foodService := services.NewFoodService(foodRepo, recipeRepo)
	if foodService == nil {
		log.Fatalf("Failed to initialize food service")
	}
	log.Println("Food service initialized")

	recipeService := services.NewRecipeService(recipeRepo, foodRepo)
	if recipeService == nil {
		log.Fatalf("Failed to initialize recipe service")
	}
	log.Println("Recipe service initialized")

	purchaseService := services.NewPurchaseService(purchaseRepo, foodRepo, foodService)
	if purchaseService == nil {
		log.Fatalf("Failed to initialize purchase service")
	}
	log.Println("Purchase service initialized")

	foodHandler = handlers.NewFoodHandler(foodService)
	if foodHandler == nil {
		log.Fatalf("Failed to initialize food handler")
	}
	log.Println("Food handler initialized")

	recipeHandler = handlers.NewRecipeHandler(*recipeService)
	if recipeHandler == nil {
		log.Fatalf("Failed to initialize recipe handler")
	}
	log.Println("Recipe handler initialized")

	purchaseHandler = handlers.NewPurchaseHandler(purchaseService)
	if purchaseHandler == nil {
		log.Fatalf("Failed to initialize purchase handler")
	}
	log.Println("Purchase handler initialized")
}
