DB_URL=postgresql://root:root@localhost:5432/ecommerce_microservice?sslmode=disable

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

rm_pb:	
	rm -f internal/pb/*.go

proto:	
	protoc --go_out=. --go-grpc_out=. ./internal/pb/product.proto

user_server:
	go run cmd/main.go

sqlc:
	sqlc generate




.PHONY: migrateup migratedown migrateup1 migratedown1 new_migration user_server proto rm_pb sqlc