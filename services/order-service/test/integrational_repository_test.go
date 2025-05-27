package test

import (
	"context"
	"testing"
	"time"

	"order-service/domain"
	"order-service/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateOrder(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("orders_db")
	orderCollection := db.Collection("orders")

	repo := repository.NewOrderRepository(db)

	order := &domain.Order{
		UserID: 1,
		Items: []domain.OrderItem{
			{ProductID: 101, Quantity: 2},
			{ProductID: 102, Quantity: 1},
		},
	}

	order.ID = int(time.Now().Unix())

	// Удаление возможных старых данных
	_, _ = orderCollection.DeleteMany(ctx, bson.M{"user_id": order.UserID, "status": "pending"})

	// Вставка заказа
	err = repo.CreateOrder(order)
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}


	// Получение заказа
	storedOrder, err := repo.GetOrderByID(order.ID)
	if err != nil {
		t.Fatalf("GetOrderByID failed: %v", err)
	}

	if storedOrder.UserID != order.UserID {
		t.Errorf("Expected UserID %d, got %d", order.UserID, storedOrder.UserID)
	}

	if storedOrder.Status != order.Status {
		t.Errorf("Expected Status %s, got %s", order.Status, storedOrder.Status)
	}

	if len(storedOrder.Items) != len(order.Items) {
		t.Fatalf("Expected %d items, got %d", len(order.Items), len(storedOrder.Items))
	}

	for i, item := range storedOrder.Items {
		if item.ProductID != order.Items[i].ProductID {
			t.Errorf("Item %d: expected ProductID %d, got %d", i, order.Items[i].ProductID, item.ProductID)
		}
		if item.Quantity != order.Items[i].Quantity {
			t.Errorf("Item %d: expected Quantity %d, got %d", i, order.Items[i].Quantity, item.Quantity)
		}
	}

	// Очистка данных
	_, err = orderCollection.DeleteOne(ctx, bson.M{"_id": order.ID})
	if err != nil {
		t.Logf("Failed to delete test order with ID %v: %v", order.ID, err)
	}
}
