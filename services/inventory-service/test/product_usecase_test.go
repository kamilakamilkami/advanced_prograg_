package test

import (
	"inventory-service/domain"
	// "inventory-service/internal/repository"
	"inventory-service/internal/usecase"
	"testing"
)

func TestDecreaseStockLogic(t *testing.T) {
	usecase := &usecase.ProductUsecase{} 

	t.Run("successfully decreases stock", func(t *testing.T) {
		product := &domain.Product{
			ID:    1,
			Name:  "Test Product",
			Stock: 10,
		}

		err := usecase.DecreaseStockLogic(product, 3)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if product.Stock != 7 {
			t.Errorf("expected stock to be 7, got %d", product.Stock)
		}
	})

	t.Run("fails if not enough stock", func(t *testing.T) {
		product := &domain.Product{
			ID:    2,
			Name:  "Another Product",
			Stock: 2,
		}

		err := usecase.DecreaseStockLogic(product, 5)
		if err == nil {
			t.Errorf("expected error due to insufficient stock, got nil")
		}
		if product.Stock != 2 {
			t.Errorf("stock should remain unchanged, expected 2, got %d", product.Stock)
		}
	})
}
