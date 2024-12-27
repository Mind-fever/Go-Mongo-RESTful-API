package repositories

import (
	"context"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/errors"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PurchaseRepositoryInterface interface {
	GetPurchases(filter bson.M) ([]*models.Purchase, error)
	PostPurchase(purchase *models.Purchase) (*mongo.InsertOneResult, error)
}

type PurchaseRepository struct {
	db DB
}

func (repository *PurchaseRepository) getPurchaseCollection() *mongo.Collection {
	return getPurchaseCollection(repository.db)
}

func NewPurchaseRepository(db DB) *PurchaseRepository {
	return &PurchaseRepository{
		db: db,
	}
}

func (repository PurchaseRepository) GetPurchases(filter bson.M) ([]*models.Purchase, error) {
	collection := repository.getPurchaseCollection()

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, errors.ErrDatabaseGet
	}

	defer cursor.Close(context.Background())

	var purchases []*models.Purchase
	for cursor.Next(context.Background()) {
		var purchase models.Purchase
		err := cursor.Decode(&purchase)
		if err != nil {
			return nil, errors.ErrCursorDecodeFailed
		}
		purchases = append(purchases, &purchase)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.ErrCursorIterationFailed
	}

	return purchases, nil
}

func (repository PurchaseRepository) PostPurchase(purchase *models.Purchase) (*mongo.InsertOneResult, error) {
	collection := repository.getPurchaseCollection()

	result, err := collection.InsertOne(context.Background(), purchase)
	if err != nil {
		return nil, errors.ErrDatabasePost
	}

	return result, nil
}
