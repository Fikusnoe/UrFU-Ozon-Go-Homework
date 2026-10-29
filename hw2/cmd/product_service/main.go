package main

import (
	"fmt"
	"hw2/internal/product"
	"hw2/internal/storage"
	"log"
	"net"

	pb "hw2/pkg/api/product"

	"google.golang.org/grpc"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("Service Product is down: %v", err)
	}
}

func run() error {
	productStorage, err := storage.NewJsonProductStorage("products.json")
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	productService, err := product.NewProductService(productStorage)
	if err != nil {
		return fmt.Errorf("failed to create product service: %w", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("failed to listen on port 50051: %w", err)
	}

	server := grpc.NewServer()
	pb.RegisterProductServiceServer(server, productService)
	log.Println("Service Product is up")

	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}
