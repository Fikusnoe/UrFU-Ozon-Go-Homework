package cart

import (
	"context"
	"errors"
	"fmt"
	"hw2/internal/domain"
)

func (cs *CartService) AddItem(ctx context.Context, userID domain.UserID, sku int64, count uint32) error {
	err := ValidateInput(userID, sku, count)
	if err != nil {
		return err
	}

	response, err := cs.MakeProductRPCRequest(ctx, sku)
	if err != nil{
		return err
	}

	cs.mu.Lock()
	defer cs.mu.Unlock()

	userCart, err := cs.CartStorage.LoadCart(userID)
	if err != nil {
		if errors.Is(err, ErrCartNotFound) {
			userCart = domain.Cart{Items: map[int64]domain.CartItem{}}
		} else {
			return err
		}
	}

	item, exist := userCart.Items[sku]
	if exist {
		item.Count += count
		userCart.Items[sku] = item
	} else {
		userCart.Items[sku] = domain.CartItem{Sku: response.Sku,
			Count: count,
			Price: uint32(response.Price),
			Name:  response.Name,
		}
	}

	err = cs.CartStorage.SaveCart(userID, userCart)
	if err != nil {
		return fmt.Errorf("failed to save cart: %w", err)
	}

	return nil
}