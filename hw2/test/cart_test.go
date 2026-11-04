package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hw2/internal/cart"
	"hw2/internal/domain"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserIDGenerator struct {
	mu      *sync.RWMutex
	counter domain.UserID
}

func (g *UserIDGenerator) Next() domain.UserID {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.counter++

	return g.counter
}

func TestCartService(t *testing.T) {
	t.Parallel()

	userID := UserIDGenerator{
		mu:      &sync.RWMutex{},
		counter: 1,
	}

	t.Run("fail - get empty cart via HTTP", func(t *testing.T) {
		t.Parallel()

		resp, err := http.Get(cartServicePort + fmt.Sprintf("/user/%v/cart", userID.Next()))
		defer resp.Body.Close()
		require.Error(t, err)
	})

	t.Run("success - get cart via HTTP", func(t *testing.T) {
		t.Parallel()

		userID := userID.Next()

		wantCount := 1
		bodyData := cart.AddRequestBody{Count: int32(wantCount)}

		jsonData, err := json.Marshal(bodyData)
		require.NoError(t, err)

		body := bytes.NewReader(jsonData)

		http.Post(fmt.Sprintf("/user/%v/cart/%v", userID, 1076963), "application/json", body)

		gotResp, err := http.Get(cartServicePort + fmt.Sprintf("/user/%v/cart", userID))
		require.NoError(t, err)
		defer gotResp.Body.Close()

		require.Equal(t, http.StatusOK, gotResp.StatusCode)

		items := []domain.CartItem{}
		err = json.NewDecoder(gotResp.Body).Decode(&items)
		require.NoError(t, err)

		require.Equal(t, items[0].Sku, 1076963)
		require.Equal(t, items[0].Count, wantCount)
		require.Equal(t, items[0].Name, "Теория нравственных чувств | Смит Адам")
		require.Equal(t, items[0].Price, 3379)
	})
}
