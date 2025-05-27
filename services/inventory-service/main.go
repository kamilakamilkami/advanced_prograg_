package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"log"

	// "github.com/gin-gonic/gin"
	grpc_hand "inventory-service/delivery/grpc"
	"proto/inventorypb"
	"google.golang.org/grpc"
	"net"
	"fmt"
	// "inventory-service/delivery/http"
	"inventory-service/internal/repository"
	"inventory-service/internal/usecase"
	"inventory-service/nats"
	"inventory-service/internal/redis"
)

func main() {

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

	db := client.Database("inventory_db")
	productCollection := db.Collection("products")
	redisClient := redis.NewClient()

    repo := repository.NewProductRepository(productCollection.Database())
    usecase := usecase.NewProductUsecase(repo, redisClient)

	grpcServer := grpc.NewServer()

	inventoryServiceServer := grpc_hand.NewInventoryServer(usecase)

	inventorypb.RegisterInventoryServiceServer(grpcServer, inventoryServiceServer)

	go nats.InitNATSConsumer(usecase)


	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("Inventary Service running on :5001")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	  
}
