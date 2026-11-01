package product

import (
	"context"
	"hw2/internal/stub"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "hw2/pkg/api/product"
)

func TestGetProduct(t *testing.T) {
	t.Parallel()
	st := stub.NewStubProductStorage()

	tests := []struct {
		name    string
		request *pb.GetProductRequest
		code    codes.Code
	}{
		{
			name:    "negative sku",
			request: &pb.GetProductRequest{Sku: -1},
			code:    codes.InvalidArgument,
		},
		{
			name:    "zero sku",
			request: &pb.GetProductRequest{Sku: 0},
			code:    codes.InvalidArgument,
		},
		{
			name:    "not found",
			request: &pb.GetProductRequest{Sku: 311111111},
			code:    codes.NotFound,
		},
		{
			name:    "nil request",
			request: nil,
			code:    codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc, err := NewProductService(st)
			require.NoError(t, err)

			got, err := svc.GetProduct(context.Background(), tt.request)

			require.Error(t, err)
			require.Nil(t, got)
			require.Equal(t, tt.code, status.Code(err))
		})
	}
}

func TestGetProductSuccess(t *testing.T) {
	t.Parallel()
	st := stub.NewStubProductStorage()

	tests := []struct {
		name  string
		title string
		sku   int64
		price int64
	}{
		{
			name:  "success1",
			title: "Теория нравственных чувств | Смит Адам",
			sku:   1076963,
			price: 3379,
		},
		{
			name:  "success2",
			title: "Планетарный проект. От устойчивого развития к управляемой гармонии",
			sku:   135399359,
			price: 1571,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc, err := NewProductService(st)
			require.NoError(t, err)

			got, err := svc.GetProduct(context.Background(), &pb.GetProductRequest{Sku: tt.sku})

			require.NoError(t, err)
			require.Equal(t,
				&pb.GetProductResponse{
					Sku:   tt.sku,
					Name:  tt.title,
					Price: int32(tt.price),
				}, got)
		})
	}
}
