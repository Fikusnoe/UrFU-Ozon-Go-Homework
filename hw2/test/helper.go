//go:build api

package test

import (
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "hw2/pkg/api/product"
)

const productServicePort = "localhost:50051"
const cartServicePort = "http://localhost:8080"

func newProductServiceClient(t *testing.T) pb.ProductServiceClient {
	t.Helper()

	conn, err := grpc.NewClient(productServicePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to create grpc client: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	return pb.NewProductServiceClient(conn)
}
