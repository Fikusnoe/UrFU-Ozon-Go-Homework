package product

import(
	"hw2/internal/domain"
)

type Storage interface {
	LoadProducts() (map[int64]domain.Product, error)
}