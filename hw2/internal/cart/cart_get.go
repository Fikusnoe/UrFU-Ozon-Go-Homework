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

	cs.mu.RLock()

	cart, err := cs.CartStorage.LoadCart(userID)
	if err != nil {
		cs.mu.RUnlock()
		if errors.Is(err, ErrCartNotFound) {
			return nil, 0, ErrCartNotFound
		}

		return nil, 0, fmt.Errorf("failed to get cart: %w", err)
	}
	cs.mu.RUnlock()

	cartLen := len(cart.Items)
	if cartLen == 0 {
		return nil, 0, ErrCartNotFound
	}

	cxtWithDeadline, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	RPCResults := make([]RPCResult, cartLen)

	var wg sync.WaitGroup
	index := 0
	for i, _ := range cart.Items {
		wg.Add(1)

		go func(index int) {
			response, err := cs.MakeProductRPCRequest(cxtWithDeadline, i)
			RPCResults[index] = RPCResult{data: response, err: err}
			defer wg.Done()
		}(index)
		index++
	}
	wg.Wait()

	for _, result := range RPCResults {
		if result.err != nil {
			return nil, 0, result.err
		}
	}

	totalPrice := uint32(0)
	for _, item := range cart.Items {
		totalPrice += item.Price * item.Count
	}

	items := SortItems(&cart)

	return items, totalPrice, nil
}

func SortItems(c *Cart) []CartItem {
	items := make([]CartItem, 0, len(c.Items))

	skus := make([]int64, 0, len(c.Items))
	for sku := range c.Items {
		skus = append(skus, sku)
	}

	slices.Sort(skus)

	for _, sku := range skus {
		items = append(items, c.Items[sku])
	}

	return items
}
