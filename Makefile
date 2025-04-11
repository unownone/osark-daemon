run:
	go run cmd/main.go
build:
	go build -o osark-daemon cmd/main.go
format:
	go fmt ./...
lint:
	golangci-lint run ./...
test:
	go test -v ./...
