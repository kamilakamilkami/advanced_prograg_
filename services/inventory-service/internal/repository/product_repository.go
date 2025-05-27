package repository

import (
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
    // "go.mongodb.org/mongo-driver/bson/primitive"
    "context"
	"errors"
    "inventory-service/domain"
	"fmt"
)

type ProductRepository struct {
    collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
    return &ProductRepository{
        collection: db.Collection("products"),
    }
}

func (r *ProductRepository) GetProductByID(id int32) (*domain.Product, error) {
	filter := bson.M{"_id": id}

	var product domain.Product
	err := r.collection.FindOne(context.Background(), filter).Decode(&product)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}

	return &product, nil
}



func (r *ProductRepository) CreateProduct(product *domain.Product) (int32, error) {
    var last domain.Product
	opts := options.FindOne().SetSort(bson.D{{"_id", -1}})

	err := r.collection.FindOne(context.Background(), bson.M{}, opts).Decode(&last)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			product.ID = 1
		} else {
			return 0, err
		}
	} else {
		product.ID = last.ID + 1
	}


	_, err = r.collection.InsertOne(context.Background(), product)
	if err != nil {
		return 0, err
	}

	return int32(product.ID), nil
}

func (r *ProductRepository) UpdateProduct(product *domain.Product) error {
    filter := bson.M{"_id": product.ID}
	update := bson.M{"$set": bson.M{
		"name":        product.Name,
		"price":       product.Price,
		"stock":       product.Stock,
		"category_id": product.CategoryID,
	}}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *ProductRepository) DeleteProduct(id int32) error {
    _, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *ProductRepository) GetAllProducts(name string, category, limit, offset int) ([]domain.Product, error) {
    filter := bson.M{}
	if category != 0 {
		filter["category_id"] = category
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	cursor, err := r.collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var products []domain.Product
	for cursor.Next(context.Background()) {
		var p domain.Product
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil

}