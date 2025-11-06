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
	@echo "  make build          — собрать бинарник"
	@echo "  make build-flags    — собрать бинарник with flags"
	@echo "  make run            — запустить приложение"
	@echo "  make test           — запустить unit-тесты"
	@echo "  make deps           — обновить зависимости (go mod tidy && go mod vendor)"
	@echo "  make clean-build    — очистить билды"
	@echo "  make install-lint   — установить golangci-lint"
	@echo "  make lint           — запустить линтер"
	@echo "  make lint-fix       — запустить линтер с автоисправлением"
	@echo "  make pre-commit     — подготовка к коммиту (deps + lint + test)"

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