package repository

import (
	"fmt"
	"user-service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
    // "go.mongodb.org/mongo-driver/bson/primitive"
	"errors"
    "context"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryMongo struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepositoryMongo {
	return &UserRepositoryMongo{
		collection: db.Collection("users"),
	}
}


func (repo *UserRepositoryMongo) Create(user *domain.User) (string, error) {
	_, err := repo.collection.InsertOne(context.Background(), bson.M{
		"user_id":  user.UserID,
		"username": user.Username,
		"password": user.Password,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create user: %v", err)
	}
	return user.UserID, nil

}
func (repo *UserRepositoryMongo) Authenticate(username, password string) (string, error) {
	var user domain.User
	filter := bson.M{"username": username}

	err := repo.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("user not found")
			return "", fmt.Errorf("authentication failed: user not found")
		}
		fmt.Println("failed to find user")
		return "", fmt.Errorf("authentication failed: %v", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		fmt.Println("incorrect password")
		return "", fmt.Errorf("authentication failed: incorrect password")
	}

	return user.UserID, nil
}


func (repo *UserRepositoryMongo) GetUserByID(userID string) (*domain.User, error) {
	// var user domain.User
	// err := repo.db.QueryRow(`SELECT user_id, username, password FROM users WHERE user_id=$1`, userID).Scan(&user.UserID, &user.Username, &user.Password)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return nil, fmt.Errorf("user not found")
	// 	}
	// 	return nil, fmt.Errorf("failed to get user by ID: %v", err)
	// }
	// return &user, nil
	var user domain.User

	filter := bson.M{"user_id": userID}
	err := repo.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %v", err)
	}

	return &user, nil

}

