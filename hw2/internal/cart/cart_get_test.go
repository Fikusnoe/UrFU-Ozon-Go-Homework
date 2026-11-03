package cart_test

import (
	"context"
	"hw2/internal/cart"
	"hw2/internal/domain"
	"hw2/internal/fake"
	"testing"

	"github.com/stretchr/testify/require"
)

type getStep struct {
	userID         int64
	sku            int64
	count          int64
	name           string
	price          uint32
	wantSku        int64
	wantName       string
	wantTotalPrice uint32
	wantCount      uint32
}

func TestGetItem(t *testing.T) {
	t.Parallel()
	stProduct := fake.NewProductStorage()

	tests := []struct {
		name   string
		userID int64
		isErr  error
		asErr  **cart.ValidationError
	}{
		{
			name:   "negative user id",
			userID: -1,
			asErr:  new(*cart.ValidationError),
		},
		{
			name:   "zero user id",
			userID: 0,
			asErr:  new(*cart.ValidationError),
		},
		{
			name:   "cart not found",
			userID: 1,
			isErr:  cart.ErrCartNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stCart := fake.NewCartStorage()
			client := fake.NewProductServiceClient(*stProduct)

			svc := cart.NewCartService(stCart, client)

			_, _, err := svc.GetCart(context.Background(), domain.UserID(tt.userID))
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

func TestGetItemSuccess(t *testing.T) {
	t.Parallel()
	stProduct := fake.NewProductStorage()

	tests := []struct {
		name  string
		steps []getStep
	}{
		{
			name: "successOneUserOneItem",
			steps: []getStep{
				{
					userID:         1,
					sku:            1076963,
					count:          1,
					name:           "Теория нравственных чувств | Смит Адам",
					wantSku:        1076963,
					wantName:       "Теория нравственных чувств | Смит Адам",
					wantTotalPrice: 3379,
					wantCount:      1,
				},
			},
		},
		{
			name: "successOneUserSameItem",
			steps: []getStep{
				{
					userID:         1,
					sku:            1076963,
					count:          1,
					name:           "Теория нравственных чувств | Смит Адам",
					wantSku:        1076963,
					wantName:       "Теория нравственных чувств | Смит Адам",
					wantTotalPrice: 3379,
					wantCount:      1,
				},
				{
					userID:         1,
					sku:            1076963,
					count:          2,
					name:           "Теория нравственных чувств | Смит Адам",
					wantSku:        1076963,
					wantName:       "Теория нравственных чувств | Смит Адам",
					wantTotalPrice: 10137,
					wantCount:      3,
				},
			},
		},
		{
			name: "successSameUserAnotherSkuTotalPrice",
			steps: []getStep{
				{
					userID:         1,
					sku:            1076963,
					count:          1,
					name:           "Теория нравственных чувств | Смит Адам",
					wantSku:        1076963,
					wantName:       "Теория нравственных чувств | Смит Адам",
					wantTotalPrice: 3379,
					wantCount:      1,
				},
				{
					userID:         1,
					sku:            1148162,
					count:          1,
					name:           "Кулинар Гуров",
					wantSku:        1148162,
					wantName:       "Кулинар Гуров",
					wantTotalPrice: 6310,
					wantCount:      1,
				},
			},
		},
		{
			name: "successDiffUser",
			steps: []getStep{
				{
					userID:         1,
					sku:            1076963,
					count:          1,
					name:           "Теория нравственных чувств | Смит Адам",
					wantSku:        1076963,
					wantName:       "Теория нравственных чувств | Смит Адам",
					wantTotalPrice: 3379,
					wantCount:      1,
				},
				{
					userID:         2,
					sku:            1148162,
					count:          1,
					name:           "Кулинар Гуров",
					wantSku:        1148162,
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
				stCart.AddItemDirectly(domain.UserID(step.userID), step.sku, uint32(step.count), step.price, step.name)

				gotItems, gotTotalPrice, err := svc.GetCart(context.Background(), domain.UserID(step.userID))
				require.NoError(t, err)

				require.Equal(t, step.wantTotalPrice, gotTotalPrice)

				haveItem := false
				for _, item := range gotItems {
					if item.Sku == step.wantSku {
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

func TestGetItemOrder(t *testing.T) {
	t.Parallel()
	stProduct := fake.NewProductStorage()

	tests := []struct {
		name         string
		wantSkuOrder []uint32
		steps        []getStep
	}{
		{
			name:         "successSameUserSkuOrder",
			wantSkuOrder: []uint32{1076963, 1148162, 135147479},
			steps: []getStep{
				{
					userID:         1,
					sku:            1076963,
					count:          1,
					name:           "Теория нравственных чувств | Смит Адам",
					wantSku:        1076963,
					wantName:       "Теория нравственных чувств | Смит Адам",
					wantTotalPrice: 3379,
					wantCount:      1,
				},
				{
					userID:         1,
					sku:            1148162,
					count:          1,
					name:           "Кулинар Гуров",
					wantSku:        1148162,
					wantName:       "Кулинар Гуров",
					wantTotalPrice: 6310,
					wantCount:      1,
				},
				{
					userID:         1,
					sku:            135147479,
					count:          1,
					name:           "Онлайн.ru. Сидишь в Интернете? Прочти! Тут о тебе! | Игошев Андрей",
					wantSku:        135147479,
					wantName:       "Онлайн.ru. Сидишь в Интернете? Прочти! Тут о тебе! | Игошев Андрей",
					wantTotalPrice: 9291,
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
				stCart.AddItemDirectly(domain.UserID(step.userID), step.sku, uint32(step.count), step.price, step.name)
			}

			gotItems, _, err := svc.GetCart(context.Background(), domain.UserID(tt.steps[0].userID))
			require.NoError(t, err)

			gotItemOrder := []uint32{}
			for _, item := range gotItems {
				gotItemOrder = append(gotItemOrder, uint32(item.Sku))
			}

			require.Equal(t, tt.wantSkuOrder, gotItemOrder)
		})
	}
}
