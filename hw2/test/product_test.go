////  //go:build api

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "hw2/pkg/api/product"
)

func TestProductServiceOverNetwork(t *testing.T) {
	t.Parallel()

	client := newProductServiceClient(t)
	ctx := context.Background()

	post, err := client.PostItem(ctx, &pb.PostItemRequest{Title: "Book", Description: "Very Great Book"})
	require.NoError(t, err)
	require.NotZero(t, post.Id)

	got, err := client.GetItem(ctx, &pb.GetItemRequest{Id: post.Id})
	require.NoError(t, err)
	require.Equal(t, "Book", got.Title)
	require.Equal(t, "Very Great Book", got.Description)

	_, err = client.GetItem(ctx, &pb.GetItemRequest{Id: 0})
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	_, err = client.GetItem(ctx, &pb.GetItemRequest{Id: 999999999})
	require.Equal(t, codes.NotFound, status.Code(err))
}
