DB_URL=postgresql://root:root@localhost:5432/ecommerce_microservice?sslmode=disable
PROTO_DIR=../../proto-shared
PROTO_FILE=$(PROTO_DIR)/user.proto
OUT_DIR=./internal/pb

migrateup:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir internal/db/migrations -seq $(name)

proto:	
	rm -f internal/pb/*.go
	protoc -I $(PROTO_DIR) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILE)

user_server:
	go run cmd/main.go



.PHONY: migrateup migratedown migrateup1 migratedown1 new_migration user_server proto