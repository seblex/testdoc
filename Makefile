.PHONY: build build-cli test test-race test-short test-cover lint fmt vet tidy clean install demo examples docs docker-check docker-build docker-run help

# Переменные
BINARY_NAME=testdoc
CLI_PATH=./cmd/testdoc
COVERAGE_FILE=coverage.out
DOCKER_IMAGE=seblex5/testdoc

# Docker команда (только docker, без sudo)
DOCKER_CMD := docker

# Сборка CLI приложения
build:
	@echo "🔨 Сборка CLI приложения..."
	@go build -o bin/$(BINARY_NAME) $(CLI_PATH)
	@echo "✅ Сборка завершена: bin/$(BINARY_NAME)"

# Сборка с версией
build-cli:
	@echo "🔨 Сборка CLI с версией..."
	@mkdir -p bin
	@go build -ldflags "-X main.version=$(shell git describe --tags --always --dirty)" -o bin/$(BINARY_NAME) $(CLI_PATH)
	@echo "✅ Сборка завершена: bin/$(BINARY_NAME)"

# Установка локально
install:
	@echo "📦 Установка $(BINARY_NAME)..."
	@go install $(CLI_PATH)
	@echo "✅ $(BINARY_NAME) установлен в $$(go env GOPATH)/bin/"

# Запуск всех тестов
test:
	@echo "🧪 Запуск всех тестов..."
	@go test ./pkg/... ./ 

# Запуск тестов с race detector
test-race:
	@echo "🏃 Запуск тестов с race detector..."
	@go test -race ./...

# Запуск только быстрых тестов
test-short:
	@echo "⚡ Запуск быстрых тестов..."
	@go test -short ./...

# Тесты с покрытием
test-cover:
	@echo "📊 Тесты с анализом покрытия..."
	@go test -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@go tool cover -func=$(COVERAGE_FILE)
	@echo "📈 HTML отчет: make cover-html"

# HTML отчет покрытия
cover-html: test-cover
	@go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "📊 Отчет покрытия: coverage.html"

# Линтинг
lint:
	@echo "🔍 Запуск линтера..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Установите golangci-lint: https://golangci-lint.run/usage/install/"; exit 1; }
	@golangci-lint run

# Форматирование кода
fmt:
	@echo "🎨 Форматирование кода..."
	@go fmt ./...

# Статический анализ
vet:
	@echo "🔍 Статический анализ..."
	@go vet ./...

# Обновление зависимостей
tidy:
	@echo "📦 Обновление go.mod..."
	@go mod tidy
	@go mod verify

# Полная проверка качества
check: fmt vet lint test-race
	@echo "✅ Все проверки пройдены"

# Очистка
clean:
	@echo "🧹 Очистка файлов..."
	@rm -rf bin/
	@rm -f $(COVERAGE_FILE) coverage.html
	@rm -f *-documentation.md
	@echo "✅ Очистка завершена"

# Демонстрация работы
demo: build
	@echo "🚀 Демонстрация TestDoc"
	@echo ""
	@echo "1. Анализ примеров тестов:"
	@./bin/$(BINARY_NAME) examples/_examples
	@echo ""
	@echo "2. Первые строки сгенерированной документации:"
	@head -20 test-documentation.md || echo "Документация не найдена"
	@echo ""
	@echo "📖 Полная документация: test-documentation.md"

# Генерация документации для примеров
docs: build
	@echo "📚 Генерация документации..."
	@./bin/$(BINARY_NAME) examples/_examples -output docs/examples.md
	@echo "✅ Документация сгенерирована: docs/examples.md"

# Запуск примеров
examples: build
	@echo "🔧 Запуск примеров использования..."
	@cd cmd/examples/basic && go run main.go
	@echo ""
	@cd cmd/examples/advanced && go run main.go
	@echo "✅ Примеры выполнены"

# Docker сборка
docker-build:
	@echo "🐳 Сборка Docker образа..."
	@echo "ℹ️  Используется: $(DOCKER_CMD)"
	@if ! $(DOCKER_CMD) build -t $(DOCKER_IMAGE):latest . 2>/dev/null; then \
		echo ""; \
		echo "❌ Ошибка сборки Docker образа!"; \
		echo ""; \
		echo "💡 Возможные решения:"; \
		echo "   1. Проверьте права доступа к Docker:"; \
		echo "      sudo usermod -aG docker $$USER && newgrp docker"; \
		echo ""; \
		echo "   2. Используйте готовый образ:"; \
		echo "      docker pull seblex5/testdoc:latest"; \
		echo ""; \
		echo "   3. Или используйте локальную сборку:"; \
		echo "      make build && ./bin/$(BINARY_NAME) examples/_examples"; \
		echo ""; \
		exit 1; \
	fi
	@echo "✅ Docker образ собран: $(DOCKER_IMAGE):latest"

# Проверка Docker
docker-check:
	@echo "🔍 Проверка Docker..."
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "❌ Docker не установлен!"; \
		echo "   Установите Docker: https://docs.docker.com/get-docker/"; \
		exit 1; \
	fi
	@if ! docker version >/dev/null 2>&1; then \
		echo "❌ Docker недоступен для текущего пользователя!"; \
		echo ""; \
		echo "📋 Для исправления выполните:"; \
		echo "   sudo usermod -aG docker $$USER"; \
		echo "   newgrp docker"; \
		echo ""; \
		echo "   Или перелогиньтесь в систему."; \
		echo ""; \
		echo "⚠️  Подробнее: https://docs.docker.com/engine/install/linux-postinstall/"; \
		exit 1; \
	fi
	@if ! docker build --help >/dev/null 2>&1; then \
		echo "❌ Команда 'docker build' недоступна!"; \
		echo "   Проверьте установку Docker и права доступа."; \
		exit 1; \
	fi
	@echo "✅ Docker готов к работе"

# Docker запуск
docker-run: docker-check docker-build
	@echo "🐳 Запуск в Docker..."
	@$(DOCKER_CMD) run --rm -v $$(pwd):/workspace $(DOCKER_IMAGE):latest /workspace/examples/_examples

# Подготовка к релизу
release-prep: clean fmt vet lint test-cover
	@echo "🚀 Подготовка к релизу..."
	@git status --porcelain | grep -q . && echo "❌ Есть неcommitted изменения" && exit 1 || echo "✅ Рабочая директория чистая"
	@echo "✅ Готово к релизу"

# CI проверки
ci: check test-cover
	@echo "🤖 CI проверки завершены"

# Развертывание локальной разработки
dev-setup:
	@echo "🛠️ Настройка среды разработки..."
	@go mod download
	@go mod verify
	@command -v golangci-lint >/dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	@echo "✅ Среда разработки настроена"

# Бенчмарки
bench:
	@echo "⚡ Запуск бенчмарков..."
	@go test -bench=. -benchmem ./...

# Профилирование
profile:
	@echo "📊 Профилирование..."
	@go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./pkg/parser
	@echo "Используйте: go tool pprof cpu.prof или go tool pprof mem.prof"

# Генерация моков (если используются)
generate:
	@echo "🔄 Генерация кода..."
	@go generate ./...

# Справка
help:
	@echo "TestDoc - Генератор документации для Go тестов"
	@echo ""
	@echo "🏗️  Сборка:"
	@echo "  build           - Собрать CLI приложение"
	@echo "  build-cli       - Собрать с версией из git"
	@echo "  install         - Установить локально"
	@echo ""
	@echo "🧪 Тестирование:"
	@echo "  test            - Запуск всех тестов"
	@echo "  test-race       - Тесты с race detector"
	@echo "  test-short      - Только быстрые тесты"
	@echo "  test-cover      - Тесты с покрытием"
	@echo "  cover-html      - HTML отчет покрытия"
	@echo "  bench           - Бенчмарки"
	@echo ""
	@echo "🔍 Качество кода:"
	@echo "  lint            - Линтинг"
	@echo "  fmt             - Форматирование"
	@echo "  vet             - Статический анализ"
	@echo "  check           - Полная проверка"
	@echo ""
	@echo "📚 Документация:"
	@echo "  demo            - Демонстрация работы"
	@echo "  docs            - Генерация документации"
	@echo "  examples        - Запуск примеров"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  docker-check    - Проверка настройки Docker"
	@echo "  docker-build    - Сборка Docker образа"
	@echo "  docker-run      - Запуск в Docker"
	@echo ""
	@echo "🛠️  Разработка:"
	@echo "  dev-setup       - Настройка среды разработки"
	@echo "  tidy            - Обновление зависимостей"
	@echo "  generate        - Генерация кода"
	@echo "  clean           - Очистка файлов"
	@echo ""
	@echo "🚀 Релиз:"
	@echo "  release-prep    - Подготовка к релизу"
	@echo "  ci              - CI проверки"
	@echo ""
	@echo "Примеры:"
	@echo "  make demo                    # Быстрая демонстрация"
	@echo "  make check                   # Полная проверка кода"
	@echo "  make test-cover             # Тесты с покрытием"
	@echo "  make docker-run             # Запуск в Docker"
