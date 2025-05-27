package usecase

import (
    "inventory-service/domain"
    "inventory-service/internal/repository"
	"fmt"
	"github.com/redis/go-redis/v9"
	local_redis "inventory-service/internal/redis"
	"encoding/json"
	"time"
)

type ProductUsecase struct {
    repo *repository.ProductRepository
	redis *redis.Client
}

func NewProductUsecase(repo *repository.ProductRepository, redisClient *redis.Client) *ProductUsecase {
    return &ProductUsecase{
		repo: repo,
		redis: redisClient,
	}
}

func (uc *ProductUsecase) GetProduct(id int32) (*domain.Product, error) {
    key := fmt.Sprintf("product:%d", id)
    
    // Check cache
    val, err := uc.redis.Get(local_redis.Ctx, key).Result()
    if err == nil {
        var cachedProduct domain.Product
        json.Unmarshal([]byte(val), &cachedProduct)
        return &cachedProduct, nil
    }
	product, err := uc.repo.GetProductByID(id)
    if err != nil {
        return nil, err
    }
	jsonData, _ := json.Marshal(product)
    uc.redis.Set(local_redis.Ctx, key, jsonData, time.Minute*5)

    return product, nil
}

func (uc *ProductUsecase) CreateProduct(product domain.Product) (domain.Product, error) {
	id, err := uc.repo.CreateProduct(&product)
	if err != nil {
		return domain.Product{}, err
	}
	product.ID = int(id)
	return product, nil
}


func (uc *ProductUsecase) UpdateProduct(product *domain.Product) error {
	uc.redis.Del(local_redis.Ctx, fmt.Sprintf("product:%d", product.ID))
	return uc.repo.UpdateProduct(product)
}

func (uc *ProductUsecase) DeleteProduct(id int32) error {
	uc.redis.Del(local_redis.Ctx, fmt.Sprintf("product:%v", id))
	return uc.repo.DeleteProduct(id)
}
func (uc *ProductUsecase) GetAllProducts(name string, category, limit, offset int) ([]domain.Product, error) {
    return uc.repo.GetAllProducts(name, category, limit, offset)
}

func (u *ProductUsecase) DecreaseStockLogic(product *domain.Product, qty int32) error {
	if product.Stock < qty {
		return fmt.Errorf("not enough stock for product ID %d", product.ID)
	}
	product.Stock -= qty
	return nil
}

func (u *ProductUsecase) DecreaseStock(productID int, qty int32) error {
	product, err := u.GetProduct(int32(productID))
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	err = u.DecreaseStockLogic(product, qty)
	if err != nil {
		return err
	}

	err = u.UpdateProduct(product)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}
