package repositories

import (
	"context"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/errors"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipeRepositoryInterface interface {
	GetRecipes(userID string) ([]*models.Recipe, error)
	GetRecipesByFilter(filters map[string]interface{}) ([]*models.Recipe, error)
	PostRecipe(recipe *models.Recipe) (*mongo.InsertOneResult, error)
	PutRecipe(filters bson.M, update bson.M) (*mongo.UpdateResult, error)
	DeleteRecipe(filters bson.M) (*mongo.DeleteResult, error)
	GetRecipeByID(filters bson.M) (*models.Recipe, error)
}

type RecipeRepository struct {
	db DB
}

func (repository *RecipeRepository) getRecipeCollection() *mongo.Collection {
	return getRecipeCollection(repository.db)
}

func NewRecipeRepository(db DB) *RecipeRepository {
	return &RecipeRepository{
		db: db,
	}
}

func (repository *RecipeRepository) GetRecipes(userID string) ([]*models.Recipe, error) {
	collection := repository.getRecipeCollection()
	filter := bson.M{"user_id": userID}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := cursor.Close(context.Background()); err != nil {
		}
	}()

	var recipes []*models.Recipe
	for cursor.Next(context.Background()) {
		var recipe models.Recipe
		err := cursor.Decode(&recipe)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, &recipe)
	}
	return recipes, err
}

func (repository RecipeRepository) GetRecipesByFilter(filters map[string]interface{}) ([]*models.Recipe, error) {
	collection := repository.getRecipeCollection()

	cursor, err := collection.Find(context.Background(), filters)
	if err != nil {
		return nil, errors.ErrNoRecipesFound
	}

	defer cursor.Close(context.Background())

	var recipes []*models.Recipe
	for cursor.Next(context.Background()) {
		var recipe models.Recipe
		if err := cursor.Decode(&recipe); err != nil {
			return nil, errors.ErrCursorDecodeFailed
		}
		recipes = append(recipes, &recipe)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.ErrCursorIterationFailed
	}

	return recipes, nil
}

func (repository RecipeRepository) PostRecipe(recipe *models.Recipe) (*mongo.InsertOneResult, error) {
	collection := repository.getRecipeCollection()

	result, err := collection.InsertOne(context.Background(), recipe)
	if err != nil {
		return nil, errors.ErrDatabasePost
	}

	return result, nil
}

func (repository RecipeRepository) PutRecipe(filters bson.M, update bson.M) (*mongo.UpdateResult, error) {
	collection := repository.getRecipeCollection()

	result, err := collection.UpdateOne(context.Background(), filters, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repository RecipeRepository) DeleteRecipe(filters bson.M) (*mongo.DeleteResult, error) {
	collection := repository.getRecipeCollection()

	result, err := collection.DeleteOne(context.TODO(), filters)
	if err != nil {
		return nil, errors.ErrDatabaseDelete
	}

	return result, nil
}

func (repository *RecipeRepository) GetRecipeByID(filters bson.M) (*models.Recipe, error) {
	collection := repository.getRecipeCollection()
	var recipe models.Recipe

	err := collection.FindOne(context.TODO(), filters).Decode(&recipe)
	if err != nil {
		return nil, errors.ErrDatabaseGet
	}

	return &recipe, nil
}
