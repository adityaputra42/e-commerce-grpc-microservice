PROTO_DIR=../proto-shared
OUT_DIR=./internal/pb

proto:	
	rm -f internal/pb/*.go
	protoc -I $(PROTO_DIR) -I $(PROTO_DIR)/google \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(OUT_DIR) --grpc-gateway_opt=paths=source_relative \
		$(PROTO_DIR)/auth.proto \
		$(PROTO_DIR)/user.proto \
		$(PROTO_DIR)/cars.proto \
		$(PROTO_DIR)/order.proto \
		$(PROTO_DIR)/payment.proto \
		$(PROTO_DIR)/payment_wallet.proto 

gateaway_server:
	go run cmd/main.go

.PHONY: proto gateaway_server