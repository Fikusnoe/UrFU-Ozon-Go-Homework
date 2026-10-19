.PHONY: run-all:

.PHONY: generate:
	protoc --go_out=. \
	--go-grpc_out=. \
	--go_opt=module=github.com/Fikusnoe/hw2 \
	--go-grpc_opt=module=github.com/Fikusnoe/hw2 \
	.\proto\product_service.proto