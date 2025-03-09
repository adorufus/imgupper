.PHONY: build run test clean migration-up migration-down docker-build docker-run

# Build the application
build:
	go build -o bin/apiserver main.go

# Run the application
run: build
	./bin/apiserver

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Database migrations
migration-up:
	migrate -path migrations -database "postgres://HidungBelalang:Kuatcore141!@localhost:5432/imgupper?sslmode=disable" up

migration-down:
	migrate -path migrations -database "postgres://HidungBelalang:Kuatcore141!@localhost:5432/imgupper?sslmode=disable" down

# Docker commands
docker-build:
	docker build -t myapi:latest .

docker-run:
	docker run -p 8080:8080 myapi:latest