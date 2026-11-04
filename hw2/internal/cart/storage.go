package cart

import "hw2/internal/domain"

type Storage interface {
	SaveCart(domain.UserID, domain.Cart) error
	LoadCart(domain.UserID) (domain.Cart, error)
}