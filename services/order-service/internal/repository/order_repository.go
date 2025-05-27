package repository

import (
	"order-service/domain"
    "context"
	"go.mongodb.org/mongo-driver/mongo"
    // "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/bson"
	// "time"
	// "fmt"
	"errors"
	"time"
)

type OrderRepository struct {
    collection *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) *OrderRepository {
    return &OrderRepository{
        collection: db.Collection("orders"),
    }
}

func (r *OrderRepository) CreateOrder(order *domain.Order) error {
    ctx := context.Background()

	order.ID = int(time.Now().Unix())

	_, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return err
	}

	return nil

}
// func (r *OrderRepository) CreateOrder(order *domain.Order) error {
//     // 1. Начинаем сессию
//     session, err := r.collection.Database().Client().StartSession()
//     if err != nil {
//         return fmt.Errorf("failed to start session: %w", err)
//     }
//     defer session.EndSession(context.Background())

//     // 2. Транзакция
//     callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
//         order.ID = int(time.Now().Unix())

//         _, err := r.collection.InsertOne(sessCtx, order)
//         if err != nil {
//             return nil, fmt.Errorf("failed to insert order with transaction: %w", err)
//         }

//         return nil, nil
//     }

//     // 3. Выполнение транзакции
//     _, err = session.WithTransaction(context.Background(), callback)
//     if err != nil {
//         return fmt.Errorf("transaction failed: %w", err)
//     }

//     return nil
// }



func (r *OrderRepository) GetOrderByID(id int) (*domain.Order, error) {
    ctx := context.Background()
	var order domain.Order

	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil

}

func (r *OrderRepository) UpdateOrderStatus(id int, status string) error {
    ctx := context.Background()
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *OrderRepository) GetOrdersByUser(userID int) ([]domain.Order, error) {
    ctx := context.Background()
	filter := bson.M{"user_id": userID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []domain.Order
	for cursor.Next(ctx) {
		var order domain.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return orders, nil

}

func (r *OrderRepository) getOrderItems(orderID int) ([]domain.OrderItem, error) {
	order, err := r.GetOrderByID(orderID)
	if err != nil || order == nil {
		return nil, err
	}
	return order.Items, nil
}
