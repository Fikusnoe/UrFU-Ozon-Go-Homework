package cart

import (
	"context"
	"fmt"
	"slices"

	pb "hw2/pkg/api/product"
)

type UserID int64
type User struct {
	userID UserID
}

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
}

func NewCartService(cs Storage, client pb.ProductServiceClient) CartService {
	return CartService{CartStorage: cs, Client: client}
}

func (cs *CartService) AddItem(ctx context.Context, userID UserID, sku int64, count uint32) error {
	err := ValidateInput(userID, sku, count)
	if err != nil {
		return err
	}

	cart, err := cs.CartStorage.LoadCart(userID)
	if err != nil {
		return err
	}

	item, exist := cart.Items[sku]
	if exist {
		item.Count += count
		cart.Items[sku] = item
	} else {
		request := pb.GetProductRequest{Sku: sku}
		response, err := cs.Client.GetProduct(ctx, &request)
		if err != nil {
			return err
		}

		cart.Items[sku] = CartItem{Sku: response.Sku,
			Count: count,
			Price: uint32(response.Price),
			Name:  response.Name,
		}
	}

	cs.CartStorage.SaveCart(userID, cart)

	return nil
}

func (cs *CartService) GetCart(userID UserID) ([]CartItem, uint32, error) {
	cart, err := cs.CartStorage.LoadCart(userID)
	if err != nil {
		return make([]CartItem, 0), 0, err
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
	if userID <= 0 {
		return fmt.Errorf("Идентификатор пользователя должен быть натуральным числом (больше нуля)")
	} else if sku <= 0 {
		return fmt.Errorf("SKU должен быть натуральным числом (больше нуля)")
	} else if count <= uint32(0) {
		return fmt.Errorf("Количество должно быть натуральным числом (больше нуля)")
	}

	return nil
}
