package main

import (
	"context"
	"log"
	"hw2/internal/product"
	"hw2/internal/storage"
	"net"

	pb "hw2/pkg/api/product"

	"google.golang.org/grpc"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("Сервис Product упал: %v", err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	productStorage, err := storage.NewJsonProductStorage("products.json")
	if err != nil {
		return err
	}

	productService, err := product.NewProductService(productStorage)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterProductServiceServer(server, productService)
	log.Println("Сервис Product запущен")

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Serve(lis)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		go func() {
			server.GracefulStop()
		}()
	}

	return nil
}
