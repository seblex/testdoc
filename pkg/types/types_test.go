package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestType_String(t *testing.T) {
	tests := []struct {
		name     string
		testType TestType
		expected string
	}{
		{"unit", UnitTest, "unit"},
		{"integration", IntegrationTest, "integration"},
		{"functional", FunctionalTest, "functional"},
		{"e2e", E2ETest, "e2e"},
		{"performance", PerformanceTest, "performance"},
		{"security", SecurityTest, "security"},
		{"regression", RegressionTest, "regression"},
		{"smoke", SmokeTest, "smoke"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.testType.String())
		})
	}
}

func TestTestType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		testType TestType
		expected bool
	}{
		{"valid_unit", UnitTest, true},
		{"valid_integration", IntegrationTest, true},
		{"valid_functional", FunctionalTest, true},
		{"valid_e2e", E2ETest, true},
		{"valid_performance", PerformanceTest, true},
		{"valid_security", SecurityTest, true},
		{"valid_regression", RegressionTest, true},
		{"valid_smoke", SmokeTest, true},
		{"invalid_empty", TestType(""), false},
		{"invalid_unknown", TestType("unknown"), false},
		{"invalid_random", TestType("random_type"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.testType.IsValid())
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	require.NotNil(t, config)
	assert.Equal(t, "Test Documentation", config.Title)
	assert.Equal(t, "Generated automatically", config.Author)
	assert.Equal(t, "1.0.0", config.Version)
	assert.True(t, config.IncludeSkipped)
	assert.True(t, config.GroupByType)
	assert.False(t, config.GroupByPackage)
	assert.NotNil(t, config.CustomTemplates)
	assert.NotNil(t, config.ExcludePatterns)
	assert.NotNil(t, config.IncludePatterns)
	assert.Contains(t, config.IncludePatterns, "*_test.go")
}

func TestParseResult_CalculateStats(t *testing.T) {
	// Создаем тестовые данные
	testInfo1 := TestInfo{
		Name:    "TestUnit1",
		Type:    UnitTest,
		Skipped: false,
	}
	testInfo2 := TestInfo{
		Name:    "TestUnit2",
		Type:    UnitTest,
		Skipped: true,
	}
	testInfo3 := TestInfo{
		Name:    "TestIntegration1",
		Type:    IntegrationTest,
		Skipped: false,
	}

	pkg1 := &PackageInfo{
		Name:  "package1",
		Tests: []TestInfo{testInfo1, testInfo2},
	}
	pkg2 := &PackageInfo{
		Name:  "package2",
		Tests: []TestInfo{testInfo3},
	}

	result := &ParseResult{
		Packages: map[string]*PackageInfo{
			"package1": pkg1,
			"package2": pkg2,
		},
	}

	// Вычисляем статистику
	result.CalculateStats()

	// Проверяем результаты
	assert.Equal(t, 3, result.Stats.TotalTests)
	assert.Equal(t, 2, result.Stats.ActiveTests)
	assert.Equal(t, 1, result.Stats.SkippedTests)
	assert.Equal(t, 2, result.Stats.PackageCount)

	assert.Equal(t, 2, result.Stats.TypeDistribution[UnitTest])
	assert.Equal(t, 1, result.Stats.TypeDistribution[IntegrationTest])
	assert.Equal(t, 0, result.Stats.TypeDistribution[FunctionalTest])
}

func TestTestInfo_Creation(t *testing.T) {
	now := time.Now()

	testInfo := TestInfo{
		Name:        "TestExample",
		Type:        UnitTest,
		Description: "Test description",
		Package:     "example",
		File:        "example_test.go",
		Line:        10,
		Author:      "John Doe",
		Created:     now,
		Tags:        []string{"unit", "example"},
		Metadata:    map[string]string{"priority": "high"},
	}

	assert.Equal(t, "TestExample", testInfo.Name)
	assert.Equal(t, UnitTest, testInfo.Type)
	assert.Equal(t, "Test description", testInfo.Description)
	assert.Equal(t, "example", testInfo.Package)
	assert.Equal(t, "example_test.go", testInfo.File)
	assert.Equal(t, 10, testInfo.Line)
	assert.Equal(t, "John Doe", testInfo.Author)
	assert.Equal(t, now, testInfo.Created)
	assert.Contains(t, testInfo.Tags, "unit")
	assert.Contains(t, testInfo.Tags, "example")
	assert.Equal(t, "high", testInfo.Metadata["priority"])
}

func TestTestCase_Creation(t *testing.T) {
	step1 := Step{
		Action:   "Do something",
		Expected: "Expected result",
	}
	step2 := Step{
		Action:      "Do another thing",
		Description: "Step description",
		Expected:    "Another expected result",
	}

	testCase := TestCase{
		Name:        "Test case example",
		Description: "Test case description",
		Input:       "input data",
		Expected:    "expected output",
		Steps:       []Step{step1, step2},
	}

	assert.Equal(t, "Test case example", testCase.Name)
	assert.Equal(t, "Test case description", testCase.Description)
	assert.Equal(t, "input data", testCase.Input)
	assert.Equal(t, "expected output", testCase.Expected)
	assert.Len(t, testCase.Steps, 2)
	assert.Equal(t, "Do something", testCase.Steps[0].Action)
	assert.Equal(t, "Expected result", testCase.Steps[0].Expected)
	assert.Equal(t, "Do another thing", testCase.Steps[1].Action)
	assert.Equal(t, "Step description", testCase.Steps[1].Description)
}

func TestPackageInfo_Creation(t *testing.T) {
	testInfo := TestInfo{
		Name: "TestExample",
		Type: UnitTest,
	}

	pkg := PackageInfo{
		Name:        "example",
		Path:        "/path/to/example",
		Description: "Example package",
		Tests:       []TestInfo{testInfo},
		Coverage:    85.5,
		TestTypes:   []TestType{UnitTest, IntegrationTest},
	}

	assert.Equal(t, "example", pkg.Name)
	assert.Equal(t, "/path/to/example", pkg.Path)
	assert.Equal(t, "Example package", pkg.Description)
	assert.Len(t, pkg.Tests, 1)
	assert.Equal(t, 85.5, pkg.Coverage)
	assert.Contains(t, pkg.TestTypes, UnitTest)
	assert.Contains(t, pkg.TestTypes, IntegrationTest)
}

func TestStatistics_EmptyResults(t *testing.T) {
	result := &ParseResult{
		Packages: map[string]*PackageInfo{},
	}

	result.CalculateStats()

	assert.Equal(t, 0, result.Stats.TotalTests)
	assert.Equal(t, 0, result.Stats.ActiveTests)
	assert.Equal(t, 0, result.Stats.SkippedTests)
	assert.Equal(t, 0, result.Stats.PackageCount)
	assert.Empty(t, result.Stats.TypeDistribution)
}
