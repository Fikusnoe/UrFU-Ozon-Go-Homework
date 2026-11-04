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

func TestProductGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		sku  int64
		err  codes.Code
	}{
		{
			name: "invalid sku",
			sku:  341212341,
			err:  codes.NotFound,
		},
		{
			name: "invalid sku",
			sku:  -1,
			err:  codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		client := newProductServiceClient(t)
		reqCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err := client.GetProduct(reqCtx, &pb.GetProductRequest{Sku: tt.sku})
		require.Error(t, err)
		require.Equal(t, tt.err, status.Code(err))
	}
}

func TestProductGetNil(t *testing.T) {
	t.Parallel()

	client := newProductServiceClient(t)
	reqCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := client.GetProduct(reqCtx, nil)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestProductGetSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		sku  int64
		err  codes.Code
	}{
		{
			name: "get item",
			sku:  1076963,
		},
		{
			name: "get same item",
			sku:  1076963,
		},
		{
			name: "get another item",
			sku:  135147479,
		},
	}

	client := newProductServiceClient(t)
	reqCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, tt := range tests {
		got, err := client.GetProduct(reqCtx, &pb.GetProductRequest{Sku: tt.sku})
		require.NoError(t, err)

		require.Equal(t, testItems[tt.sku].Sku, got.Sku)
		require.Equal(t, testItems[tt.sku].Name, got.Name)
		require.Equal(t, testItems[tt.sku].Price, uint32(got.Price))
	}
}
