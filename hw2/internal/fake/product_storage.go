package stub

import (
	"hw2/internal/domain"
)

type StubProductStorage struct {
}

func (s *StubProductStorage) LoadProducts() (map[int64]domain.Product, error) {
	dataMap := map[int64]domain.Product{
		1076963: {
			Name:  "Теория нравственных чувств | Смит Адам",
			Price: 3379,
		},
		1148162: {Name: "Кулинар Гуров", Price: 2931},
		135147479: {
			Name:  "Онлайн.ru. Сидишь в Интернете? Прочти! Тут о тебе! | Игошев Андрей",
			Price: 2981,
		},
		135399359: {
			Name:  "Планетарный проект. От устойчивого развития к управляемой гармонии",
			Price: 1571,
		},
		135466315: {Name: "Менялы | Хейли Артур", Price: 3322},
		135567670: {
			Name:  "Комплект постельного белья Verossa, Сатин, 1,5 спальный, наволочки 70x70",
			Price: 2458,
		},
		135585483: {
			Name:  "Ваш Андрей Петров. Композитор в воспоминаниях современников | Коллектив авторов",
			Price: 696,
		},
	}

	return dataMap, nil
}

func NewStubProductStorage() *StubProductStorage {
	storage := StubProductStorage{}
	return &storage
}
