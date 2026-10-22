package product

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "hw2/pkg/api/product"
)

type Price uint32

type Product struct {
	Name  string `json:"name"`
	Price Price `json:"price"`
}

type ProductService struct {
	Storage  Storage
	Products map[int64]Product

	pb.UnimplementedProductServiceServer
}

func NewProductService(storage Storage) *ProductService {
	ProductService := ProductService{
		Storage: storage,
	}

	return &ProductService
}

func (ps *ProductService) LoadProducts() error {
	loadedProducts := map[int64]Product{}

	loadedProducts, err := ps.Storage.LoadProducts()
	if err != nil {
		return err
	}

	ps.Products = loadedProducts
	return nil
}

func (ps *ProductService) GetProduct(context context.Context, request *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	if request.Sku <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Идентификатор товара должен быть натуральным числом (больше нуля)")
	}

	product, ok := ps.Products[request.Sku]
	if !ok {
		return nil, status.Error(codes.NotFound, "Информацию можно получить только по существующему товару")
	}

	return &pb.GetProductResponse{Sku: request.Sku, Name: product.Name, Price: int32(product.Price)}, nil
}
