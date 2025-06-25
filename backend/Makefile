BINARY_NAME=server
MAIN_PATH=./cmd/server
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

.PHONY: help build test clean run deps lint coverage benchmark docker-build docker-run

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

deps: ## Установить зависимости
	@echo "Установка зависимостей..."
	go mod download
	go mod tidy

build: ## Собрать приложение
	@echo "Сборка приложения..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

run: build ## Запустить приложение
	@echo "Запуск приложения..."
	./bin/$(BINARY_NAME)

test: ## Запустить тесты
	@echo "Запуск тестов..."
	go test -v ./...

test-race: ## Запустить тесты с проверкой гонок
	@echo "Запуск тестов с проверкой race conditions..."
	go test -race -v ./...

coverage: ## Генерация отчета о покрытии кода
	@echo "Генерация отчета о покрытии..."
	go test -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Отчет сохранен в $(COVERAGE_HTML)"

coverage-text: ## Показать покрытие в консоли
	@echo "Покрытие кода:"
	go test -cover ./...

benchmark: ## Запустить бенчмарки
	@echo "Запуск бенчмарков..."
	go test -bench=. -benchmem ./...

lint: ## Проверка кода линтером
	@echo "Проверка кода линтером..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout 5m; \
	else \
		echo "golangci-lint не найден. Используем стандартные инструменты Go..."; \
		go fmt ./...; \
		go vet ./...; \
		echo "Для полной проверки установите golangci-lint: make install-golangci-lint"; \
	fi

install-golangci-lint: ## Установить golangci-lint
	@echo "Установка golangci-lint..."
	@if [ "$$(uname)" = "Linux" ] || [ "$$(uname)" = "Darwin" ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2; \
	elif [ "$$(uname -o 2>/dev/null)" = "Msys" ] || [ "$$(uname -s)" = "MINGW64_NT-10.0" ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2; \
	else \
		echo "Скачайте golangci-lint с https://github.com/golangci/golangci-lint/releases"; \
	fi

fmt: ## Форматирование кода
	@echo "Форматирование кода..."
	go fmt ./...

vet: ## Статический анализ кода
	@echo "Статический анализ..."
	go vet ./...

clean: ## Очистка сгенерированных файлов
	@echo "Очистка..."
	rm -f bin/$(BINARY_NAME)
	rm -f $(COVERAGE_FILE)
	rm -f $(COVERAGE_HTML)
	rm -rf bin/

docker-build: ## Собрать Docker образ
	@echo "Сборка Docker образа..."
	docker build -t codegrader-backend .

docker-run: docker-build ## Запустить в Docker
	@echo "Запуск в Docker..."
	docker run -p 8080:8080 codegrader-backend

quality-check: fmt vet lint test coverage ## Полная проверка качества кода
	@echo "Проверка качества кода завершена!"

prepare-prod: clean deps quality-check build ## Подготовка к развертыванию
	@echo "Подготовка к продакшену завершена!"
