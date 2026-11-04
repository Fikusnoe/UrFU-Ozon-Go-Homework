//  //go:build api

package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "hw2/pkg/api/product"
)

func TestProductGetSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		sku int64
		err    int
	}{
		{
			name:   "get item",
			sku: 1076963,
			err:    http.StatusNotFound,
		},
		{
			name:   "get another item",
			sku: 135147479,
			err:    http.StatusBadRequest,
		},
	}

	client := newProductServiceClient(t)
	reqCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, tt := range tests {
		got, err := client.GetItem(ctx, &pb.GetProductRequest{tt.})
		require.NoError(t, err)

		require.Equal(t, "Book", got.Title)
		require.Equal(t, "Very Great Book", got.Description)
	}
}
