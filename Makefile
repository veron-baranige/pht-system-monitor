BINARY_NAME=springboot-app-monitor

build:
	@if [ -d bin ]; then rm -rf bin; fi
	@go clean
	@CGO_ENABLED=0 go build -o bin/${BINARY_NAME} ./cmd/app_monitor

run:
	@go run cmd/app_monitor/main.go