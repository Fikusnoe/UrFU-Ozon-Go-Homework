package main

import (
	"fmt"

	"hw2/internal/storage"
	"hw2/internal/cart"
	"google.golang.org/grpc"

	pb "hw2/pkg/api/product"
)

func main() {

}

func run() error {
	conn, err := grpc.NewClient("localhost:50051")
	if err != nil {
		return fmt.Errorf("Не удалось подключиться к product service: %w", err)
	}
	defer conn.Close()

	productClient := pb.NewProductServiceClient(conn)
	cartStorage := storage.NewCartStorage()
	cartService := cart.NewCartService(cartStorage, productClient)


	return nil
}