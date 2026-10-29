package storage

import (
	"hw2/internal/cart"
)

type CartStorage struct {
	carts map[cart.UserID]cart.Cart
}

func NewCartStorage() *CartStorage {
	cartStorage := CartStorage{carts: map[cart.UserID]cart.Cart{}}
	return &cartStorage
}

func (cs *CartStorage) SaveCart(userID cart.UserID, updatedCart cart.Cart) error {
	cs.carts[userID] = updatedCart
	return nil
}

func (cs *CartStorage) LoadCart(userID cart.UserID) (cart.Cart, error) {
	userCart, ok := cs.carts[userID]
	if !ok  {
		return cart.Cart{}, cart.ErrCartNotFound
	}

	return userCart, nil
}
