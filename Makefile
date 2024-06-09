BINARY_NAME=springboot-app-monitor

build:
	@if [ -d bin ]; then rm -rf bin; fi
	@go clean
	@CGO_ENABLED=0 go build -o bin/${BINARY_NAME} ./cmd/app_monitor

run:
	@go run cmd/app_monitor/main.go

dist:
	@if [ -d dist/${BINARY_NAME} ]; then rm -rf dist/dist/${BINARY_NAME}; fi
	@go clean
	@CGO_ENABLED=0 go build -o dist/${BINARY_NAME} ./cmd/app_monitor
	@cp ./assets/logo.png ./dist/assets/
	@mkdir -p ./dist/config
	@touch ./dist/config/.env
	@echo "MONITOR_INTERVAL_MINUTES=5\n\
	SPRINGBOOT_APPLICATION_BASE_URLS=\n\n\
	CPU_USAGE_WARN_THRESHOLD=90\n\
	JVM_MEMORY_USAGE_WARN_THRESHOLD=80\n\n\
	ENABLE_DESKTOP_ALERTS=true\n\n\
	ENABLE_EMAIL_ALERTS=false\n\
	EMAIL_ALERT_RECIPIENTS=\n\n\
	SMTP_HOST=smtp.gmail.com\n\
	SMTP_PORT=587\n\
	SMTP_USER=\n\
	SMTP_PASSWORD=" > ./dist/config/.env