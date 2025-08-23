// Package types содержит основные типы данных для testdoc библиотеки
package types

import "time"

// TestType определяет тип теста
type TestType string

const (
	UnitTest        TestType = "unit"
	IntegrationTest TestType = "integration"
	FunctionalTest  TestType = "functional"
	E2ETest         TestType = "e2e"
	PerformanceTest TestType = "performance"
	SecurityTest    TestType = "security"
	RegressionTest  TestType = "regression"
	SmokeTest       TestType = "smoke"
)

// String возвращает строковое представление типа теста
func (t TestType) String() string {
	return string(t)
}

// IsValid проверяет, является ли тип теста валидным
func (t TestType) IsValid() bool {
	switch t {
	case UnitTest, IntegrationTest, FunctionalTest, E2ETest,
		PerformanceTest, SecurityTest, RegressionTest, SmokeTest:
		return true
	default:
		return false
	}
}

// TestInfo содержит информацию о тесте
type TestInfo struct {
	Name        string            `json:"name" yaml:"name"`
	Type        TestType          `json:"type" yaml:"type"`
	Description string            `json:"description" yaml:"description"`
	TestCases   []TestCase        `json:"test_cases" yaml:"test_cases"`
	Skipped     bool              `json:"skipped" yaml:"skipped"`
	SkipReason  string            `json:"skip_reason,omitempty" yaml:"skip_reason,omitempty"`
	Package     string            `json:"package" yaml:"package"`
	File        string            `json:"file" yaml:"file"`
	Line        int               `json:"line" yaml:"line"`
	Tags        []string          `json:"tags" yaml:"tags"`
	Author      string            `json:"author,omitempty" yaml:"author,omitempty"`
	Created     time.Time         `json:"created,omitempty" yaml:"created,omitempty"`
	Updated     time.Time         `json:"updated,omitempty" yaml:"updated,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// TestCase представляет отдельный тест-кейс
type TestCase struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Input       string `json:"input,omitempty" yaml:"input,omitempty"`
	Expected    string `json:"expected,omitempty" yaml:"expected,omitempty"`
	Steps       []Step `json:"steps,omitempty" yaml:"steps,omitempty"`
}

// Step представляет шаг в тест-кейсе
type Step struct {
	Action      string `json:"action" yaml:"action"`
	Description string `json:"description" yaml:"description"`
	Expected    string `json:"expected,omitempty" yaml:"expected,omitempty"`
}

// TestAnnotation представляет аннотацию в комментарии
type TestAnnotation struct {
	Type  string `json:"type" yaml:"type"`
	Value string `json:"value" yaml:"value"`
}

// PackageInfo содержит информацию о пакете с тестами
type PackageInfo struct {
	Name        string     `json:"name" yaml:"name"`
	Path        string     `json:"path" yaml:"path"`
	Description string     `json:"description" yaml:"description"`
	Tests       []TestInfo `json:"tests" yaml:"tests"`
	Coverage    float64    `json:"coverage,omitempty" yaml:"coverage,omitempty"`
	TestTypes   []TestType `json:"test_types" yaml:"test_types"`
}

// Config содержит настройки генерации документации
type Config struct {
	Title           string            `yaml:"title"`
	Author          string            `yaml:"author"`
	Version         string            `yaml:"version"`
	IncludeSkipped  bool              `yaml:"include_skipped"`
	GroupByType     bool              `yaml:"group_by_type"`
	GroupByPackage  bool              `yaml:"group_by_package"`
	CustomTemplates map[string]string `yaml:"custom_templates"`
	ExcludePatterns []string          `yaml:"exclude_patterns"`
	IncludePatterns []string          `yaml:"include_patterns"`
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		Title:           "Test Documentation",
		Author:          "Generated automatically",
		Version:         "1.0.0",
		IncludeSkipped:  true,
		GroupByType:     true,
		GroupByPackage:  false,
		CustomTemplates: make(map[string]string),
		ExcludePatterns: []string{},
		IncludePatterns: []string{"*_test.go"},
	}
}

// ParseResult содержит результат парсинга тестов
type ParseResult struct {
	Packages map[string]*PackageInfo `json:"packages" yaml:"packages"`
	Stats    Statistics              `json:"statistics" yaml:"statistics"`
}

// Statistics содержит статистику тестов
type Statistics struct {
	TotalTests       int              `json:"total_tests" yaml:"total_tests"`
	ActiveTests      int              `json:"active_tests" yaml:"active_tests"`
	SkippedTests     int              `json:"skipped_tests" yaml:"skipped_tests"`
	PackageCount     int              `json:"package_count" yaml:"package_count"`
	TypeDistribution map[TestType]int `json:"type_distribution" yaml:"type_distribution"`
}

// CalculateStats вычисляет статистику из результата парсинга
func (pr *ParseResult) CalculateStats() {
	stats := Statistics{
		TypeDistribution: make(map[TestType]int),
	}

	for _, pkg := range pr.Packages {
		stats.PackageCount++
		for _, test := range pkg.Tests {
			stats.TotalTests++
			if test.Skipped {
				stats.SkippedTests++
			} else {
				stats.ActiveTests++
			}
			stats.TypeDistribution[test.Type]++
		}
	}

	pr.Stats = stats
}
