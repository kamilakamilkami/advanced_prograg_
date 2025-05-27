package test

import (
	"context"
	"testing"
	"time"

	"inventory-service/domain"
	"inventory-service/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateAndGetProduct(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("inventory_db")
	productCollection := db.Collection("products")
	repo := repository.NewProductRepository(db)

	input := &domain.Product{
		ID:         int(time.Now().Unix()), // или UUID, если используешь строки
		Name:       "Test Product",
		Price:      100,
		Stock:      10,
		CategoryID: 1,
	}

	// Удаление возможного дубликата
	_, _ = productCollection.DeleteMany(ctx, bson.M{"_id": input.ID})

	_, err = repo.CreateProduct(input)
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}

	got, err := repo.GetProductByID(int32(input.ID))
	if err != nil {
		t.Fatalf("GetProductByID failed: %v", err)
	}

	if got.Name != input.Name || got.Price != input.Price || got.Stock != input.Stock || got.CategoryID != input.CategoryID {
		t.Errorf("Product mismatch.\nExpected: %+v\nGot: %+v", input, got)
	}

	// Очистка
	_, err = productCollection.DeleteOne(ctx, bson.M{"_id": input.ID})
	if err != nil {
		t.Logf("Failed to delete test product with ID %d: %v", input.ID, err)
	}
}
