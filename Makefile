.PHONY: dist

run:
	@UID=$(shell id -u) go run cmd/app_monitor/main.go

dist:
	@if [ -d dist/bin ]; then rm -f dist/bin; fi
	@go clean
	@CGO_ENABLED=0 go build -o dist/bin ./cmd/app_monitor
	@chmod +x ./dist/bin
	@cp -r ./assets/* ./dist/assets/