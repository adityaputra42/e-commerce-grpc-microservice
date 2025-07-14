
createdb:
	docker exec -it postgres16 createdb --username=root --owner=root ecommerce_microservice

dropdb:
	docker exec -it postgres16 dropdb ecommerce_microservice

run_server:

	cd api-gateaway && make gateaway_server &
	cd service/auth && make auth_server &
	cd service/user && make user_server &
	cd service/cars && make car_server &
	cd service/order && make order_server &
	cd service/payment && make payment_server &

run_services:
	cd service/auth && make auth_server &
	cd service/user && make user_server &
	cd service/cars && make car_server &
	cd service/order && make order_server &
	cd service/payment && make payment_server &

kill_ports:
	@for port in 8080 50051 50052 50053 50054 50055; do \
		pid=$$(lsof -t -i :$$port); \
		if [ -n "$$pid" ]; then \
			echo "Killing process on port $$port (PID: $$pid)"; \
			kill -9 $$pid; \
		else \
			echo "No process found on port $$port"; \
		fi \
	done


.PHONY: createdb dropdb run_server run_services kill_ports