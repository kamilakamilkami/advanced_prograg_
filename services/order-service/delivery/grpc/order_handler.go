package grpc

import (
	"context"
	"fmt"
	// "order-service/domain"
	"order-service/internal/usecase"
	"proto/orderpb"
)

type OrderServer struct {
	orderpb.UnimplementedOrderServiceServer
	Usecase *usecase.OrderUsecase
}

func NewOrderServer(usecase *usecase.OrderUsecase) *OrderServer {
	return &OrderServer{Usecase: usecase}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	fmt.Println((req))
	order, err := s.Usecase.CreateOrderFromRequest(req)
	if err != nil {
		return nil, err
	}
	fmt.Println(order)
	respItems := []*orderpb.OrderItem{}
	for _, i := range order.Items {
		respItems = append(respItems, &orderpb.OrderItem{
			ProductId: int32(i.ProductID),
			Quantity:  int32(i.Quantity),
		})
	}
	fmt.Println(respItems)
	return &orderpb.CreateOrderResponse{
		Order: &orderpb.Order{
			Id:     int32(order.ID),
			UserId: int32(order.UserID),
			Status: order.Status,
			Items:  respItems,
		},
	}, nil
}


func (s *OrderServer) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	order, err := s.Usecase.GetOrderByIDPB(int(req.Id))
	if err != nil {
		return nil, err
	}
	return &orderpb.GetOrderResponse{Order: order}, nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*orderpb.UpdateOrderStatusResponse, error) {
	err := s.Usecase.UpdateOrderStatus(int(req.Id), req.Status)
	if err != nil {
		return nil, err
	}
	return &orderpb.UpdateOrderStatusResponse{Success: true}, nil
}

func (s *OrderServer) GetOrdersByUser(ctx context.Context, req *orderpb.GetOrdersByUserRequest) (*orderpb.GetOrdersByUserResponse, error) {
	orders, err := s.Usecase.GetOrdersByUserPB(int(req.UserId))
	if err != nil {
		return nil, err
	}
	return &orderpb.GetOrdersByUserResponse{Orders: orders}, nil
}

