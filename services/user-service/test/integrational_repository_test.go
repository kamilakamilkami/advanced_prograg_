package test

import (
	// "database/sql"
	// "fmt"
	"context"
	"log"
	"testing"
	"user-service/domain"
	"user-service/internal/repository"
	"user-service/internal/usecase"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateAndGetUserByID(t *testing.T) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("❌ Failed to ping MongoDB:", err)
	} else {
		log.Println("✅ Successfully connected to MongoDB")
	}

	db := client.Database("user_db")
	productCollection := db.Collection("users")

	repo := repository.NewUserRepository(productCollection.Database())

	user := &domain.User{
		UserID:   usecase.GenerateUserID(),
		Username: "testuser",
		Password: "hashedpassword123",
	}

	_, _ = productCollection.DeleteOne(context.TODO(), map[string]interface{}{"user_id": user.UserID})

	_, err = repo.Create(user)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	got, err := repo.GetUserByID(user.UserID)
	if err != nil {
		t.Fatalf("GetUserByID() error: %v", err)
	}

	if got.UserID != user.UserID || got.Username != user.Username || got.Password != user.Password {
		t.Errorf("Expected %+v, got %+v", user, got)
	}
}
