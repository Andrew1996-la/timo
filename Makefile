include .env
export

service-run:
	go run ./cmd/timo

migrate-up:
	migrate -path migration -database ${CONN_STRING} up

migrate-down:
	migrate -path migration -database ${CONN_STRING} down
