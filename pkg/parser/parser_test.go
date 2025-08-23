package parser

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/testdoc-org/testdoc/pkg/types"
)

func TestParser_New(t *testing.T) {
	parser := New()
	assert.NotNil(t, parser)
	assert.NotNil(t, parser.fileSet)
}

func TestParser_isTestFunction(t *testing.T) {
	parser := New()

	tests := []struct {
		name     string
		funcName string
		expected bool
	}{
		{"test_function", "TestExample", true},
		{"benchmark_function", "BenchmarkExample", true},
		{"example_function", "ExampleExample", true},
		{"regular_function", "RegularFunction", false},
		{"test_prefix_but_lowercase", "testExample", false},
		{"empty_name", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.isTestFunction(tt.funcName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParser_isTestFile(t *testing.T) {
	parser := New()

	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{"test_file", "example_test.go", true},
		{"test_file_with_path", "/path/to/example_test.go", true},
		{"regular_file", "example.go", false},
		{"test_in_middle", "test_example.go", false},
		{"empty_filename", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.isTestFile(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParser_shouldExcludeFile(t *testing.T) {
	parser := New()

	tests := []struct {
		name     string
		path     string
		config   *types.Config
		expected bool
	}{
		{
			name: "include_pattern_match",
			path: "/path/to/example_test.go",
			config: &types.Config{
				IncludePatterns: []string{"*_test.go"},
				ExcludePatterns: []string{},
			},
			expected: false,
		},
		{
			name: "include_pattern_no_match",
			path: "/path/to/example.go",
			config: &types.Config{
				IncludePatterns: []string{"*_test.go"},
				ExcludePatterns: []string{},
			},
			expected: true,
		},
		{
			name: "exclude_pattern_match",
			path: "/path/to/bench_test.go",
			config: &types.Config{
				IncludePatterns: []string{"*_test.go"},
				ExcludePatterns: []string{"bench_*"},
			},
			expected: true,
		},
		{
			name: "no_patterns",
			path: "/path/to/example_test.go",
			config: &types.Config{
				IncludePatterns: []string{},
				ExcludePatterns: []string{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.shouldExcludeFile(tt.path, tt.config)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParser_ParseFile(t *testing.T) {
	// Создаем временный тест-файл
	testCode := `package testpkg

import "testing"

// @type: unit
// @author: Test Author
// @created: 2024-01-15
// @tags: unit,example
// @testcase: Simple test - tests basic functionality
// @step: Call function - should return expected result
// TestExample tests basic functionality
func TestExample(t *testing.T) {
	// Test implementation
}

// @type: integration
// @skip_reason: Database not available
// TestIntegration tests database integration
func TestIntegration(t *testing.T) {
	t.Skip("Database not available")
}

// Regular function, should be ignored
func RegularFunction() {
	// Not a test
}

// @type: performance
// BenchmarkPerformance benchmarks performance
func BenchmarkPerformance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Benchmark implementation
	}
}
`

	// Создаем временный файл
	tmpDir, err := ioutil.TempDir("", "parser_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "example_test.go")
	err = ioutil.WriteFile(testFile, []byte(testCode), 0644)
	require.NoError(t, err)

	// Парсим файл
	parser := New()
	tests, err := parser.ParseFile(testFile)
	require.NoError(t, err)

	// Проверяем результаты
	assert.Len(t, tests, 3) // TestExample, TestIntegration, BenchmarkPerformance

	// Проверяем первый тест (TestExample)
	testExample := findTestByName(tests, "TestExample")
	require.NotNil(t, testExample)
	assert.Equal(t, "TestExample", testExample.Name)
	assert.Equal(t, types.UnitTest, testExample.Type)
	assert.Equal(t, "Test Author", testExample.Author)
	assert.Equal(t, "testpkg", testExample.Package)
	assert.Equal(t, "example_test.go", testExample.File)
	assert.False(t, testExample.Skipped)
	assert.Contains(t, testExample.Tags, "unit")
	assert.Contains(t, testExample.Tags, "example")
	assert.Len(t, testExample.TestCases, 1)
	assert.Equal(t, "Simple test", testExample.TestCases[0].Name)
	assert.Equal(t, "tests basic functionality", testExample.TestCases[0].Description)

	// Проверяем created дату
	expectedDate, _ := time.Parse("2006-01-02", "2024-01-15")
	assert.Equal(t, expectedDate, testExample.Created)

	// Проверяем второй тест (TestIntegration)
	testIntegration := findTestByName(tests, "TestIntegration")
	require.NotNil(t, testIntegration)
	assert.Equal(t, "TestIntegration", testIntegration.Name)
	assert.Equal(t, types.IntegrationTest, testIntegration.Type)
	assert.True(t, testIntegration.Skipped)
	assert.Equal(t, "Database not available", testIntegration.SkipReason)

	// Проверяем benchmark
	benchmarkTest := findTestByName(tests, "BenchmarkPerformance")
	require.NotNil(t, benchmarkTest)
	assert.Equal(t, "BenchmarkPerformance", benchmarkTest.Name)
	assert.Equal(t, types.PerformanceTest, benchmarkTest.Type)
	assert.False(t, benchmarkTest.Skipped)
}

func TestParser_ParseFile_InvalidFile(t *testing.T) {
	parser := New()

	// Тестируем несуществующий файл
	_, err := parser.ParseFile("/non/existent/file.go")
	assert.Error(t, err)
}

func TestParser_parseAnnotation(t *testing.T) {
	parser := New()
	testInfo := &types.TestInfo{
		Tags:     []string{},
		Metadata: make(map[string]string),
	}

	tests := []struct {
		name  string
		line  string
		check func(*testing.T, *types.TestInfo)
	}{
		{
			name: "type_annotation",
			line: "@type: integration",
			check: func(t *testing.T, ti *types.TestInfo) {
				assert.Equal(t, types.IntegrationTest, ti.Type)
			},
		},
		{
			name: "author_annotation",
			line: "@author: John Doe",
			check: func(t *testing.T, ti *types.TestInfo) {
				assert.Equal(t, "John Doe", ti.Author)
			},
		},
		{
			name: "tags_annotation",
			line: "@tags: api,database,unit",
			check: func(t *testing.T, ti *types.TestInfo) {
				assert.Contains(t, ti.Tags, "api")
				assert.Contains(t, ti.Tags, "database")
				assert.Contains(t, ti.Tags, "unit")
			},
		},
		{
			name: "testcase_annotation",
			line: "@testcase: Test name - Test description",
			check: func(t *testing.T, ti *types.TestInfo) {
				require.Len(t, ti.TestCases, 1)
				assert.Equal(t, "Test name", ti.TestCases[0].Name)
				assert.Equal(t, "Test description", ti.TestCases[0].Description)
			},
		},
		{
			name: "custom_metadata",
			line: "@priority: high",
			check: func(t *testing.T, ti *types.TestInfo) {
				assert.Equal(t, "high", ti.Metadata["priority"])
			},
		},
		{
			name: "created_date",
			line: "@created: 2024-01-15",
			check: func(t *testing.T, ti *types.TestInfo) {
				expectedDate, _ := time.Parse("2006-01-02", "2024-01-15")
				assert.Equal(t, expectedDate, ti.Created)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сбрасываем состояние testInfo для каждого теста
			testInfo = &types.TestInfo{
				Tags:     []string{},
				Metadata: make(map[string]string),
			}

			parser.parseAnnotation(tt.line, testInfo)
			tt.check(t, testInfo)
		})
	}
}

func TestParser_ParseDirectory(t *testing.T) {
	// Создаем временную директорию с тест-файлами
	tmpDir, err := ioutil.TempDir("", "parser_dir_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Создаем тест-файл
	testCode := `package testpkg

import "testing"

// @type: unit
// TestUnit is a unit test
func TestUnit(t *testing.T) {
	// Test implementation
}
`

	testFile := filepath.Join(tmpDir, "example_test.go")
	err = ioutil.WriteFile(testFile, []byte(testCode), 0644)
	require.NoError(t, err)

	// Создаем обычный файл (не тест)
	regularCode := `package testpkg

func RegularFunction() {
	// Regular implementation
}
`

	regularFile := filepath.Join(tmpDir, "example.go")
	err = ioutil.WriteFile(regularFile, []byte(regularCode), 0644)
	require.NoError(t, err)

	// Парсим директорию
	parser := New()
	config := types.DefaultConfig()
	result, err := parser.ParseDirectory(tmpDir, config)
	require.NoError(t, err)

	// Проверяем результаты
	assert.NotNil(t, result)
	assert.Len(t, result.Packages, 1)
	assert.Contains(t, result.Packages, "testpkg")

	pkg := result.Packages["testpkg"]
	assert.Equal(t, "testpkg", pkg.Name)
	assert.Len(t, pkg.Tests, 1)
	assert.Equal(t, "TestUnit", pkg.Tests[0].Name)
	assert.Equal(t, types.UnitTest, pkg.Tests[0].Type)

	// Проверяем статистику
	assert.Equal(t, 1, result.Stats.TotalTests)
	assert.Equal(t, 1, result.Stats.ActiveTests)
	assert.Equal(t, 0, result.Stats.SkippedTests)
	assert.Equal(t, 1, result.Stats.PackageCount)
}

// Вспомогательная функция для поиска теста по имени
func findTestByName(tests []types.TestInfo, name string) *types.TestInfo {
	for i := range tests {
		if tests[i].Name == name {
			return &tests[i]
		}
	}
	return nil
}
