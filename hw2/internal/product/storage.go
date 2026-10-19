package product

type Storage interface {
	LoadProducts() (map[int64]Product, error)
}