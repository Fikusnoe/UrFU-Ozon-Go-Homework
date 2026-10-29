package main

import (
	"context"
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
		log.Fatalf("Service Product is down: %v", err)
	}
}

func run() error {
	ctx, stop :=  context.WithCancel(context.Background())
	defer stop()

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
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	fmt.Println("Сервис cart запущен")

	errCh := make(chan error, 1)
	go func() {
		errCh <- httpServer.ListenAndServe()
	}()

	select {
		case err := <-errCh:
		if err != http.ErrServerClosed {
			return err
		}
		return nil
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err := httpServer.Shutdown(shutdownCtx)
			if  err != nil {
				return err
			}
		}

	return nil
}