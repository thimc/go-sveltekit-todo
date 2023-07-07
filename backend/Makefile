.PHONY: run test swagger

build:
	@go build -o bin/api

run: swagger
	@go run .

test:
	@go test -v ./...

swagger:
	swag init --dir ./,./api

docker:
	@docker build -t backend .
	@docker run -p 1111:1111 backend

