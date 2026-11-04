package fake

import (
	"context"
	"hw2/internal/product"
	pb "hw2/pkg/api/product"

	"google.golang.org/grpc"
)

type ProductServiceClient struct {
	storage ProductStorage
}

func NewProductServiceClient(storage ProductStorage) *ProductServiceClient {
	return &ProductServiceClient{
		storage: storage,
	}
}

func (c *ProductServiceClient) GetProduct(ctx context.Context, in *pb.GetProductRequest, opts ...grpc.CallOption) (*pb.GetProductResponse, error) {
	svc, err := product.NewProductService(&c.storage)
	if err != nil {
		return nil, err
	}

	out, err := svc.GetProduct(ctx, in)
	if err != nil {
		return nil, err
	}

	return out, nil
}
