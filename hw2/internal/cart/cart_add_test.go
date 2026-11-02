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

type addStep struct {
	userID         int64
	sku            int64
	count          int64
	wantName       string
	wantTotalPrice uint32
	wantCount      uint32
}

func TestAddItem(t *testing.T) {
	t.Parallel()
	stProduct := fake.NewProductStorage()

	tests := []struct {
		name   string
		userID int64
		sku    int64
		count  int64
		isErr  error
		asErr  **cart.ValidationError
	}{
		{
			name:   "negative user id",
			userID: -1,
			sku:    1076963,
			count:  1,
			asErr:  new(*cart.ValidationError),
		},
		{
			name:   "negative sku",
			userID: 1,
			sku:    -1,
			count:  1,
			asErr:  new(*cart.ValidationError),
		},
		{
			name:   "zero count",
			userID: 1,
			sku:    1076963,
			count:  0,
			asErr:  new(*cart.ValidationError),
		},
		{
			name:   "not found",
			userID: 1,
			sku:    311111111,
			count:  1,
			isErr:  cart.ErrProductNotFound,
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

func TestAddItemSuccess(t *testing.T) {
	t.Parallel()
	stProduct := fake.NewProductStorage()

	tests := []struct {
		name  string
		steps []addStep
	}{
		{
			name: "successOneUserOneItem",
			steps: []addStep{
				{
					userID:         1,
					sku:            1076963,
					count:          1,
					wantName:       "Теория нравственных чувств | Смит Адам",
					wantTotalPrice: 3379,
					wantCount:      1,
				},
			},
		},
		{
			name: "successOneUserSameItem",
			steps: []addStep{
				{
					userID:         1,
					sku:            1076963,
					count:          1,
					wantName:       "Теория нравственных чувств | Смит Адам",
					wantTotalPrice: 3379,
					wantCount:      1,
				},
				{
					userID:         1,
					sku:            1076963,
					count:          2,
					wantName:       "Теория нравственных чувств | Смит Адам",
					wantTotalPrice: 10137,
					wantCount:      3,
				},
			},
		},
		{
			name: "successSameUserAnotherSku",
			steps: []addStep{
				{
					userID:         1,
					sku:            1076963,
					count:          1,
					wantName:       "Теория нравственных чувств | Смит Адам",
					wantTotalPrice: 3379,
					wantCount:      1,
				},
				{
					userID:         1,
					sku:            1148162,
					count:          1,
					wantName:       "Кулинар Гуров",
					wantTotalPrice: 6310,
					wantCount:      1,
				},
			},
		},
		{
			name: "successDiffUser",
			steps: []addStep{
				{
					userID:         1,
					sku:            1076963,
					count:          1,
					wantName:       "Теория нравственных чувств | Смит Адам",
					wantTotalPrice: 3379,
					wantCount:      1,
				},
				{
					userID:         2,
					sku:            1148162,
					count:          1,
					wantName:       "Кулинар Гуров",
					wantTotalPrice: 2931,
					wantCount:      1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			stCart := fake.NewCartStorage()
			client := fake.NewProductServiceClient(*stProduct)
			svc := cart.NewCartService(stCart, client)

			for _, step := range tt.steps {

				err := svc.AddItem(context.Background(), domain.UserID(step.userID), step.sku, uint32(step.count))
				require.NoError(t, err)

				gotItems, gotTotalPrice, err := svc.GetCart(context.Background(), domain.UserID(step.userID))
				require.NoError(t, err)

				require.Equal(t, step.wantTotalPrice, gotTotalPrice)

				haveItem := false
				for _, item := range gotItems{
					if item.Sku == step.sku{
						haveItem = true
						require.Equal(t, step.wantCount, item.Count)
						require.Equal(t, step.wantName, item.Name)
					}
				}

				if !haveItem {
					require.Fail(t, "Item not added")
				}
			}
		})
	}
}
