package grpc

import (
	"context"
	"inventory-service/domain"
	"inventory-service/internal/usecase"
	"proto/inventorypb"
	"fmt"
)

type InventoryServer struct {
	inventorypb.UnimplementedInventoryServiceServer
	Usecase *usecase.ProductUsecase
}

func NewInventoryServer(usecase *usecase.ProductUsecase) *InventoryServer {
	return &InventoryServer{Usecase: usecase}
}

func (s *InventoryServer) CreateProduct(ctx context.Context, req *inventorypb.CreateProductRequest) (*inventorypb.CreateProductResponse, error) {
	p := req.Product
	product := domain.Product{
		Name:       p.Name,
		Price:      p.Price,
		Stock:      p.Stock,
		CategoryID: p.CategoryId,
	}
	createdProduct, err := s.Usecase.CreateProduct(product)
	if err != nil {
		return nil, err
	}
	return &inventorypb.CreateProductResponse{Id: int32(createdProduct.ID)}, nil
}

func (s *InventoryServer) GetProduct(ctx context.Context, req *inventorypb.GetProductRequest) (*inventorypb.GetProductResponse, error) {
	fmt.Println(req)
	product, err := s.Usecase.GetProduct(req.Id)
	if err != nil {
		return nil, err
	}
	return &inventorypb.GetProductResponse{
		Product: &inventorypb.Product{
			Id:         int32(product.ID),
			Name:       product.Name,
			Price:      product.Price,
			Stock:      product.Stock,
			CategoryId: product.CategoryID,
		},
	}, nil
}

func (s *InventoryServer) ListProducts(ctx context.Context, req *inventorypb.ListProductsRequest) (*inventorypb.ListProductsResponse, error) {
	products, err := s.Usecase.GetAllProducts(req.Name, int(req.Category), int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	var pbProducts []*inventorypb.Product
	for _, product := range products {
		pbProducts = append(pbProducts, &inventorypb.Product{
			Id:         int32(product.ID),
			Name:       product.Name,
			Price:      product.Price,
			Stock:      product.Stock,
			CategoryId: product.CategoryID,
		})
	}

	return &inventorypb.ListProductsResponse{Products: pbProducts}, nil
}

func (s *InventoryServer) UpdateProduct(ctx context.Context, req *inventorypb.UpdateProductRequest) (*inventorypb.UpdateProductResponse, error) {
	p := req.Product
	product := &domain.Product{ 
		ID:         int(p.Id),
		Name:       p.Name,
		Price:      p.Price,
		Stock:      p.Stock,
		CategoryID: p.CategoryId,
	}
	fmt.Println("Update")
	fmt.Println(product)
	err := s.Usecase.UpdateProduct(product)
	if err != nil {
		return nil, err
	}

	return &inventorypb.UpdateProductResponse{
		Product: &inventorypb.Product{
			Id:         int32(product.ID),
			Name:       product.Name,
			Price:      product.Price,
			Stock:      product.Stock,
			CategoryId: product.CategoryID,
		},
	}, nil
}


func (s *InventoryServer) DeleteProduct(ctx context.Context, req *inventorypb.DeleteProductRequest) (*inventorypb.DeleteProductResponse, error) {
	err := s.Usecase.DeleteProduct(req.Id)
	if err != nil {
		return nil, err
	}
	return &inventorypb.DeleteProductResponse{Success: true}, nil
}
