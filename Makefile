# Переменные
APP_NAME = walletapp
DOCKER_IMAGE = $(APP_NAME):latest
DOCKER_COMPOSE = docker-compose

# Команды
.PHONY: build test run docker-build docker-run docker-down clean

# Сборка приложения
build:
	@echo "Building the application..."
	go build -o $(APP_NAME) ./cmd/main.go

# Тестирование приложения
test:
	@echo "Running tests..."
	go test ./...

# Запуск приложения
run: build
	@echo "Running the application..."
	./$(APP_NAME)

# Сборка Docker образа
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# Запуск приложения в Docker
docker-run: docker-build
	@echo "Starting application in Docker..."
	$(DOCKER_COMPOSE) up --build

# Остановка и удаление контейнеров Docker
docker-down:
	@echo "Stopping and removing Docker containers..."
	$(DOCKER_COMPOSE) down

# Очистка скомпилированных файлов и образов
clean:
	@echo "Cleaning up..."
	rm -f $(APP_NAME)
	docker rmi $(DOCKER_IMAGE) || true