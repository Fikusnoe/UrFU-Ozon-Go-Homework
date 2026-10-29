package cart

type Storage interface {
	SaveCart(UserID, Cart) error
	LoadCart(UserID) (Cart, error)
}