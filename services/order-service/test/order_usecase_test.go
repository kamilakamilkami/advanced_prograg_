package test

import (
	"proto/orderpb"
	"order-service/internal/usecase"
	"testing"
)

func TestBuildOrderFromRequest(t *testing.T) {
	u := usecase.OrderUsecase{} 

	req := &orderpb.CreateOrderRequest{
		UserId: 42,
		Items: []*orderpb.OrderItem{
			{ProductId: 1, Quantity: 2},
			{ProductId: 3, Quantity: 4},
		},
	}

	order := u.BuildOrderFromRequest(req)

	if order.UserID != int(req.UserId) {
		t.Errorf("expected UserID %d, got %d", req.UserId, order.UserID)
	}

	if order.Status != "pending" {
		t.Errorf("expected Status 'pending', got '%s'", order.Status)
	}

	if len(order.Items) != len(req.Items) {
		t.Fatalf("expected %d items, got %d", len(req.Items), len(order.Items))
	}

	for i, item := range order.Items {
		if item.ProductID != int(req.Items[i].ProductId) {
			t.Errorf("item %d: expected ProductID %d, got %d", i, req.Items[i].ProductId, item.ProductID)
		}
		if item.Quantity != int(req.Items[i].Quantity) {
			t.Errorf("item %d: expected Quantity %d, got %d", i, req.Items[i].Quantity, item.Quantity)
		}
	}
}
