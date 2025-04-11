run:
	go run cmd/main.go -ldflags="-X 'main.OsarkServerURL=http://localhost:8080'"
build:
	go build -o osark-daemon cmd/main.go
