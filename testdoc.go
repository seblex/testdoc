// Package testdoc предоставляет инструменты для автоматической генерации
// документации Go тестов на основе комментариев и анализа кода.
//
// Основные возможности:
//   - Классификация тестов по типам (unit, integration, functional, e2e и др.)
//   - Извлечение тест-кейсов из комментариев
//   - Определение пропущенных тестов с причинами
//   - Генерация детальной Markdown документации
//
// Пример использования:
//
//	import "github.com/seblex/testdoc"
//
//	// Создание конфигурации
//	config := testdoc.DefaultConfig()
//	config.Title = "Документация тестов моего проекта"
//	config.Author = "Команда разработки"
//
//	// Генерация документации
//	doc, err := testdoc.GenerateFromDirectory("./tests", config)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Сохранение в файл
//	err = testdoc.WriteToFile(doc, "test-documentation.md")
//	if err != nil {
//		log.Fatal(err)
//	}
package testdoc

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/seblex/testdoc/pkg/generator"
	"github.com/seblex/testdoc/pkg/parser"
	"github.com/seblex/testdoc/pkg/types"
)

// DefaultConfig возвращает конфигурацию по умолчанию для генерации документации
func DefaultConfig() *types.Config {
	return types.DefaultConfig()
}

// LoadConfig загружает конфигурацию из YAML файла
func LoadConfig(filename string) (*types.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config types.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig сохраняет конфигурацию в YAML файл
func SaveConfig(config *types.Config, filename string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// ParseDirectory анализирует директорию и возвращает информацию о тестах
func ParseDirectory(path string, config *types.Config) (*types.ParseResult, error) {
	if config == nil {
		config = DefaultConfig()
	}

	p := parser.New()
	return p.ParseDirectory(path, config)
}

// ParseFile анализирует один тест-файл и возвращает информацию о тестах
func ParseFile(filename string) ([]types.TestInfo, error) {
	p := parser.New()
	return p.ParseFile(filename)
}

// GenerateMarkdown генерирует Markdown документацию из результата парсинга
func GenerateMarkdown(result *types.ParseResult, config *types.Config) string {
	if config == nil {
		config = DefaultConfig()
	}

	g := generator.New(config)
	return g.GenerateMarkdown(result)
}

// GenerateFromDirectory анализирует директорию и генерирует документацию
func GenerateFromDirectory(path string, config *types.Config) (string, error) {
	result, err := ParseDirectory(path, config)
	if err != nil {
		return "", err
	}

	return GenerateMarkdown(result, config), nil
}

// WriteToFile записывает документацию в файл
func WriteToFile(content, filename string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

// AppendToFile добавляет документацию к существующему файлу
func AppendToFile(content, filename string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

// ValidateConfig проверяет корректность конфигурации
func ValidateConfig(config *types.Config) error {
	if config.Title == "" {
		config.Title = "Test Documentation"
	}
	if config.Author == "" {
		config.Author = "Generated automatically"
	}
	if config.Version == "" {
		config.Version = "1.0.0"
	}
	if config.Language == "" {
		config.Language = "ru"
	}

	// Инициализируем пустые слайсы если они nil
	if config.ExcludePatterns == nil {
		config.ExcludePatterns = []string{}
	}
	if config.IncludePatterns == nil {
		config.IncludePatterns = []string{"*_test.go"}
	}
	if config.CustomTemplates == nil {
		config.CustomTemplates = make(map[string]string)
	}

	return nil
}

// GetSupportedTestTypes возвращает список поддерживаемых типов тестов
func GetSupportedTestTypes() []types.TestType {
	return []types.TestType{
		types.UnitTest,
		types.IntegrationTest,
		types.FunctionalTest,
		types.E2ETest,
		types.PerformanceTest,
		types.SecurityTest,
		types.RegressionTest,
		types.SmokeTest,
	}
}

// IsValidTestType проверяет, является ли тип теста валидным
func IsValidTestType(testType types.TestType) bool {
	return testType.IsValid()
}

// Statistics содержит утилиты для работы со статистикой
type Statistics struct{}

// CalculateTestCoverage вычисляет покрытие тестами по типам
func (s *Statistics) CalculateTestCoverage(result *types.ParseResult) map[types.TestType]float64 {
	coverage := make(map[types.TestType]float64)

	if result.Stats.TotalTests == 0 {
		return coverage
	}

	for testType, count := range result.Stats.TypeDistribution {
		coverage[testType] = float64(count) / float64(result.Stats.TotalTests) * 100
	}

	return coverage
}

// GetMostCommonTestType возвращает наиболее часто встречающийся тип тестов
func (s *Statistics) GetMostCommonTestType(result *types.ParseResult) (types.TestType, int) {
	maxCount := 0
	var mostCommon types.TestType

	for testType, count := range result.Stats.TypeDistribution {
		if count > maxCount {
			maxCount = count
			mostCommon = testType
		}
	}

	return mostCommon, maxCount
}

// Filter предоставляет утилиты для фильтрации тестов
type Filter struct{}

// ByType фильтрует тесты по типу
func (f *Filter) ByType(result *types.ParseResult, testType types.TestType) *types.ParseResult {
	filtered := &types.ParseResult{
		Packages: make(map[string]*types.PackageInfo),
	}

	for pkgName, pkg := range result.Packages {
		filteredPkg := &types.PackageInfo{
			Name:        pkg.Name,
			Path:        pkg.Path,
			Description: pkg.Description,
			Tests:       []types.TestInfo{},
			TestTypes:   []types.TestType{},
		}

		for _, test := range pkg.Tests {
			if test.Type == testType {
				filteredPkg.Tests = append(filteredPkg.Tests, test)
			}
		}

		if len(filteredPkg.Tests) > 0 {
			filteredPkg.TestTypes = append(filteredPkg.TestTypes, testType)
			filtered.Packages[pkgName] = filteredPkg
		}
	}

	filtered.CalculateStats()
	return filtered
}

// ByTags фильтрует тесты по тегам
func (f *Filter) ByTags(result *types.ParseResult, tags []string) *types.ParseResult {
	filtered := &types.ParseResult{
		Packages: make(map[string]*types.PackageInfo),
	}

	for pkgName, pkg := range result.Packages {
		filteredPkg := &types.PackageInfo{
			Name:        pkg.Name,
			Path:        pkg.Path,
			Description: pkg.Description,
			Tests:       []types.TestInfo{},
			TestTypes:   []types.TestType{},
		}

		for _, test := range pkg.Tests {
			if f.hasAnyTag(test.Tags, tags) {
				filteredPkg.Tests = append(filteredPkg.Tests, test)

				// Добавляем тип теста в список типов пакета
				found := false
				for _, t := range filteredPkg.TestTypes {
					if t == test.Type {
						found = true
						break
					}
				}
				if !found {
					filteredPkg.TestTypes = append(filteredPkg.TestTypes, test.Type)
				}
			}
		}

		if len(filteredPkg.Tests) > 0 {
			filtered.Packages[pkgName] = filteredPkg
		}
	}

	filtered.CalculateStats()
	return filtered
}

// ByAuthor фильтрует тесты по автору
func (f *Filter) ByAuthor(result *types.ParseResult, author string) *types.ParseResult {
	filtered := &types.ParseResult{
		Packages: make(map[string]*types.PackageInfo),
	}

	for pkgName, pkg := range result.Packages {
		filteredPkg := &types.PackageInfo{
			Name:        pkg.Name,
			Path:        pkg.Path,
			Description: pkg.Description,
			Tests:       []types.TestInfo{},
			TestTypes:   []types.TestType{},
		}

		for _, test := range pkg.Tests {
			if test.Author == author {
				filteredPkg.Tests = append(filteredPkg.Tests, test)

				// Добавляем тип теста в список типов пакета
				found := false
				for _, t := range filteredPkg.TestTypes {
					if t == test.Type {
						found = true
						break
					}
				}
				if !found {
					filteredPkg.TestTypes = append(filteredPkg.TestTypes, test.Type)
				}
			}
		}

		if len(filteredPkg.Tests) > 0 {
			filtered.Packages[pkgName] = filteredPkg
		}
	}

	filtered.CalculateStats()
	return filtered
}

// hasAnyTag проверяет, есть ли у теста хотя бы один из указанных тегов
func (f *Filter) hasAnyTag(testTags, filterTags []string) bool {
	for _, filterTag := range filterTags {
		for _, testTag := range testTags {
			if testTag == filterTag {
				return true
			}
		}
	}
	return false
}

// NewStatistics создает новый экземпляр Statistics
func NewStatistics() *Statistics {
	return &Statistics{}
}

// NewFilter создает новый экземпляр Filter
func NewFilter() *Filter {
	return &Filter{}
}
