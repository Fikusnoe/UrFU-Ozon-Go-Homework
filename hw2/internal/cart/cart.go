package cart

import (
	"context"
	"errors"
	"fmt"
	"slices"
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

func NewCartService(cs Storage, client pb.ProductServiceClient) *CartService {
	return &CartService{CartStorage: cs, Client: client, mu: &sync.RWMutex{}}
}

func (cs *CartService) AddItem(ctx context.Context, userID UserID, sku int64, count uint32) error {
	err := ValidateInput(userID, sku, count)
	if err != nil {
		return err
	}

	cs.mu.Lock()
	defer cs.mu.Unlock()

	userCart, err := cs.CartStorage.LoadCart(userID)
	if err != nil {
		if errors.Is(err, ErrCartNotFound) {
			userCart = Cart{Items: map[int64]CartItem{}}
		} else {
			return err
		}
	}

	request := pb.GetProductRequest{Sku: sku}
	response, err := cs.Client.GetProduct(ctx, &request)
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return ErrProductNotFound

		case codes.InvalidArgument:
			return &ValidationError{
				Message: "SKU должен быть натуральным числом (больше нуля)",
			}

		default:
			return fmt.Errorf("get product: %w", err)
		}
	}

	item, exist := userCart.Items[sku]
	if exist {
		item.Count += count
		userCart.Items[sku] = item
	} else {
		userCart.Items[sku] = CartItem{Sku: response.Sku,
			Count: count,
			Price: uint32(response.Price),
			Name:  response.Name,
		}
	}

	err = cs.CartStorage.SaveCart(userID, userCart)
	if err != nil {
		return fmt.Errorf("failed to save cart: %w", err)
	}

	return nil
}

func (cs *CartService) GetCart(userID UserID) ([]CartItem, uint32, error) {
	err := ValidateUserId(userID)
	if err != nil {
		return nil, 0, err
	}

	cs.mu.RLock()
	defer cs.mu.RUnlock()

	cart, err := cs.CartStorage.LoadCart(userID)
	if err != nil {
		if errors.Is(err, ErrCartNotFound) {
			return nil, 0, ErrCartNotFound
		}

		return nil, 0, fmt.Errorf("failed to get cart: %w", err)
	}

	if len(cart.Items) == 0 {
		return nil, 0, ErrCartNotFound
	}

	totalPrice := uint32(0)
	for _, item := range cart.Items {
		totalPrice += item.Price * item.Count
	}

	items := SortItems(&cart)

	return items, totalPrice, nil
}

func SortItems(c *Cart) []CartItem {
	items := make([]CartItem, 0, len(c.Items))

	skus := make([]int64, 0, len(c.Items))
	for sku := range c.Items {
		skus = append(skus, sku)
	}

	slices.Sort(skus)

	for _, sku := range skus {
		items = append(items, c.Items[sku])
	}

	return items
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
