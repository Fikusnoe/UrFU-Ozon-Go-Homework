.PHONY: run-all generate

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