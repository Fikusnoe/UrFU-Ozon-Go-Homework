package cart

type Storage interface {
	SaveCart(UserID, Cart)
	LoadCart(UserID) (Cart, error)
}