package fake

import (
	"errors"
	"hw2/internal/cart"
	"hw2/internal/domain"
)

type CartStorage struct {
	carts map[domain.UserID]domain.Cart
}

func NewCartStorage() *CartStorage {
	cartStorage := CartStorage{carts: map[domain.UserID]domain.Cart{}}
	return &cartStorage
}

func (cs *CartStorage) SaveCart(userID domain.UserID, updatedCart domain.Cart) error {
	cs.carts[userID] = updatedCart
	return nil
}

func (cs *CartStorage) LoadCart(userID domain.UserID) (domain.Cart, error) {
	userCart, ok := cs.carts[userID]
	if !ok {
		return domain.Cart{}, cart.ErrCartNotFound
	}

	return userCart, nil
}

func (cs *CartStorage) AddItemDirectly(userID domain.UserID, sku int64, count uint32, price uint32, name string) error {
	userCart, err := cs.LoadCart(userID)
	if err != nil {
		if errors.Is(err, cart.ErrCartNotFound) {
			userCart = domain.Cart{Items: map[int64]domain.CartItem{}}
		} else {
			return err
		}
	}

	item, exist := userCart.Items[sku]
	if exist {
		item.Count += count
		userCart.Items[sku] = item
	} else {
		userCart.Items[sku] = domain.CartItem{
			Sku:   sku,
			Count: count,
			Price: uint32(price),
			Name:  name,
		}
	}

	_ = cs.SaveCart(userID, userCart)

	return nil
}
