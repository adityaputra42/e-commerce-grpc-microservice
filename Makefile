
createdb:
	docker exec -it postgres16 createdb --username=root --owner=root ecommerce_microservice

dropdb:
	docker exec -it postgres16 dropdb ecommerce_microservice

.PHONY: createdb dropdb 