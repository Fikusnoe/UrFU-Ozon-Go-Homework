package test

type testCartItem struct {
	Sku   int64  `json:"sku"`
	Price uint32 `json:"price"`
	Name  string `json:"name"`
}

var testItems = map[int64]testCartItem{
	1076963: testCartItem{
		Sku:   1076963,
		Name:  "Теория нравственных чувств | Смит Адам",
		Price: 3379,
	},
	1148162: testCartItem{
		Sku:   1148162,
		Name:  "Кулинар Гуров",
		Price: 2931,
	},
	135147479: testCartItem{
		Sku:   135147479,
		Name:  "Онлайн.ru. Сидишь в Интернете? Прочти! Тут о тебе! | Игошев Андрей",
		Price: 2981,
	},
}