package product

import (
	"context"
	"fmt"
	"hw2/internal/domain"
	
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "hw2/pkg/api/product"
)

type ProductService struct {
	Storage  Storage
	Products map[int64]domain.Product

	pb.UnimplementedProductServiceServer
}

func NewProductService(storage Storage) (*ProductService, error) {
	loadedProducts, err := storage.LoadProducts()
	if err != nil {
		return nil, fmt.Errorf("Failed to load Products: %w", err)
	}

	productService := ProductService{
		Storage: storage,
		Products: loadedProducts,
	}

	return &productService, nil
}

func (ps *ProductService) GetProduct(ctx context.Context, request *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	
	if request.Sku <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Идентификатор товара должен быть натуральным числом (больше нуля)")
	}

	product, ok := ps.Products[request.Sku]
	if !ok {
		return nil, status.Error(codes.NotFound, "Информацию можно получить только по существующему товару")
	}

	return &pb.GetProductResponse{Sku: request.Sku, Name: product.Name, Price: int32(product.Price)}, nil
}
