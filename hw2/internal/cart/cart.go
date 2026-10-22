package cart

import (
	"context"
	"fmt"

	pb "hw2/pkg/api/product"
)

type UserID int64
type User struct {
	UserID UserID
}

type CartItem struct {
	Sku   int64
	Count uint32
	Price uint32
	Name  string
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

func (cs *CartService) GetCart(userID UserID) (Cart, error) {
	cart, err := cs.CartStorage.LoadCart(userID)
	if err != nil {
		return Cart{}, err
	}

	return cart, nil
}

func ValidateInput(userID UserID, sku int64, count uint32) error {
	if userID <= 0 {
		return fmt.Errorf("Идентификатор пользователя должен быть натуральным числом (больше нуля)")
	} else if sku <= 0 {
		return fmt.Errorf("SKU должен быть натуральным числом (больше нуля)")
	} else if count <= 0 {
		return fmt.Errorf("Количество должно быть натуральным числом (больше нуля)")
	}

	return nil
}