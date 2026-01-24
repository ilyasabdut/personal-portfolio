.PHONY: run build docker-up docker-down clean

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

clean:
	rm -rf bin/
