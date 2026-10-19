package product

import "hw2/internal/product"

type Storage interface {
	LoadProducts() ([]product.Product, error)
}