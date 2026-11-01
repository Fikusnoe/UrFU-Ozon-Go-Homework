package fake

import (
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
	if !ok  {
		return domain.Cart{}, cart.ErrCartNotFound
	}

	return userCart, nil
}
