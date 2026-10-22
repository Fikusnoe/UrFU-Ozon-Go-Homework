package cart

import (
	"context"
	pb "hw2/pkg/api/product"
)

type UserID int64
type User struct {
	UserID UserID
}

type CartItem struct {
	Sku     int64
	Count   uint32
	Price uint32
	Name string
}

type Cart struct {
	Items map[int64]CartItem
}

type CartService struct {
	CartStorage Storage
	Client pb.ProductServiceClient
}

func NewCartService(cs Storage, client pb.ProductServiceClient) CartService {
	return CartService{ CartStorage: cs, Client: client}
}

func (cs CartService) AddItem(ctx context.Context, userID UserID, sku int64) error {
	cart, err := cs.CartStorage.LoadCart(userID)
	if err != nil{
		return err
	}

	item, exist := cart.Items[sku]
	if exist {
		item.Count++
		cart.Items[sku] = item
	} else {
		request := pb.GetProductRequest{Sku: sku}
		response, err := cs.Client.GetProduct(ctx, &request)
		if err != nil {
			return err
		}

		cart.Items[sku] = CartItem{ Sku: response.Sku, 
			Count: 1, 
			Price: uint32(response.Price), 
			Name: response.Name,
		}
	}

	return nil
}

func (c CartService) GetCart(userID UserID) error {
	cart, err := c.CartStorage.LoadCart(userID)
	if err != nil {

	}

	

	return nil
}

