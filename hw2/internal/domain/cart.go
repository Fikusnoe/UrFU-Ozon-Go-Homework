package domain

type UserID int64

type CartItem struct {
	Sku   int64  `json:"sku"`
	Count uint32 `json:"count"`
	Price uint32 `json:"price"`
	Name  string `json:"name"`
}

type Cart struct {
	Items map[int64]CartItem
}