package testdoc

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/seblex/testdoc/pkg/types"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	require.NotNil(t, config)
	assert.Equal(t, "Test Documentation", config.Title)
	assert.Equal(t, "Generated automatically", config.Author)
	assert.Equal(t, "1.0.0", config.Version)
	assert.True(t, config.IncludeSkipped)
	assert.True(t, config.GroupByType)
	assert.False(t, config.GroupByPackage)
}

func TestLoadConfig(t *testing.T) {
	// Создаем временный файл конфигурации
	configContent := `title: "Custom Test Documentation"
author: "Custom Author"
version: "2.0.0"
include_skipped: false
group_by_type: false
group_by_package: true
include_patterns:
  - "*_test.go"
exclude_patterns:
  - "*_bench_test.go"
custom_templates:
  header: "Custom header template"
`

	tmpDir, err := ioutil.TempDir("", "testdoc_config_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	configFile := filepath.Join(tmpDir, "config.yaml")
	err = ioutil.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// Загружаем конфигурацию
	config, err := LoadConfig(configFile)
	require.NoError(t, err)

	assert.Equal(t, "Custom Test Documentation", config.Title)
	assert.Equal(t, "Custom Author", config.Author)
	assert.Equal(t, "2.0.0", config.Version)
	assert.False(t, config.IncludeSkipped)
	assert.False(t, config.GroupByType)
	assert.True(t, config.GroupByPackage)
	assert.Contains(t, config.IncludePatterns, "*_test.go")
	assert.Contains(t, config.ExcludePatterns, "*_bench_test.go")
	assert.Equal(t, "Custom header template", config.CustomTemplates["header"])
}

func TestLoadConfig_NonExistentFile(t *testing.T) {
	_, err := LoadConfig("/non/existent/file.yaml")
	assert.Error(t, err)
}

func TestSaveConfig(t *testing.T) {
	config := &types.Config{
		Title:           "Test Config",
		Author:          "Test Author",
		Version:         "1.0.0",
		IncludeSkipped:  true,
		GroupByType:     true,
		GroupByPackage:  false,
		IncludePatterns: []string{"*_test.go"},
		ExcludePatterns: []string{"*_bench_test.go"},
		CustomTemplates: map[string]string{
			"header": "Custom header",
		},
	}

	tmpDir, err := ioutil.TempDir("", "testdoc_save_config_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	configFile := filepath.Join(tmpDir, "config.yaml")
	err = SaveConfig(config, configFile)
	require.NoError(t, err)

	// Проверяем, что файл создан и содержит правильные данные
	_, err = os.Stat(configFile)
	assert.NoError(t, err)

	// Загружаем обратно и проверяем
	loadedConfig, err := LoadConfig(configFile)
	require.NoError(t, err)

	assert.Equal(t, config.Title, loadedConfig.Title)
	assert.Equal(t, config.Author, loadedConfig.Author)
	assert.Equal(t, config.Version, loadedConfig.Version)
	assert.Equal(t, config.IncludeSkipped, loadedConfig.IncludeSkipped)
}

func TestParseFile(t *testing.T) {
	// Создаем временный тест-файл
	testCode := `package testpkg

import "testing"

// @type: unit
// @author: Test Author
// TestExample tests basic functionality
func TestExample(t *testing.T) {
	// Test implementation
}
`

	tmpDir, err := ioutil.TempDir("", "testdoc_parse_file_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "example_test.go")
	err = ioutil.WriteFile(testFile, []byte(testCode), 0644)
	require.NoError(t, err)

	// Парсим файл
	tests, err := ParseFile(testFile)
	require.NoError(t, err)

	assert.Len(t, tests, 1)
	assert.Equal(t, "TestExample", tests[0].Name)
	assert.Equal(t, types.UnitTest, tests[0].Type)
	assert.Equal(t, "Test Author", tests[0].Author)
}

func TestParseDirectory(t *testing.T) {
	// Создаем временную директорию с тест-файлами
	tmpDir, err := ioutil.TempDir("", "testdoc_parse_dir_test")
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

// @type: integration
// TestIntegration is an integration test
func TestIntegration(t *testing.T) {
	t.Skip("Skipped for demo")
}
`

	testFile := filepath.Join(tmpDir, "example_test.go")
	err = ioutil.WriteFile(testFile, []byte(testCode), 0644)
	require.NoError(t, err)

	// Парсим директорию
	result, err := ParseDirectory(tmpDir, nil)
	require.NoError(t, err)

	assert.NotNil(t, result)
	assert.Len(t, result.Packages, 1)
	assert.Contains(t, result.Packages, "testpkg")

	pkg := result.Packages["testpkg"]
	assert.Len(t, pkg.Tests, 2)

	// Проверяем статистику
	assert.Equal(t, 2, result.Stats.TotalTests)
	assert.Equal(t, 1, result.Stats.ActiveTests)
	assert.Equal(t, 1, result.Stats.SkippedTests)
	assert.Equal(t, 1, result.Stats.PackageCount)
}

func TestGenerateFromDirectory(t *testing.T) {
	// Создаем временную директорию с тест-файлами
	tmpDir, err := ioutil.TempDir("", "testdoc_generate_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Создаем тест-файл
	testCode := `package testpkg

import "testing"

// @type: unit
// @author: Test Author
// TestExample tests basic functionality
func TestExample(t *testing.T) {
	// Test implementation
}
`

	testFile := filepath.Join(tmpDir, "example_test.go")
	err = ioutil.WriteFile(testFile, []byte(testCode), 0644)
	require.NoError(t, err)

	// Генерируем документацию
	config := DefaultConfig()
	config.Title = "Test Documentation"

	markdown, err := GenerateFromDirectory(tmpDir, config)
	require.NoError(t, err)

	assert.Contains(t, markdown, "# Test Documentation")
	assert.Contains(t, markdown, "TestExample")
	assert.Contains(t, markdown, "Test Author")
	assert.Contains(t, markdown, "## Статистика тестов")
}

func TestWriteToFile(t *testing.T) {
	content := "# Test Documentation\n\nThis is a test document."

	tmpDir, err := ioutil.TempDir("", "testdoc_write_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	outputFile := filepath.Join(tmpDir, "test.md")
	err = WriteToFile(content, outputFile)
	require.NoError(t, err)

	// Проверяем, что файл создан и содержит правильные данные
	writtenContent, err := ioutil.ReadFile(outputFile)
	require.NoError(t, err)
	assert.Equal(t, content, string(writtenContent))
}

func TestAppendToFile(t *testing.T) {
	initialContent := "# Test Documentation\n\n"
	additionalContent := "## Additional Section\n\nMore content."

	tmpDir, err := ioutil.TempDir("", "testdoc_append_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	outputFile := filepath.Join(tmpDir, "test.md")

	// Записываем начальный контент
	err = WriteToFile(initialContent, outputFile)
	require.NoError(t, err)

	// Добавляем дополнительный контент
	err = AppendToFile(additionalContent, outputFile)
	require.NoError(t, err)

	// Проверяем результат
	finalContent, err := ioutil.ReadFile(outputFile)
	require.NoError(t, err)
	assert.Equal(t, initialContent+additionalContent, string(finalContent))
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *types.Config
		check  func(*testing.T, *types.Config)
	}{
		{
			name:   "empty_config",
			config: &types.Config{},
			check: func(t *testing.T, c *types.Config) {
				assert.Equal(t, "Test Documentation", c.Title)
				assert.Equal(t, "Generated automatically", c.Author)
				assert.Equal(t, "1.0.0", c.Version)
				assert.NotNil(t, c.ExcludePatterns)
				assert.NotNil(t, c.IncludePatterns)
				assert.NotNil(t, c.CustomTemplates)
			},
		},
		{
			name: "partial_config",
			config: &types.Config{
				Title: "Custom Title",
			},
			check: func(t *testing.T, c *types.Config) {
				assert.Equal(t, "Custom Title", c.Title)
				assert.Equal(t, "Generated automatically", c.Author)
				assert.Equal(t, "1.0.0", c.Version)
			},
		},
		{
			name: "nil_slices",
			config: &types.Config{
				Title:           "Test",
				ExcludePatterns: nil,
				IncludePatterns: nil,
				CustomTemplates: nil,
			},
			check: func(t *testing.T, c *types.Config) {
				assert.NotNil(t, c.ExcludePatterns)
				assert.NotNil(t, c.IncludePatterns)
				assert.NotNil(t, c.CustomTemplates)
				assert.Contains(t, c.IncludePatterns, "*_test.go")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			assert.NoError(t, err)
			tt.check(t, tt.config)
		})
	}
}

func TestGetSupportedTestTypes(t *testing.T) {
	testTypes := GetSupportedTestTypes()

	assert.Len(t, testTypes, 8)
	assert.Contains(t, testTypes, types.UnitTest)
	assert.Contains(t, testTypes, types.IntegrationTest)
	assert.Contains(t, testTypes, types.FunctionalTest)
	assert.Contains(t, testTypes, types.E2ETest)
	assert.Contains(t, testTypes, types.PerformanceTest)
	assert.Contains(t, testTypes, types.SecurityTest)
	assert.Contains(t, testTypes, types.RegressionTest)
	assert.Contains(t, testTypes, types.SmokeTest)
}

func TestIsValidTestType(t *testing.T) {
	tests := []struct {
		testType types.TestType
		expected bool
	}{
		{types.UnitTest, true},
		{types.IntegrationTest, true},
		{types.FunctionalTest, true},
		{types.E2ETest, true},
		{types.PerformanceTest, true},
		{types.SecurityTest, true},
		{types.RegressionTest, true},
		{types.SmokeTest, true},
		{types.TestType("unknown"), false},
		{types.TestType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.testType), func(t *testing.T) {
			result := IsValidTestType(tt.testType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStatistics_CalculateTestCoverage(t *testing.T) {
	stats := NewStatistics()

	result := &types.ParseResult{
		Stats: types.Statistics{
			TotalTests: 10,
			TypeDistribution: map[types.TestType]int{
				types.UnitTest:        6,
				types.IntegrationTest: 3,
				types.E2ETest:         1,
			},
		},
	}

	coverage := stats.CalculateTestCoverage(result)

	assert.Equal(t, 60.0, coverage[types.UnitTest])
	assert.Equal(t, 30.0, coverage[types.IntegrationTest])
	assert.Equal(t, 10.0, coverage[types.E2ETest])
}

func TestStatistics_CalculateTestCoverage_EmptyResult(t *testing.T) {
	stats := NewStatistics()

	result := &types.ParseResult{
		Stats: types.Statistics{TotalTests: 0},
	}

	coverage := stats.CalculateTestCoverage(result)
	assert.Empty(t, coverage)
}

func TestStatistics_GetMostCommonTestType(t *testing.T) {
	stats := NewStatistics()

	result := &types.ParseResult{
		Stats: types.Statistics{
			TypeDistribution: map[types.TestType]int{
				types.UnitTest:        6,
				types.IntegrationTest: 3,
				types.E2ETest:         1,
			},
		},
	}

	mostCommon, count := stats.GetMostCommonTestType(result)
	assert.Equal(t, types.UnitTest, mostCommon)
	assert.Equal(t, 6, count)
}

func TestFilter_ByType(t *testing.T) {
	filter := NewFilter()

	// Создаем тестовые данные
	testInfo1 := types.TestInfo{Name: "TestUnit", Type: types.UnitTest}
	testInfo2 := types.TestInfo{Name: "TestIntegration", Type: types.IntegrationTest}
	testInfo3 := types.TestInfo{Name: "TestUnit2", Type: types.UnitTest}

	result := &types.ParseResult{
		Packages: map[string]*types.PackageInfo{
			"pkg1": {
				Name:  "pkg1",
				Tests: []types.TestInfo{testInfo1, testInfo2},
			},
			"pkg2": {
				Name:  "pkg2",
				Tests: []types.TestInfo{testInfo3},
			},
		},
	}
	result.CalculateStats()

	// Фильтруем только unit тесты
	filtered := filter.ByType(result, types.UnitTest)

	assert.Len(t, filtered.Packages, 2) // Оба пакета содержат unit тесты
	assert.Len(t, filtered.Packages["pkg1"].Tests, 1)
	assert.Len(t, filtered.Packages["pkg2"].Tests, 1)
	assert.Equal(t, "TestUnit", filtered.Packages["pkg1"].Tests[0].Name)
	assert.Equal(t, "TestUnit2", filtered.Packages["pkg2"].Tests[0].Name)
	assert.Equal(t, 2, filtered.Stats.TotalTests)
}

func TestFilter_ByTags(t *testing.T) {
	filter := NewFilter()

	testInfo1 := types.TestInfo{
		Name: "TestAPI",
		Type: types.UnitTest,
		Tags: []string{"api", "unit"},
	}
	testInfo2 := types.TestInfo{
		Name: "TestDB",
		Type: types.IntegrationTest,
		Tags: []string{"database", "integration"},
	}
	testInfo3 := types.TestInfo{
		Name: "TestAPIIntegration",
		Type: types.IntegrationTest,
		Tags: []string{"api", "integration"},
	}

	result := &types.ParseResult{
		Packages: map[string]*types.PackageInfo{
			"pkg1": {
				Name:  "pkg1",
				Tests: []types.TestInfo{testInfo1, testInfo2, testInfo3},
			},
		},
	}
	result.CalculateStats()

	// Фильтруем тесты с тегом "api"
	filtered := filter.ByTags(result, []string{"api"})

	assert.Len(t, filtered.Packages, 1)
	assert.Len(t, filtered.Packages["pkg1"].Tests, 2)
	assert.Equal(t, "TestAPI", filtered.Packages["pkg1"].Tests[0].Name)
	assert.Equal(t, "TestAPIIntegration", filtered.Packages["pkg1"].Tests[1].Name)
	assert.Equal(t, 2, filtered.Stats.TotalTests)
}

func TestFilter_ByAuthor(t *testing.T) {
	filter := NewFilter()

	testInfo1 := types.TestInfo{
		Name:   "TestByJohn",
		Author: "John Doe",
	}
	testInfo2 := types.TestInfo{
		Name:   "TestByJane",
		Author: "Jane Smith",
	}
	testInfo3 := types.TestInfo{
		Name:   "TestByJohn2",
		Author: "John Doe",
	}

	result := &types.ParseResult{
		Packages: map[string]*types.PackageInfo{
			"pkg1": {
				Name:  "pkg1",
				Tests: []types.TestInfo{testInfo1, testInfo2, testInfo3},
			},
		},
	}
	result.CalculateStats()

	// Фильтруем тесты автора "John Doe"
	filtered := filter.ByAuthor(result, "John Doe")

	assert.Len(t, filtered.Packages, 1)
	assert.Len(t, filtered.Packages["pkg1"].Tests, 2)
	assert.Equal(t, "TestByJohn", filtered.Packages["pkg1"].Tests[0].Name)
	assert.Equal(t, "TestByJohn2", filtered.Packages["pkg1"].Tests[1].Name)
	assert.Equal(t, 2, filtered.Stats.TotalTests)
}
