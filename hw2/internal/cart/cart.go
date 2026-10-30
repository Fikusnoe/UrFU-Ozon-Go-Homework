package cart

import (
	"context"
	"errors"
	"fmt"
	"sync"

	pb "hw2/pkg/api/product"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrCartNotFound = errors.New("Cart not found")
var ErrProductNotFound = errors.New("product not found")

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type UserID int64

type CartItem struct {
	Sku   int64  `json:"sku"`
	Count uint32 `json:"count"`
	Price uint32 `json:"price"`
	Name  string `json:"name"`
}

type Cart struct {
	Items map[int64]CartItem
}

type CartService struct {
	CartStorage Storage
	Client      pb.ProductServiceClient
	mu          *sync.RWMutex
}

type rpcResult struct {
	data *pb.GetProductResponse
	err  error
}

func NewCartService(cs Storage, client pb.ProductServiceClient) *CartService {
	return &CartService{CartStorage: cs, Client: client, mu: &sync.RWMutex{}}
}

func ValidateInput(userID UserID, sku int64, count uint32) error {
	err := ValidateUserId(userID)
	if err != nil {
		return err
	}

	if sku <= 0 {
		return &ValidationError{Message: "SKU должен быть натуральным числом (больше нуля)"}
	} else if count <= uint32(0) {
		return &ValidationError{Message: "Количество должно быть натуральным числом (больше нуля)"}
	}

	return nil
}

func ValidateUserId(userID UserID) error {
	if userID <= 0 {
		return &ValidationError{"Идентификатор пользователя должен быть натуральным числом (больше нуля)"}
	}

	return nil
}

func (cs *CartService) MakeProductRPCRequest(ctx context.Context, sku int64) (*pb.GetProductResponse, error) {
	request := pb.GetProductRequest{Sku: sku}
	response, err := cs.Client.GetProduct(ctx, &request)
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return nil, ErrProductNotFound

		case codes.InvalidArgument:

			return nil,
				&ValidationError{
					Message: "SKU должен быть натуральным числом (больше нуля)",
				}

		default:
			return nil, fmt.Errorf("get product: %w", err)
		}
	}

	return response, nil
}
