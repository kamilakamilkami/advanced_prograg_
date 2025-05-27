package usecase

import (
	"fmt"
	"order-service/domain"
	// "context"
	"order-service/internal/repository"
	"order-service/nats"
	"proto/orderpb"
	"log"
	"github.com/redis/go-redis/v9"
	local_redis "order-service/internal/redis"
	"encoding/json"
	"time"
)

type OrderUsecase struct {
    repo *repository.OrderRepository
	publisher *nats.NatsPublisher
	redis *redis.Client
}

func NewOrderUsecase(repo *repository.OrderRepository, publisher *nats.NatsPublisher, redisClient *redis.Client) *OrderUsecase {
    return &OrderUsecase{
		repo: repo,
		publisher: publisher,
		redis: redisClient,
	}
}
func (u *OrderUsecase) CreateOrder(order *domain.Order) error {
    order.Status = "pending"
    return u.repo.CreateOrder(order)
}

func (uc *OrderUsecase) GetOrderByID(id int) (*domain.Order, error) {
    key := fmt.Sprintf("order:%d", id)

	val, err := uc.redis.Get(local_redis.Ctx, key).Result()
    if err == nil {
        var cachedProduct domain.Order
        json.Unmarshal([]byte(val), &cachedProduct)
        return &cachedProduct, nil
    }
    order, err := uc.repo.GetOrderByID(id)
	if err != nil {
        return nil, err
    }
	jsonData, _ := json.Marshal(order)
    uc.redis.Set(local_redis.Ctx, key, jsonData, time.Minute*5)

    return order, nil
}

func (uc *OrderUsecase) UpdateOrderStatus(id int, status string) error {
	uc.redis.Del(local_redis.Ctx, fmt.Sprintf("order:%d", id))
    return uc.repo.UpdateOrderStatus(id, status)
}

func (uc *OrderUsecase) GetOrdersByUser(userID int) ([]domain.Order, error) {
    return uc.repo.GetOrdersByUser(userID)
}

func (u *OrderUsecase) CreateOrderFromRequest(req *orderpb.CreateOrderRequest) (*domain.Order, error) {
	order := u.BuildOrderFromRequest(req)

	err := u.repo.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	for _, item := range order.Items {
		event := nats.OrderCreatedEvent{
			ID:        fmt.Sprintf("%d", order.ID),
			ProductID: fmt.Sprintf("%d", item.ProductID),
			Quantity:  item.Quantity,
		}

		err := u.publisher.PublishOrderCreated(event)
		if err != nil {
			log.Println("Warning: Failed to publish NATS message:", err)
		}
	}

	return order, nil
}

func (u *OrderUsecase) BuildOrderFromRequest(req *orderpb.CreateOrderRequest) *domain.Order {
	order := &domain.Order{
		UserID: int(req.UserId),
		Status: "pending",
	}

	for _, i := range req.Items {
		order.Items = append(order.Items, domain.OrderItem{
			ProductID: int(i.ProductId),
			Quantity:  int(i.Quantity),
		})
	}

	return order
}

func (u *OrderUsecase) GetOrderByIDPB(id int) (*orderpb.Order, error) {
	key := fmt.Sprintf("order:%d", id)

	val, err := u.redis.Get(local_redis.Ctx, key).Result()
    if err == nil {
        var cachedOrder domain.Order
        json.Unmarshal([]byte(val), &cachedOrder)
		orderItems := []*orderpb.OrderItem{}
		for _, i := range cachedOrder.Items {	
			orderItems = append(orderItems, &orderpb.OrderItem{
				ProductId: int32(i.ProductID),
				Quantity:  int32(i.Quantity),
			})
		}
        return &orderpb.Order{
			Id:     int32(cachedOrder.ID),
			UserId: int32(cachedOrder.UserID),
			Status: cachedOrder.Status,
			Items:  orderItems,
		}, nil
    }
	order, err := u.repo.GetOrderByID(id)
	if err != nil {
		return nil, err
	}
	
	jsonData, _ := json.Marshal(order)
    u.redis.Set(local_redis.Ctx, key, jsonData, time.Minute*5)

	items := []*orderpb.OrderItem{}
	for _, i := range order.Items {
		items = append(items, &orderpb.OrderItem{
			ProductId: int32(i.ProductID),
			Quantity:  int32(i.Quantity),
		})
	}

	return &orderpb.Order{
		Id:     int32(order.ID),
		UserId: int32(order.UserID),
		Status: order.Status,
		Items:  items,
	}, nil
}

func (u *OrderUsecase) GetOrdersByUserPB(userID int) ([]*orderpb.Order, error) {
	orders, err := u.repo.GetOrdersByUser(userID)
	if err != nil {
		return nil, err
	}

	var pbOrders []*orderpb.Order
	for _, o := range orders {
		items := []*orderpb.OrderItem{}
		for _, i := range o.Items {
			items = append(items, &orderpb.OrderItem{
				ProductId: int32(i.ProductID),
				Quantity:  int32(i.Quantity),
			})
		}
		pbOrders = append(pbOrders, &orderpb.Order{
			Id:     int32(o.ID),
			UserId: int32(o.UserID),
			Status: o.Status,
			Items:  items,
		})
	}

	return pbOrders, nil
}
