package generator

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/seblex/testdoc/pkg/types"
)

func TestGenerator_New(t *testing.T) {
	// Тест с конфигурацией по умолчанию
	gen := New(nil)
	assert.NotNil(t, gen)
	assert.NotNil(t, gen.config)
	assert.Equal(t, "Test Documentation", gen.config.Title)

	// Тест с пользовательской конфигурацией
	config := &types.Config{
		Title:  "Custom Title",
		Author: "Custom Author",
	}
	gen = New(config)
	assert.NotNil(t, gen)
	assert.Equal(t, "Custom Title", gen.config.Title)
	assert.Equal(t, "Custom Author", gen.config.Author)
}

func TestGenerator_GenerateMarkdown(t *testing.T) {
	// Создаем тестовые данные
	testInfo1 := types.TestInfo{
		Name:        "TestExample",
		Type:        types.UnitTest,
		Description: "Example test description",
		Package:     "example",
		File:        "example_test.go",
		Line:        10,
		Author:      "John Doe",
		Created:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Tags:        []string{"unit", "example"},
		TestCases: []types.TestCase{
			{
				Name:        "Basic test case",
				Description: "Tests basic functionality",
				Steps: []types.Step{
					{
						Action:   "Call function",
						Expected: "Return expected result",
					},
				},
			},
		},
		Metadata: map[string]string{
			"priority": "high",
		},
	}

	testInfo2 := types.TestInfo{
		Name:       "TestIntegration",
		Type:       types.IntegrationTest,
		Package:    "example",
		File:       "integration_test.go",
		Line:       20,
		Skipped:    true,
		SkipReason: "Database not available",
	}

	pkg := &types.PackageInfo{
		Name:      "example",
		Path:      "/path/to/example",
		Tests:     []types.TestInfo{testInfo1, testInfo2},
		TestTypes: []types.TestType{types.UnitTest, types.IntegrationTest},
	}

	result := &types.ParseResult{
		Packages: map[string]*types.PackageInfo{
			"example": pkg,
		},
		Stats: types.Statistics{
			TotalTests:   2,
			ActiveTests:  1,
			SkippedTests: 1,
			PackageCount: 1,
			TypeDistribution: map[types.TestType]int{
				types.UnitTest:        1,
				types.IntegrationTest: 1,
			},
		},
	}

	config := &types.Config{
		Title:       "Test Documentation",
		Author:      "Test Author",
		Version:     "1.0.0",
		GroupByType: true,
	}

	gen := New(config)
	markdown := gen.GenerateMarkdown(result)

	// Проверяем содержание
	assert.Contains(t, markdown, "# Test Documentation")
	assert.Contains(t, markdown, "**Автор:** Test Author")
	assert.Contains(t, markdown, "**Версия:** 1.0.0")
	assert.Contains(t, markdown, "## Оглавление")
	assert.Contains(t, markdown, "## Статистика тестов")
	assert.Contains(t, markdown, "- **Всего тестов:** 2")
	assert.Contains(t, markdown, "- **Активных тестов:** 1")
	assert.Contains(t, markdown, "- **Пропущенных тестов:** 1")
	assert.Contains(t, markdown, "### TestExample")
	assert.Contains(t, markdown, "### TestIntegration")
	assert.Contains(t, markdown, "⏭️ Пропущен")
	assert.Contains(t, markdown, "Database not available")
	assert.Contains(t, markdown, "John Doe")
	assert.Contains(t, markdown, "2024-01-15")
	assert.Contains(t, markdown, "`unit`, `example`")
	assert.Contains(t, markdown, "Basic test case")
	assert.Contains(t, markdown, "Call function → Return expected result")
	assert.Contains(t, markdown, "**Priority:** high")
}

func TestGenerator_getTestTypeDisplayName(t *testing.T) {
	gen := New(nil)

	tests := []struct {
		testType types.TestType
		expected string
	}{
		{types.UnitTest, "Модульные"},
		{types.IntegrationTest, "Интеграционные"},
		{types.FunctionalTest, "Функциональные"},
		{types.E2ETest, "E2E"},
		{types.PerformanceTest, "Производительности"},
		{types.SecurityTest, "Безопасности"},
		{types.RegressionTest, "Регрессионные"},
		{types.SmokeTest, "Дымовые"},
		{types.TestType("unknown"), "unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.testType), func(t *testing.T) {
			result := gen.getTestTypeDisplayName(tt.testType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerator_groupTestsByType(t *testing.T) {
	gen := New(nil)

	testInfo1 := types.TestInfo{Name: "TestUnit1", Type: types.UnitTest}
	testInfo2 := types.TestInfo{Name: "TestUnit2", Type: types.UnitTest}
	testInfo3 := types.TestInfo{Name: "TestIntegration", Type: types.IntegrationTest}

	packages := map[string]*types.PackageInfo{
		"pkg1": {
			Tests: []types.TestInfo{testInfo1, testInfo3},
		},
		"pkg2": {
			Tests: []types.TestInfo{testInfo2},
		},
	}

	groups := gen.groupTestsByType(packages)

	assert.Len(t, groups, 2)
	assert.Len(t, groups[types.UnitTest], 2)
	assert.Len(t, groups[types.IntegrationTest], 1)

	// Проверяем сортировку
	unitTests := groups[types.UnitTest]
	assert.Equal(t, "TestUnit1", unitTests[0].Name)
	assert.Equal(t, "TestUnit2", unitTests[1].Name)
}

func TestGenerator_generateTOCByType(t *testing.T) {
	gen := New(&types.Config{GroupByType: true})

	testInfo1 := types.TestInfo{Name: "TestUnit", Type: types.UnitTest}
	testInfo2 := types.TestInfo{Name: "TestIntegration", Type: types.IntegrationTest}

	packages := map[string]*types.PackageInfo{
		"example": {
			Tests: []types.TestInfo{testInfo1, testInfo2},
		},
	}

	var sb strings.Builder
	gen.generateTOCByType(&sb, packages)

	toc := sb.String()
	assert.Contains(t, toc, "- [Интеграционные тесты](#integration-тесты)")
	assert.Contains(t, toc, "- [Модульные тесты](#unit-тесты)")
	assert.Contains(t, toc, "- [TestUnit](#testunit)")
	assert.Contains(t, toc, "- [TestIntegration](#testintegration)")
}

func TestGenerator_generateStatistics(t *testing.T) {
	gen := New(nil)

	stats := &types.Statistics{
		TotalTests:   10,
		ActiveTests:  8,
		SkippedTests: 2,
		PackageCount: 3,
		TypeDistribution: map[types.TestType]int{
			types.UnitTest:        6,
			types.IntegrationTest: 3,
			types.E2ETest:         1,
		},
	}

	var sb strings.Builder
	gen.generateStatistics(&sb, stats)

	output := sb.String()
	assert.Contains(t, output, "- **Всего тестов:** 10")
	assert.Contains(t, output, "- **Активных тестов:** 8")
	assert.Contains(t, output, "- **Пропущенных тестов:** 2")
	assert.Contains(t, output, "- **Пакетов:** 3")
	assert.Contains(t, output, "### Распределение по типам")
	assert.Contains(t, output, "**E2E:** 1 (10.0%)")
	assert.Contains(t, output, "**Интеграционные:** 3 (30.0%)")
	assert.Contains(t, output, "**Модульные:** 6 (60.0%)")
}

func TestGenerator_generateTestSection(t *testing.T) {
	gen := New(nil)

	testInfo := types.TestInfo{
		Name:        "TestExample",
		Type:        types.UnitTest,
		Description: "Example test description",
		Package:     "example",
		File:        "example_test.go",
		Line:        10,
		Author:      "John Doe",
		Created:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Updated:     time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		Tags:        []string{"unit", "example"},
		Skipped:     false,
		TestCases: []types.TestCase{
			{
				Name:        "Basic functionality",
				Description: "Tests basic operations",
				Input:       "test input",
				Expected:    "expected output",
				Steps: []types.Step{
					{
						Action:   "Setup test data",
						Expected: "Data initialized",
					},
					{
						Action:   "Execute function",
						Expected: "Function returns result",
					},
				},
			},
		},
		Metadata: map[string]string{
			"priority":   "high",
			"complexity": "low",
		},
	}

	var sb strings.Builder
	gen.generateTestSection(&sb, testInfo)

	output := sb.String()

	// Проверяем основную информацию
	assert.Contains(t, output, "### TestExample")
	assert.Contains(t, output, "| **Тип** | Модульные |")
	assert.Contains(t, output, "| **Пакет** | `example` |")
	assert.Contains(t, output, "| **Файл** | `example_test.go:10` |")
	assert.Contains(t, output, "| **Статус** | ✅ Активен |")
	assert.Contains(t, output, "| **Автор** | John Doe |")
	assert.Contains(t, output, "| **Создан** | 2024-01-15 |")
	assert.Contains(t, output, "| **Обновлен** | 2024-01-20 |")
	assert.Contains(t, output, "| **Теги** | `unit`, `example` |")

	// Проверяем описание
	assert.Contains(t, output, "**Описание:**")
	assert.Contains(t, output, "Example test description")

	// Проверяем тест-кейсы
	assert.Contains(t, output, "#### Тест-кейсы")
	assert.Contains(t, output, "**1. Basic functionality**")
	assert.Contains(t, output, "Tests basic operations")
	assert.Contains(t, output, "- **Входные данные:** test input")
	assert.Contains(t, output, "- **Ожидаемый результат:** expected output")
	assert.Contains(t, output, "**Шаги:**")
	assert.Contains(t, output, "1. Setup test data → Data initialized")
	assert.Contains(t, output, "2. Execute function → Function returns result")

	// Проверяем метаданные
	assert.Contains(t, output, "#### Дополнительная информация")
	assert.Contains(t, output, "- **Priority:** high")
	assert.Contains(t, output, "- **Complexity:** low")

	// Проверяем разделитель
	assert.Contains(t, output, "---")
}

func TestGenerator_generateTestSection_Skipped(t *testing.T) {
	gen := New(nil)

	testInfo := types.TestInfo{
		Name:       "TestSkipped",
		Type:       types.IntegrationTest,
		Package:    "example",
		File:       "example_test.go",
		Line:       20,
		Skipped:    true,
		SkipReason: "External service unavailable",
	}

	var sb strings.Builder
	gen.generateTestSection(&sb, testInfo)

	output := sb.String()
	assert.Contains(t, output, "| **Статус** | ⏭️ Пропущен |")
	assert.Contains(t, output, "| **Причина пропуска** | External service unavailable |")
}

func TestGenerator_generateContentByPackage(t *testing.T) {
	config := &types.Config{GroupByPackage: true}
	gen := New(config)

	testInfo := types.TestInfo{
		Name:    "TestExample",
		Type:    types.UnitTest,
		Package: "example",
		File:    "example_test.go",
		Line:    10,
	}

	packages := map[string]*types.PackageInfo{
		"example": {
			Name:        "example",
			Path:        "/path/to/example",
			Description: "Example package description",
			Tests:       []types.TestInfo{testInfo},
		},
	}

	var sb strings.Builder
	gen.generateContentByPackage(&sb, packages)

	output := sb.String()
	assert.Contains(t, output, "## Пакет example")
	assert.Contains(t, output, "Example package description")
	assert.Contains(t, output, "**Путь:** `/path/to/example`")
	assert.Contains(t, output, "### TestExample")
}

func TestGenerator_generateContentByType(t *testing.T) {
	config := &types.Config{GroupByType: true}
	gen := New(config)

	testInfo1 := types.TestInfo{Name: "TestUnit", Type: types.UnitTest}
	testInfo2 := types.TestInfo{Name: "TestIntegration", Type: types.IntegrationTest}

	packages := map[string]*types.PackageInfo{
		"example": {
			Tests: []types.TestInfo{testInfo1, testInfo2},
		},
	}

	var sb strings.Builder
	gen.generateContentByType(&sb, packages)

	output := sb.String()
	assert.Contains(t, output, "## Интеграционные тесты")
	assert.Contains(t, output, "## Модульные тесты")
	assert.Contains(t, output, "### TestUnit")
	assert.Contains(t, output, "### TestIntegration")
}

func TestGenerator_generateSimpleContent(t *testing.T) {
	gen := New(nil)

	testInfo1 := types.TestInfo{Name: "TestB", Type: types.UnitTest}
	testInfo2 := types.TestInfo{Name: "TestA", Type: types.IntegrationTest}

	packages := map[string]*types.PackageInfo{
		"example": {
			Tests: []types.TestInfo{testInfo1, testInfo2},
		},
	}

	var sb strings.Builder
	gen.generateSimpleContent(&sb, packages)

	output := sb.String()
	assert.Contains(t, output, "## Тесты")

	// Проверяем, что тесты отсортированы по имени
	testAPos := strings.Index(output, "### TestA")
	testBPos := strings.Index(output, "### TestB")
	assert.True(t, testAPos < testBPos, "TestA should come before TestB")
}

func TestGenerator_getLanguage(t *testing.T) {
	tests := []struct {
		name     string
		language string
		expected string
	}{
		{
			name:     "russian_short",
			language: "ru",
			expected: "ru",
		},
		{
			name:     "russian_full",
			language: "russian",
			expected: "ru",
		},
		{
			name:     "english_short",
			language: "en",
			expected: "en",
		},
		{
			name:     "english_full",
			language: "english",
			expected: "en",
		},
		{
			name:     "unknown_language",
			language: "unknown",
			expected: "ru", // По умолчанию русский
		},
		{
			name:     "empty_language",
			language: "",
			expected: "ru", // По умолчанию русский
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &types.Config{Language: tt.language}
			gen := New(config)
			lang := gen.getLanguage()
			assert.Equal(t, tt.expected, lang.String())
		})
	}
}

func TestGenerator_LanguageInMetadata(t *testing.T) {
	// Тест для русского языка (по умолчанию)
	testInfo := types.TestInfo{
		Name:    "TestExample",
		Type:    types.UnitTest,
		Package: "example",
		File:    "example_test.go",
		Line:    10,
		Metadata: map[string]string{
			"priority":   "высокий",
			"complexity": "низкий",
		},
	}

	configRu := &types.Config{Language: "ru"}
	genRu := New(configRu)

	var sbRu strings.Builder
	genRu.generateTestSection(&sbRu, testInfo)
	outputRu := sbRu.String()

	// Для русского языка заголовки должны быть с заглавной буквы по правилам русского языка
	assert.Contains(t, outputRu, "- **Priority:** высокий")
	assert.Contains(t, outputRu, "- **Complexity:** низкий")

	// Тест для английского языка
	testInfoEn := types.TestInfo{
		Name:    "TestExample",
		Type:    types.UnitTest,
		Package: "example",
		File:    "example_test.go",
		Line:    10,
		Metadata: map[string]string{
			"priority":   "high",
			"complexity": "low",
		},
	}

	configEn := &types.Config{Language: "en"}
	genEn := New(configEn)

	var sbEn strings.Builder
	genEn.generateTestSection(&sbEn, testInfoEn)
	outputEn := sbEn.String()

	// Для английского языка заголовки должны быть с заглавной буквы по правилам английского языка
	assert.Contains(t, outputEn, "- **Priority:** high")
	assert.Contains(t, outputEn, "- **Complexity:** low")
}
