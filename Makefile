include .env
export

service-run-http:
	go run ./cmd/timo --http

service-run-cli:
	go run ./cmd/timo --cli
