.PHONY: run-all generate test test-race test-api

# Для Linux, использовал wsl
run-all:
	go run ./cmd/product_service/main.go & \
	go run ./cmd/cart/main.go

generate:
	protoc --go_out=. \
	--go-grpc_out=. \
	--go_opt=module=github.com/Fikusnoe/hw2 \
	--go-grpc_opt=module=github.com/Fikusnoe/hw2 \
	.\proto\product_service.proto

test:
	go test ./...

test-race:
	go test -race ./...

test-api:
	go test -v -count=1 -tags=api ./test/