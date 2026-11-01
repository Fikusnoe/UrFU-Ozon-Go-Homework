package cart_test

import (
	"context"
	"hw2/internal/cart"
	"hw2/internal/domain"
	"hw2/internal/fake"
	"testing"

	"github.com/stretchr/testify/require"

	//pb "hw2/pkg/api/product"
)

func TestAddItem(t *testing.T) {
	t.Parallel()
	stProduct := fake.NewProductStorage()

	tests := []struct {
		name string
		userID int64
		sku int64
		count int64
		isErr   error
    	asErr   **cart.ValidationError
	}{
		{
			name: "negative user id",
			userID: -1,
			sku: 1076963,
			count: 1,
			asErr: new(*cart.ValidationError),
		},
		{
			name: "negative sku",
			userID: 1,
			sku: -1,
			count: 1,
			asErr: new(*cart.ValidationError),
		},
		{
			name: "zero count",
			userID: 1,
			sku: 1076963,
			count: 0,
			asErr: new(*cart.ValidationError),
		},
		{
			name: "not found",
			userID: 1,
			sku: 311111111,
			count: 1,
			isErr: cart.ErrProductNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stCart := fake.NewCartStorage()
			client := fake.NewProductServiceClient(*stProduct)

			svc := cart.NewCartService(stCart, client)

			err := svc.AddItem(context.Background(), domain.UserID(tt.userID), tt.sku, uint32(tt.count))

			require.Error(t, err)
			
			if tt.isErr != nil {
				require.ErrorIs(t, err, tt.isErr)
			}
			if tt.asErr != nil {
				require.ErrorAs(t, err, tt.asErr)
			}
		})
	}
}

//func TestGetProductSuccess(t *testing.T) {
//	t.Parallel()
//	st := stub.NewStubProductStorage()

//	tests := []struct {
//		name  string
//		title string
//		sku   int64
//		price int64
//	}{
//		{
//			name:  "success1",
//			title: "Теория нравственных чувств | Смит Адам",
//			sku:   1076963,
//			price: 3379,
//		},
//		{
//			name:  "success2",
//			title: "Планетарный проект. От устойчивого развития к управляемой гармонии",
//			sku:   135399359,
//			price: 1571,
//		},
//	}

//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			t.Parallel()

//			svc, err := NewProductService(st)
//			require.NoError(t, err)

//			got, err := svc.GetProduct(context.Background(), &pb.GetProductRequest{Sku: tt.sku})

//			require.NoError(t, err)
//			require.Equal(t,
//				&pb.GetProductResponse{
//					Sku:   tt.sku,
//					Name:  tt.title,
//					Price: int32(tt.price),
//				}, got)
//		})
//	}
//}
