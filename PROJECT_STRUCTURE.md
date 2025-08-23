# Структура проекта TestDoc

Этот документ описывает организацию файлов и директорий в проекте TestDoc.

## 📁 Структура директорий

```
testdoc/
├── 📁 .github/                    # GitHub интеграция
│   └── 📁 workflows/              # GitHub Actions
│       ├── ci.yml                 # Continuous Integration
│       └── release.yml            # Автоматические релизы
├── 📁 cmd/                        # CLI приложения
│   └── 📁 testdoc/                # Основное CLI приложение
│       └── main.go                # Точка входа CLI
├── 📁 pkg/                        # Библиотечные пакеты
│   ├── 📁 generator/              # Генерация документации
│   │   ├── generator.go           # Основная логика генерации
│   │   └── generator_test.go      # Unit тесты
│   ├── 📁 parser/                 # Парсинг Go файлов
│   │   ├── parser.go              # AST анализ тестов
│   │   └── parser_test.go         # Unit тесты
│   └── 📁 types/                  # Типы данных
│       ├── types.go               # Определения типов
│       └── types_test.go          # Unit тесты
├── 📁 examples/                   # Примеры использования
│   ├── 📁 basic/                  # Базовый пример
│   │   └── main.go                # Простое использование
│   ├── 📁 advanced/               # Продвинутый пример
│   │   └── main.go                # Фильтрация и статистика
│   └── 📁 _examples/              # Тестовые данные
│       ├── user_service_test.go   # Примеры тестов
│       └── payment_test.go        # Примеры тестов
├── 📄 testdoc.go                  # Публичный API библиотеки
├── 📄 testdoc_test.go            # Тесты публичного API
├── 📄 go.mod                      # Go модуль
├── 📄 go.sum                      # Зависимости
├── 📄 config.yaml                 # Пример конфигурации
├── 📄 Dockerfile                  # Docker образ
├── 📄 .dockerignore              # Docker игнорирование
├── 📄 .gitignore                 # Git игнорирование
├── 📄 Makefile                   # Автоматизация сборки
├── 📄 README.md                  # Основная документация
├── 📄 CONTRIBUTING.md            # Руководство для участников
├── 📄 CHANGELOG.md               # История изменений
├── 📄 LICENSE                    # MIT лицензия
└── 📄 PROJECT_STRUCTURE.md       # Этот файл
```

## 🏗️ Архитектура

### 🎯 Основные компоненты

1. **📦 Core Library** (`testdoc.go`)
   - Публичный API для внешнего использования
   - Упрощенные функции для быстрого старта
   - Фасад для внутренних пакетов

2. **🔍 Parser** (`pkg/parser/`)
   - AST анализ Go тест-файлов
   - Извлечение аннотаций из комментариев
   - Определение пропущенных тестов

3. **📝 Generator** (`pkg/generator/`)
   - Генерация Markdown документации
   - Шаблонизация и форматирование
   - Группировка и сортировка тестов

4. **📋 Types** (`pkg/types/`)
   - Определения структур данных
   - Конфигурация и настройки
   - Статистика и метаданные

5. **⚡ CLI** (`cmd/testdoc/`)
   - Интерфейс командной строки
   - Парсинг аргументов и флагов
   - Интеграция всех компонентов

### 🔄 Поток данных

```
Go Test Files
     ↓
   Parser ──→ ParseResult
     ↓             ↓
   Types ←──→ Generator
     ↓             ↓
Configuration  Markdown
```

## 📦 Пакеты

### `github.com/seblex5/testdoc`

**Назначение**: Основной пакет с публичным API

**Экспортируемые функции**:
- `DefaultConfig()` - конфигурация по умолчанию
- `LoadConfig()/SaveConfig()` - работа с конфигурацией
- `ParseDirectory()/ParseFile()` - парсинг тестов
- `GenerateMarkdown()` - генерация документации
- `WriteToFile()/AppendToFile()` - сохранение файлов
- `NewFilter()/NewStatistics()` - утилиты

### `github.com/seblex5/testdoc/pkg/types`

**Назначение**: Типы данных и структуры

**Основные типы**:
- `TestType` - типы тестов (unit, integration, etc.)
- `TestInfo` - информация о тесте
- `TestCase` - тест-кейс с шагами
- `PackageInfo` - информация о пакете
- `Config` - конфигурация генерации
- `ParseResult` - результат парсинга

### `github.com/seblex5/testdoc/pkg/parser`

**Назначение**: Анализ Go тест-файлов

**Ключевые возможности**:
- AST парсинг Go файлов
- Извлечение аннотаций (`@type`, `@author`, etc.)
- Анализ пропущенных тестов (`t.Skip()`)
- Обход директорий с фильтрацией

### `github.com/seblex5/testdoc/pkg/generator`

**Назначение**: Генерация Markdown документации

**Возможности**:
- Создание оглавления
- Статистика тестов
- Группировка по типам/пакетам
- Форматирование тест-кейсов
- Поддержка пользовательских шаблонов

## 🛠️ Инструменты разработки

### Makefile цели

| Цель | Описание |
|------|----------|
| `build` | Сборка CLI приложения |
| `test` | Запуск всех тестов |
| `test-cover` | Тесты с покрытием |
| `lint` | Статический анализ |
| `demo` | Демонстрация работы |
| `docker-build` | Сборка Docker образа |
| `release-prep` | Подготовка к релизу |

### GitHub Actions

- **CI Pipeline**: тестирование на Go 1.20-1.22
- **Security**: Gosec сканирование
- **Release**: автоматическая сборка для всех платформ
- **Docker**: multi-arch образы

## 📊 Метрики качества

- **Покрытие тестами**: >85% для всех пакетов
- **Линтинг**: golangci-lint с строгими правилами
- **Документация**: 100% покрытие публичного API
- **Безопасность**: регулярное сканирование зависимостей

## 🚀 Deployment

### Релизы

1. **Автоматическая сборка** для Linux, macOS, Windows
2. **Docker образы** с поддержкой multi-arch
3. **Homebrew formula** для macOS/Linux
4. **Go modules** для программной интеграции

### Установка

```bash
# Go install
go install github.com/seblex5/testdoc/cmd/testdoc@latest

# Homebrew
brew install testdoc

# Docker
docker pull seblex5/testdoc:latest

# Binary download
curl -L https://github.com/seblex5/testdoc/releases/latest/download/testdoc-linux-amd64.tar.gz | tar xz
```

## 🔧 Расширение

### Добавление нового типа теста

1. Добавить в `pkg/types/types.go`:
   ```go
   ContractTest TestType = "contract"
   ```

2. Обновить `IsValid()` метод

3. Добавить в `generator.go`:
   ```go
   case types.ContractTest:
       return "Контрактные"
   ```

### Новая аннотация

1. Добавить парсинг в `pkg/parser/parser.go`:
   ```go
   if strings.HasPrefix(line, "@priority:") {
       // parsing logic
   }
   ```

2. Обновить структуру в `pkg/types/types.go`

3. Добавить в документацию

## 📚 Документация

- **README.md** - основная документация
- **CONTRIBUTING.md** - руководство для участников
- **CHANGELOG.md** - история изменений
- **examples/** - примеры использования
- **pkg/*/**.go - GoDoc комментарии

## 🔒 Безопасность

- **Статический анализ**: Gosec в CI
- **Зависимости**: регулярная проверка уязвимостей
- **Docker**: non-root пользователь, минимальный образ
- **Релизы**: подписание артефактов (планируется)

---

Эта структура обеспечивает:
- ✅ **Модульность** - четкое разделение ответственности
- ✅ **Тестируемость** - высокое покрытие тестами
- ✅ **Расширяемость** - простое добавление новых функций
- ✅ **Документированность** - полная документация API
- ✅ **Качество** - автоматические проверки в CI
