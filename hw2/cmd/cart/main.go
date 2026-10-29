package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"hw2/internal/cart"
	"hw2/internal/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "hw2/pkg/api/product"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("Service Cart is down: %v", err)
	}
}

func run() error {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("Не удалось подключиться к product service: %w", err)
	}
	defer conn.Close()

	productClient := pb.NewProductServiceClient(conn)
	cartStorage := storage.NewCartStorage()
	cartService := cart.NewCartService(cartStorage, productClient)

	httpHandler := cart.NewHTTPHandler(cartService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku}", httpHandler.AddItem) 
	mux.HandleFunc("GET /user/{user_id}/cart", httpHandler.GetCart)

	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Println("Service Cart is up")

	err = httpServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}