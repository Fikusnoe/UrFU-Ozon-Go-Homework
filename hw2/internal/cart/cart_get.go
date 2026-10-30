package cart

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"
)

func (cs *CartService) GetCart(ctx context.Context, userID UserID) ([]CartItem, uint32, error) {
	err := ValidateUserId(userID)
	if err != nil {
		return nil, 0, err
	}

	cartSnapshot := []CartItem{}
	cs.mu.RLock()

	cart, err := cs.CartStorage.LoadCart(userID)
	if err != nil {
		cs.mu.RUnlock()
		if errors.Is(err, ErrCartNotFound) {
			return nil, 0, ErrCartNotFound
		}

		return nil, 0, fmt.Errorf("failed to get cart: %w", err)
	}

	for sku, item := range cart.Items {
		cartSnapshot = append(cartSnapshot, CartItem{Sku: sku, Count: item.Count})
	}

	cs.mu.RUnlock()

	cartLen := len(cartSnapshot)
	if cartLen == 0 {
		return nil, 0, ErrCartNotFound
	}

	requestCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rpcResults := make([]rpcResult, cartLen)

	var wg sync.WaitGroup
	for i, item := range cartSnapshot {
		wg.Add(1)

		go func(index int, checkItem CartItem) {
			defer wg.Done()
			response, err := cs.MakeProductRPCRequest(requestCtx, checkItem.Sku)
			rpcResults[index] = rpcResult{data: response, err: err}
		}(i, item)
	}
	wg.Wait()

	for i, result := range rpcResults {
		if result.err != nil {
			return nil, 0, result.err
		}
		item := cartSnapshot[i]
		item.Name = result.data.Name
		item.Price = uint32(result.data.Price)
		cartSnapshot[i] = item
	}

	totalPrice := uint32(0)
	for _, item := range cartSnapshot {
		totalPrice += item.Price * item.Count
	}

	sortedItems := SortItems(cartSnapshot)

	return sortedItems, totalPrice, nil
}

func SortItems(items []CartItem) []CartItem {
	slices.SortFunc(items, func(a, b CartItem) int {
		if a.Sku > b.Sku {
			return 1
		} else if a.Sku < b.Sku {
			return -1
		}
		return 0
	})

	return items
}
