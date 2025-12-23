# ========================
# Конфигурация проекта
# ========================
APP_NAME ?= teamboard
BUILD_DIR = bin
BINARY_PATH = $(BUILD_DIR)/$(APP_NAME)
GOBIN ?= $(shell go env GOPATH)/bin
LINTER_VERSION = v2.6.0

# Флаги для сборки
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_FLAGS = -ldflags="-s -w -X main.version=$(VERSION)"

# ========================
# Основные команды
# ========================
.PHONY: help
help:
	@echo "Доступные команды:"
	@echo "  make build         				 — собрать бинарник"
	@echo "  make build-flags    				 — собрать бинарник with flags"
	@echo "  make run           				 — запустить приложение"
	@echo "  make run-dev       				 — запустить приложение в режиме разработки"
	@echo "  make test           				 — запустить unit-тесты"
	@echo "  make deps          				 — обновить зависимости (go mod tidy && go mod vendor)"
	@echo "  make clean-build   				 — очистить билды"
	@echo "  make install-lint   				 — установить golangci-lint"
	@echo "  make lint           				 — запустить линтер"
	@echo "  make lint-fix      				 — запустить линтер с автоисправлением"
	@echo "  make pre-commit    				 — подготовка к коммиту (deps + lint + test)"
	@echo "  make swag           				 — генерация OpenApi документации"
	@echo "  make install-swag           			 — скачивание swaggo"
	@echo "  make docker-run   				 — запуск в контейнера (c беком)"
	@echo "  make docker-dev-run  				 — запуск в контейнера локальной разработки (без бека)"
	@echo "  make migrate-create {имя файла}     		 — создание новой миграции (up && down) в папке migrate"
	@echo "  make generate-docs    			 — инициализация OpenApi документации"
	@echo "  make install-mockery   			 — установить mockery"

# Сборка бинарника
.PHONY: build
build:
	@echo "Сборка $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BINARY_PATH) ./cmd/$(APP_NAME)

# Сборка бинарника with flags
.PHONY: build-flags
build-flags:
	@echo "Сборка $(APP_NAME) с флагами: $(BUILD_FLAGS)..."
	@mkdir -p $(BUILD_DIR)
	go build $(BUILD_FLAGS) -o $(BINARY_PATH) ./cmd/$(APP_NAME)

# Запуск приложения (оба файла)
.PHONY: run
run:
	@echo "Запуск приложения..."
	go run ./cmd/$(APP_NAME)/init.go ./cmd/$(APP_NAME)/main.go

# Запуск приложения (в режиме разработки)
.PHONY: run-dev
run-dev:
	@echo "Запуск приложения..."
	go run ./cmd/$(APP_NAME)/init.go ./cmd/$(APP_NAME)/main.go --config config/dev.yaml

# Запуск тестов
.PHONY: test
test:
	@echo "Запуск тестов..."
	go test -v ./...

# Обновление зависимостей
.PHONY: deps
deps:
	@echo "Обновление зависимостей..."
	go mod tidy
	go mod vendor
	@echo "Зависимости обновлены"

# Очистка билдов
.PHONY: clean-build
clean-build:
	@echo "Очистка бинарников..."
	rm -rf $(BUILD_DIR)
	go clean
	@echo "Очистка завершена"

# Установка golangci-lint (если отсутствует)
.PHONY: install-lint
install-lint:
	@if [ ! -f "$(GOBIN)/golangci-lint" ]; then \
		echo "Устанавливаю golangci-lint $(LINTER_VERSION)..."; \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(LINTER_VERSION); \
	else \
		echo "golangci-lint уже установлен в $(GOBIN)/golangci-lint"; \
	fi

# Запуск линтера
.PHONY: lint
lint: lint
	@echo "Запуск линтера..."
	@$(GOBIN)/golangci-lint run ./...

# Запуск линтера с автоисправлением
.PHONY: lint-fix
lint-fix: lint-fix
	@echo "Запуск линтера с автоисправлением..."
	@$(GOBIN)/golangci-lint run --fix ./...

# Обмазаться всем для коммита
.PHONY: pre-commit
pre-commit: deps lint test
	@echo "Всё готово для коммита!"

# Создать миграции в migrate
# Документация: https://github.com/golang-migrate/migrate
# Пример взаимодействия с миграциями:
#	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5431/teamboard?sslmode=disable" up {num}
#	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5431/teamboard?sslmode=disable" down {num}
#	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5431/teamboard?sslmode=disable" force {num}
MIGR_NAME := $(word 2,$(MAKECMDGOALS))
.PHONY: migrate-create
migrate-create:
	@echo "Создание миграций..."
	@test -n "$(MIGR_NAME)" || { echo "Usage: make migrate-create <name>"; exit 1; }
	@mkdir -p migrations
	migrate create -ext sql -dir migrations $(MIGR_NAME)

# Запуск в контейнера (c беком)
.PHONY: docker-run
docker-run:
	@echo "Запуск в контейнера (c беком)..."
	docker-compose up -d

# Запуск в контейнера локальной разработки (без бека)
# Для работы с контейнером: go run ./cmd/teamboard --config config/dev.yaml || make run-dev
.PHONY: docker-dev-run
docker-dev-run:
	@echo "Запуск в контейнера локальной разработки (без бека)..."
	docker-compose -f docker-compose.dev.yaml up --build -d

# Генерация документации Swagger
.PHONY: swag
swag: install-swag
	@echo "Инициализация OpenApi документации..."
	swag init -g ./cmd/teamboard/init.go

# Скачивание Swagger
.PHONY: install-swag
install-swag:
	@echo "Скачивание swago..."
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: install-mockery
install-mockery:
	@echo "Скачивание mockery..."
	go install github.com/vektra/mockery/v2@latest
	@echo "\x1b[33musage: mockery --name Repo --dir src/usecase/gettask --output src/usecase/gettask/mocks \x1b[0m"