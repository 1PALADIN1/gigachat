CHAT_APP_BINARY=chat-app
AUTH_SERVICE_BINARY=auth-app

## up: запускает все контейнеры
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up-build: останавливает docker compose (если запущен), собирает проекты и запускает docker compose
up-build: build-chat-app build-auth
	@echo "Stopping docker images (if running)..."
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: останавливает все контейнеры
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build-chat-app: собирает чат микросервис под Linux
build-chat-app:
	@echo "Building chat-app binary..."
	cd ../chat-app && env GOOS=linux CGO_ENABLED=0 go build -o ./bin/${CHAT_APP_BINARY} ./cmd/gigachat/*.go
	@echo "Done!"

## build-auth: собирает сервис авторизации под Linux
build-auth:
	@echo "Building auth binary..."
	cd ../auth-service && env GOOS=linux CGO_ENABLED=0 go build -o ./bin/${AUTH_SERVICE_BINARY} ./cmd/api/*.go
	@echo "Done!"
