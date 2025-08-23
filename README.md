# TestDoc 📚

[![CI](https://github.com/seblex5/testdoc/workflows/CI/badge.svg)](https://github.com/seblex5/testdoc/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/seblex5/testdoc)](https://goreportcard.com/report/github.com/seblex5/testdoc)
[![codecov](https://codecov.io/gh/testdoc-org/testdoc/branch/main/graph/badge.svg)](https://codecov.io/gh/testdoc-org/testdoc)
[![Go Reference](https://pkg.go.dev/badge/github.com/seblex5/testdoc.svg)](https://pkg.go.dev/github.com/seblex5/testdoc)
[![Release](https://img.shields.io/github/release/testdoc-org/testdoc.svg)](https://github.com/seblex5/testdoc/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Автоматический генератор документации для Go тестов** 🚀

TestDoc анализирует ваши Go тест-файлы и создает красивую, подробную документацию в формате Markdown на основе комментариев и анализа кода.

## ✨ Основные возможности

- 🏷️ **Классификация тестов** по типам (unit, integration, functional, e2e, performance, security, regression, smoke)
- 📝 **Извлечение тест-кейсов** из аннотаций в комментариях
- ⏭️ **Автоматическое определение пропущенных тестов** с причинами
- 📊 **Детальная статистика** с распределением по типам
- 🎯 **Фильтрация** по типам, авторам, тегам
- 📚 **Красивая Markdown документация** с оглавлением и якорными ссылками
- 🔧 **Гибкая конфигурация** через YAML
- 🛠️ **CLI и библиотека** для интеграции в любой проект

## 🚀 Быстрый старт

### Установка

#### Через Go Install
```bash
go install github.com/seblex5/testdoc/cmd/testdoc@latest
```

#### Скачать бинарный файл
Скачайте последнюю версию с [страницы релизов](https://github.com/seblex5/testdoc/releases).

#### Docker
```bash
docker pull seblex5/testdoc:latest
```

### Использование

#### CLI
```bash
# Базовая генерация для текущей директории
testdoc

# Анализ конкретной директории
testdoc ./examples/_examples

# С настройками
testdoc -config config.yaml -output docs.md ./internal

# Фильтрация по типу тестов
testdoc -type unit -author "John Doe" ./pkg
```

#### Как библиотека
```go
package main

import (
    "log"
    "github.com/seblex5/testdoc"
)

func main() {
    // Создание конфигурации
    config := testdoc.DefaultConfig()
    config.Title = "Документация тестов моего проекта"
    config.Author = "Команда разработки"

    // Генерация документации
    doc, err := testdoc.GenerateFromDirectory("./examples/_examples", config)
    if err != nil {
        log.Fatal(err)
    }

    // Сохранение в файл
    err = testdoc.WriteToFile(doc, "test-documentation.md")
    if err != nil {
        log.Fatal(err)
    }
}
```

#### Docker
```bash
# Генерация документации
docker run --rm -v $(pwd):/workspace seblex5/testdoc /workspace/examples/_examples

# С настройками
docker run --rm -v $(pwd):/workspace seblex5/testdoc \
  -config /workspace/config.yaml \
  -output /workspace/docs.md \
  /workspace/examples/_examples
```

## 📝 Система аннотаций

TestDoc использует специальные аннотации в комментариях для извлечения метаданных тестов:

```go
// @type: unit
// @author: Иван Петров
// @created: 2024-01-15
// @updated: 2024-01-20
// @tags: user,validation,api
// @testcase: Валидация email - проверяет корректность email адреса
// @testcase: Пустой email - должен возвращать ошибку
// @step: Передать корректный email - функция должна вернуть true
// @step: Передать некорректный email - функция должна вернуть false
// TestValidateEmail проверяет функцию валидации email адресов
func TestValidateEmail(t *testing.T) {
    // реализация теста
}
```

### Поддерживаемые аннотации

| Аннотация | Описание | Пример |
|-----------|----------|---------|
| `@type` | Тип теста | `@type: unit` |
| `@author` | Автор теста | `@author: John Doe` |
| `@created` | Дата создания | `@created: 2024-01-15` |
| `@updated` | Дата обновления | `@updated: 2024-01-20` |
| `@tags` | Теги (через запятую) | `@tags: api,database,critical` |
| `@testcase` | Тест-кейс | `@testcase: Название - описание` |
| `@step` | Шаг тестирования | `@step: Действие - ожидаемый результат` |
| `@skip_reason` | Причина пропуска | `@skip_reason: Требует внешний API` |

### Типы тестов

- **unit** - Модульные тесты
- **integration** - Интеграционные тесты
- **functional** - Функциональные тесты
- **e2e** - End-to-end тесты
- **performance** - Тесты производительности
- **security** - Тесты безопасности
- **regression** - Регрессионные тесты
- **smoke** - Дымовые тесты

## 📊 Пример вывода

TestDoc генерирует красивую документацию с:

- 📋 **Оглавлением** с якорными ссылками
- 📊 **Статистикой** тестов по типам
- 📝 **Детальным описанием** каждого теста
- ⏭️ **Информацией о пропущенных тестах**
- 🏷️ **Метаданными** (автор, даты, теги)
- 📋 **Тест-кейсами** с пошаговыми инструкциями

```markdown
# Документация тестов

## Статистика тестов

- **Всего тестов:** 15
- **Активных тестов:** 13
- **Пропущенных тестов:** 2
- **Пакетов:** 3

### Распределение по типам
- **Модульные:** 8 (53.3%)
- **Интеграционные:** 4 (26.7%)
- **Функциональные:** 2 (13.3%)
- **E2E:** 1 (6.7%)

## Модульные тесты

### TestValidateEmail

| Параметр | Значение |
|----------|----------|
| **Тип** | Модульные |
| **Пакет** | `user` |
| **Файл** | `user_test.go:15` |
| **Статус** | ✅ Активен |
| **Автор** | Иван Петров |
| **Теги** | `user`, `validation`, `api` |

**Описание:** Проверяет функцию валидации email адресов

#### Тест-кейсы
**1. Валидация email**
Проверяет корректность email адреса

**Шаги:**
1. Передать корректный email → функция должна вернуть true
2. Передать некорректный email → функция должна вернуть false
```

## ⚙️ Конфигурация

Создайте файл `config.yaml` для настройки генерации:

```yaml
title: "Документация тестов проекта"
author: "Команда разработки"
version: "1.0.0"
include_skipped: true
group_by_type: true
group_by_package: false

# Паттерны файлов
include_patterns:
  - "*_test.go"
exclude_patterns:
  - "*_bench_test.go"
  - "vendor/*"

# Пользовательские шаблоны
custom_templates:
  test_header: "### Тест: {name}"
```

## 🔧 Интеграция в CI/CD

### GitHub Actions

```yaml
name: Generate Test Documentation

on:
  push:
    branches: [ main ]

jobs:
  docs:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Generate test documentation
      uses: testdoc-org/testdoc-action@v1
      with:
        path: './examples/_examples'
        output: 'docs/test-documentation.md'
        config: 'testdoc.yaml'
    
    - name: Commit documentation
      run: |
        git config --local user.email "action@github.com"
        git config --local user.name "GitHub Action"
        git add docs/test-documentation.md
        git commit -m "Update test documentation" || exit 0
        git push
```

### Makefile

```makefile
.PHONY: test-docs
test-docs:
	testdoc -config testdoc.yaml -output doc./examples/_examples.md ./internal

.PHONY: test-docs-check
test-docs-check:
	testdoc -config testdoc.yaml -output /tmp/test-docs.md ./internal
	diff doc./examples/_examples.md /tmp/test-docs.md || \
		(echo "❌ Документация не актуальна! Запустите 'make test-docs'" && exit 1)
```

## 🏗️ API библиотеки

### Основные функции

```go
// Парсинг и генерация
result, err := testdoc.ParseDirectory("./examples/_examples", config)
markdown := testdoc.GenerateMarkdown(result, config)

// Работа с конфигурацией
config := testdoc.DefaultConfig()
config, err := testdoc.LoadConfig("config.yaml")
err := testdoc.SaveConfig(config, "config.yaml")

// Фильтрация
filter := testdoc.NewFilter()
unitTests := filter.ByType(result, types.UnitTest)
apiTests := filter.ByTags(result, []string{"api"})
authorTests := filter.ByAuthor(result, "John Doe")

// Статистика
stats := testdoc.NewStatistics()
coverage := stats.CalculateTestCoverage(result)
mostCommon, count := stats.GetMostCommonTestType(result)
```

### Примеры использования

См. директорию [examples/](examples/) для полных примеров:

- [examples/basic/](examples/basic/) - Базовое использование
- [examples/advanced/](examples/advanced/) - Продвинутые возможности с фильтрацией

## 🧪 Разработка

### Требования

- Go 1.20 или выше
- Git

### Сборка из исходников

```bash
git clone https://github.com/seblex5/testdoc.git
cd testdoc
go build -o testdoc ./cmd/testdoc
```

### Запуск тестов

```bash
# Все тесты
go test ./...

# С покрытием
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Только unit тесты
go test -short ./...
```

### Линтинг

```bash
golangci-lint run
```

## 📈 Roadmap

- [ ] **Веб-интерфейс** для просмотра документации
- [ ] **Экспорт в другие форматы** (HTML, PDF, Confluence)
- [ ] **Анализ покрытия** интеграция
- [ ] **Шаблоны документации** для разных команд/проектов
- [ ] **API для внешних интеграций**

## 🤝 Участие в развитии

Мы приветствуем участие сообщества! См. [CONTRIBUTING.md](CONTRIBUTING.md) для деталей.

### Как помочь

- 🐛 **Сообщайте о багах** через [Issues](https://github.com/seblex5/testdoc/issues)
- 💡 **Предлагайте новые возможности** через [Discussions](https://github.com/seblex5/testdoc/discussions)
- 🔧 **Отправляйте Pull Requests**
- 📖 **Улучшайте документацию**
- ⭐ **Ставьте звезды** проекту

### Благодарности

Спасибо всем [участникам](https://github.com/seblex5/testdoc/contributors) проекта!

## 📜 Лицензия

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE) для деталей.

## 🔗 Полезные ссылки

- 📚 [Документация](https://testdoc-org.github.io/testdoc/)
- 🎯 [Примеры использования](examples/)
- 🐛 [Сообщить об ошибке](https://github.com/seblex5/testdoc/issues/new/choose)
- 💬 [Обсуждения](https://github.com/seblex5/testdoc/discussions)
- 📦 [Релизы](https://github.com/seblex5/testdoc/releases)
- 🐳 [Docker Hub](https://hub.docker.com/r/seblex5/testdoc)

---

<div align="center">

**Сделано с ❤️ командой TestDoc**

[⭐ Поставьте звезду на GitHub!](https://github.com/seblex5/testdoc)

</div>