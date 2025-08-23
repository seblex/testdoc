# Участие в развитии TestDoc 🤝

Спасибо за интерес к участию в развитии TestDoc! Мы приветствуем все виды вклада: от отчетов об ошибках до новых функций.

## 📋 Содержание

- [Code of Conduct](#code-of-conduct)
- [Как помочь](#как-помочь)
- [Сообщение об ошибках](#сообщение-об-ошибках)
- [Предложение улучшений](#предложение-улучшений)
- [Разработка](#разработка)
- [Pull Request Process](#pull-request-process)
- [Стиль кодирования](#стиль-кодирования)
- [Тестирование](#тестирование)
- [Документация](#документация)

## 📜 Code of Conduct

Этот проект придерживается [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). Участвуя, вы соглашаетесь соблюдать этот кодекс.

## 🎯 Как помочь

Существует множество способов внести свой вклад в TestDoc:

### 🐛 Сообщение об ошибках
- Сообщайте о найденных багах через [GitHub Issues](https://github.com/seblex5/testdoc/issues)
- Убедитесь, что ошибка еще не была сообщена
- Используйте шаблон для отчета об ошибке

### 💡 Предложение новых функций
- Обсуждайте идеи в [GitHub Discussions](https://github.com/seblex5/testdoc/discussions)
- Создавайте feature requests через [GitHub Issues](https://github.com/seblex5/testdoc/issues)
- Используйте шаблон для предложения функций

### 🔧 Код
- Исправление багов
- Реализация новых функций
- Улучшение производительности
- Рефакторинг кода

### 📖 Документация
- Улучшение README
- Добавление примеров
- Создание туториалов
- Исправление опечаток

### 🧪 Тестирование
- Написание unit тестов
- Создание интеграционных тестов
- Тестирование на разных платформах

## 🐛 Сообщение об ошибках

При сообщении об ошибке, пожалуйста, включите:

### Обязательная информация
- **Версия TestDoc**: `testdoc -version`
- **Версия Go**: `go version`
- **Операционная система**: (Linux, macOS, Windows)
- **Архитектура**: (amd64, arm64, etc.)

### Описание проблемы
- Четкое и краткое описание проблемы
- Шаги для воспроизведения
- Ожидаемое поведение
- Фактическое поведение
- Скриншоты (если применимо)

### Пример отчета
```markdown
## Описание
TestDoc падает при анализе файлов с Unicode символами в комментариях.

## Шаги для воспроизведения
1. Создать файл test_unicode_test.go с комментарием `// Тест с русскими символами`
2. Запустить `testdoc ./`
3. Наблюдать ошибку

## Ожидаемое поведение
TestDoc должен корректно обрабатывать Unicode символы в комментариях.

## Фактическое поведение
```
panic: runtime error: invalid memory address
...
```

## Окружение
- TestDoc версия: v1.0.0
- Go версия: go1.21.0
- ОС: Ubuntu 20.04 LTS
```

## 💡 Предложение улучшений

При предложении новых функций:

1. **Проверьте существующие issues** - возможно, функция уже обсуждается
2. **Объясните use case** - зачем нужна эта функция
3. **Предложите API** - как должна работать функция
4. **Рассмотрите альтернативы** - есть ли другие способы решения

### Шаблон предложения
```markdown
## Краткое описание
Добавить поддержку экспорта в HTML формат.

## Мотивация
Пользователи хотят просматривать документацию в браузере с интерактивными элементами.

## Детальное описание
- Новый флаг `--format html`
- HTML шаблоны с CSS стилями
- Интерактивное оглавление
- Поиск по документации

## Возможная реализация
```go
type HTMLGenerator struct {
    templates map[string]*template.Template
}

func (g *HTMLGenerator) Generate(result *ParseResult) (string, error) {
    // implementation
}
```

## Альтернативы
- Использование внешних конверторов Markdown → HTML
- Интеграция с существующими генераторами документации
```

## 🏗️ Разработка

### Требования
- **Go 1.20+**
- **Git**
- **golangci-lint** (для линтинга)

### Настройка окружения

1. **Fork репозитория**
   ```bash
   # Через GitHub UI или
   gh repo fork testdoc-org/testdoc
   ```

2. **Клонирование**
   ```bash
   git clone https://github.com/YOUR_USERNAME/testdoc.git
   cd testdoc
   ```

3. **Установка зависимостей**
   ```bash
   go mod download
   go mod verify
   ```

4. **Проверка сборки**
   ```bash
   go build ./cmd/testdoc
   ./testdoc -version
   ```

### Структура проекта

```
testdoc/
├── cmd/testdoc/           # CLI приложение
├── pkg/                   # Библиотечные пакеты
│   ├── parser/           # Парсинг Go файлов
│   ├── generator/        # Генерация документации
│   └── types/            # Типы данных
├── examples/             # Примеры использования
├── internal/             # Внутренние пакеты
├── .github/              # GitHub Actions
├── docs/                 # Документация
└── testdoc.go           # Основной API
```

### Ветки

- `main` - стабильная ветка с релизами
- `develop` - ветка разработки
- `feature/название` - ветки для новых функций
- `bugfix/название` - ветки для исправления багов
- `hotfix/название` - срочные исправления

### Workflow разработки

1. **Создайте ветку**
   ```bash
   git checkout -b feature/my-awesome-feature
   ```

2. **Делайте изменения**
   ```bash
   # Пишите код
   # Добавляйте тесты
   # Обновляйте документацию
   ```

3. **Проверяйте качество**
   ```bash
   # Тесты
   go test ./...
   
   # Линтинг
   golangci-lint run
   
   # Форматирование
   go fmt ./...
   
   # Проверка модулей
   go mod tidy
   ```

4. **Коммитьте изменения**
   ```bash
   git add .
   git commit -m "feat: add HTML export support"
   ```

5. **Отправляйте в fork**
   ```bash
   git push origin feature/my-awesome-feature
   ```

6. **Создавайте Pull Request**

## 🔄 Pull Request Process

### Подготовка PR

1. **Убедитесь в актуальности ветки**
   ```bash
   git fetch upstream
   git rebase upstream/develop
   ```

2. **Проверьте все тесты**
   ```bash
   go test ./... -race -coverprofile=coverage.out
   go tool cover -func=coverage.out
   ```

3. **Проверьте линтинг**
   ```bash
   golangci-lint run
   ```

4. **Обновите документацию** (если необходимо)

### Требования к PR

- ✅ Все тесты проходят
- ✅ Покрытие кода не уменьшается
- ✅ Линтинг проходит без ошибок
- ✅ Коммиты имеют осмысленные сообщения
- ✅ PR имеет четкое описание
- ✅ Документация обновлена (если необходимо)

### Шаблон описания PR

```markdown
## Описание
Краткое описание изменений.

## Тип изменения
- [ ] Исправление бага (non-breaking change)
- [ ] Новая функция (non-breaking change)
- [ ] Breaking change (изменение, нарушающее обратную совместимость)
- [ ] Обновление документации

## Как тестировать
1. Шаги для тестирования изменений
2. ...

## Чеклист
- [ ] Код следует стилю проекта
- [ ] Код самодокументируемый или добавлены комментарии
- [ ] Добавлены тесты для новой функциональности
- [ ] Все новые и существующие тесты проходят
- [ ] Документация обновлена

## Связанные issues
Closes #123
Fixes #456
```

### Процесс ревью

1. **Автоматические проверки** должны пройти
2. **Code review** от мейнтейнеров
3. **Обсуждение** и внесение изменений
4. **Одобрение** и merge

## 🎨 Стиль кодирования

### Go код

Следуйте стандартным соглашениям Go:

- **gofmt** для форматирования
- **golint** для проверки стиля
- **go vet** для статического анализа

### Именование

```go
// ✅ Хорошо
type TestInfo struct {
    Name string
}

func ParseTestFile(filename string) (*TestInfo, error) {
    // ...
}

// ❌ Плохо
type testinfo struct {
    name string
}

func parseTestfile(filename string) (*testinfo, error) {
    // ...
}
```

### Комментарии

```go
// ✅ Хорошо - комментарий к пакету
// Package parser предоставляет функциональность для анализа Go тест-файлов

// ParseFile анализирует один тест-файл и возвращает информацию о тестах
func ParseFile(filename string) ([]TestInfo, error) {
    // ...
}

// ❌ Плохо - нет комментариев к публичным функциям
func ParseFile(filename string) ([]TestInfo, error) {
    // ...
}
```

### Обработка ошибок

```go
// ✅ Хорошо
result, err := SomeFunction()
if err != nil {
    return fmt.Errorf("failed to process: %w", err)
}

// ❌ Плохо
result, _ := SomeFunction() // игнорирование ошибки
```

## 🧪 Тестирование

### Типы тестов

1. **Unit тесты** - тестируют отдельные функции
2. **Integration тесты** - тестируют взаимодействие компонентов
3. **End-to-end тесты** - тестируют полный workflow

### Написание тестов

```go
func TestParseAnnotation(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected TestInfo
    }{
        {
            name:  "type annotation",
            input: "@type: unit",
            expected: TestInfo{Type: UnitTest},
        },
        // ...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := parseAnnotation(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Запуск тестов

```bash
# Все тесты
go test ./...

# Конкретный пакет
go test ./pkg/parser

# С покрытием
go test -coverprofile=coverage.out ./...

# Только быстрые тесты
go test -short ./...

# С race detector
go test -race ./...

# Verbose output
go test -v ./...
```

### Требования к покрытию

- **Минимальное покрытие**: 80%
- **Новый код**: должен иметь тесты
- **Критический код**: 90%+ покрытие

## 📖 Документация

### Типы документации

1. **README.md** - общая информация о проекте
2. **GoDoc** - документация к API
3. **Примеры** - код примеры в `examples/`
4. **Туториалы** - пошаговые руководства

### Обновление документации

При изменении API:
- Обновите комментарии к функциям
- Добавьте примеры использования
- Обновите README.md
- Создайте changelog запись

### Примеры кода

```go
// Пример должен быть исполняемым
func ExampleParseFile() {
    tests, err := ParseFile("example_test.go")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d tests\n", len(tests))
    // Output: Found 3 tests
}
```

## 🏷️ Версионирование

Проект использует [Semantic Versioning](https://semver.org/):

- **MAJOR** версия для breaking changes
- **MINOR** версия для новых функций
- **PATCH** версия для bug fixes

### Changelog

Все изменения документируются в [CHANGELOG.md](CHANGELOG.md):

```markdown
## [1.2.0] - 2024-02-01

### Added
- HTML export support
- New filter options

### Changed
- Improved error messages
- Better performance

### Fixed
- Unicode handling in comments
- Memory leak in parser
```

## 🎉 Признание вклада

Все участники признаются в:
- [Contributors](https://github.com/seblex5/testdoc/contributors)
- [CHANGELOG.md](CHANGELOG.md)
- Релизных заметках

## 📞 Получение помощи

Если у вас есть вопросы:

1. 📖 Прочитайте [документацию](README.md)
2. 🔍 Поищите в [существующих issues](https://github.com/seblex5/testdoc/issues)
3. 💬 Задайте вопрос в [Discussions](https://github.com/seblex5/testdoc/discussions)
4. 📧 Напишите мейнтейнерам

## 🙏 Спасибо!

Спасибо за участие в развитии TestDoc! Каждый вклад делает проект лучше. 🚀
