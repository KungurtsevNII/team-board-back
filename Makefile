# Название бинарника линтера
LINTER = golangci-lint
# Версия линтера
LINTER_VERSION = v2.6.0
# Имя приложения
APP_NAME ?= teamboard
# Директория для билдов
BUILD_DIR = bin
# Директория, где хранится go bin (по умолчанию ~/go/bin)
GOBIN ?= $(shell go env GOPATH)/bin
# Полный путь до бинарника линтера
LINTER_PATH = $(GOBIN)/$(LINTER)

# Показать доступные команды
.PHONY: help
help:
	@echo "Доступные команды:"
	@echo "  make build          — собрать бинарник"
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
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

# Запуск приложения
.PHONY: run
run:
	@echo "Запуск приложения..."
	go run ./cmd/$(APP_NAME)

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
	@if [ ! -f "$(LINTER_PATH)" ]; then \
		echo "Устанавливаю golangci-lint $(LINTER_VERSION)..."; \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(LINTER_VERSION); \
	else \
		echo "golangci-lint уже установлен в $(LINTER_PATH)"; \
	fi

# Запуск линтера
.PHONY: lint
lint: install-lint
	@echo "Запуск линтера..."
	@$(LINTER) run ./...

# Запуск линтера с автоисправлением
.PHONY: lint-fix
lint-fix: install-lint
	@echo "Запуск линтера с автоисправлением..."
	@$(LINTER) run --fix ./...

# Обмазаться всем для коммита
.PHONY: pre-commit
pre-commit: deps lint test
	@echo "Всё готово для коммита!"