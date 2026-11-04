// //go:build api

package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hw2/internal/domain"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type addRequestBody struct {
	Count int32 `json:"count"`
}

type cartResponse struct {
	Items      []domain.CartItem `json:"items"`
	TotalPrice uint32            `json:"total_price"`
}

func GenerateUserID() domain.UserID {
	return domain.UserID(time.Now().UnixNano())
}

type step struct {
	userID    domain.UserID
	sku       int64
	count     uint32
	wantCount uint32
}

func TestCartGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		userID domain.UserID
		err    int
	}{
		{
			name:   "cartDoesntExist",
			userID: GenerateUserID(),
			err:    http.StatusNotFound,
		},
		{
			name:   "invalidUserID",
			userID: -1,
			err:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resp, err := http.Get(cartServicePort + fmt.Sprintf("/user/%v/cart", tt.userID))
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tt.err, resp.StatusCode)
		})
	}
}

func TestGetCartSuccess(t *testing.T) {
	t.Parallel()

	userID1 := GenerateUserID()

	tests := []struct {
		name  string
		steps []step
	}{
		{
			name: "success One User One Item",
			steps: []step{
				step{
					userID:    userID1,
					sku:       1076963,
					count:     1,
					wantCount: 1,
				},
			},
		},
		{
			name: "success One User Same Item",
			steps: []step{
				step{
					userID:    userID1,
					sku:       1076963,
					count:     1,
					wantCount: 2,
				},
			},
		},
		{
			name: "success Same User Another Sku",
			steps: []step{
				step{
					userID:    userID1,
					sku:       1148162,
					count:     1,
					wantCount: 1,
				},
			},
		},
		{
			name: "success Diff User",
			steps: []step{
				step{
					userID:    GenerateUserID(),
					sku:       1148162,
					count:     1,
					wantCount: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for _, step := range tt.steps {
				wantTotalPrice := uint32(0)
				bodyData := addRequestBody{Count: int32(step.count)}

				jsonData, err := json.Marshal(bodyData)
				require.NoError(t, err)

				body := bytes.NewReader(jsonData)

				postResp, err := http.Post(cartServicePort+fmt.Sprintf("/user/%v/cart/%v", step.userID, step.sku), "application/json", body)
				require.NoError(t, err)
				defer postResp.Body.Close()
				require.Equal(t, http.StatusOK, postResp.StatusCode)

				gotResp, err := http.Get(cartServicePort + fmt.Sprintf("/user/%v/cart", step.userID))
				require.NoError(t, err)
				defer gotResp.Body.Close()

				require.Equal(t, http.StatusOK, gotResp.StatusCode)

				resp := cartResponse{}
				err = json.NewDecoder(gotResp.Body).Decode(&resp)
				require.NoError(t, err)

				items := resp.Items

				haveItem := false
				for _, item := range items {
					wantTotalPrice += item.Price * item.Count
					if item.Sku == step.sku {
						require.Equal(t, step.wantCount, item.Count)
						require.Equal(t, testItems[step.sku].Name, item.Name)
						require.Equal(t, testItems[step.sku].Price, item.Price)
						haveItem = true
					}
				}

				require.Equal(t, wantTotalPrice, resp.TotalPrice)
				if !haveItem {
					require.Fail(t, "Item not added")
				}
			}
		})
	}
}

func TestCartAdd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		userID domain.UserID
		sku    int64
		count  int64
		err    int
	}{
		{
			name:   "invalid Sku",
			userID: GenerateUserID(),
			sku:    -1,
			count:  1,
			err:    http.StatusBadRequest,
		},
		{
			name:   "wrong Sku",
			userID: GenerateUserID(),
			sku:    1452541123,
			count:  1,
			err:    http.StatusPreconditionFailed,
		},
		{
			name:   "invalid Count",
			userID: GenerateUserID(),
			sku:    1076963,
			count:  -1,
			err:    http.StatusBadRequest,
		},
		{
			name:   "invalid User ID",
			userID: -1,
			sku:    1076963,
			count:  1,
			err:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bodyData := addRequestBody{Count: int32(tt.count)}

			jsonData, err := json.Marshal(bodyData)
			require.NoError(t, err)

			body := bytes.NewReader(jsonData)

			postResp, err := http.Post(cartServicePort+fmt.Sprintf("/user/%v/cart/%v", tt.userID, tt.sku), "application/json", body)
			require.NoError(t, err)
			require.Equal(t, tt.err, postResp.StatusCode)
			defer postResp.Body.Close()
		})
	}
}
