package cart

import (
	"hw2/internal/product"
)

type UserID int64
type Sku int64

type User struct {
	UserID UserID
}

type CartItem struct {
	Sku     Sku
	Count   int
	product product.Product
}

type Cart struct {
	Items map[Sku]CartItem
}

func (c Cart) AddItem(sku Sku) error {

}

func (c Cart) GetCart() error {

}
