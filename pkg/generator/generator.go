// Package generator предоставляет функциональность для генерации документации
package generator

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/seblex/testdoc/pkg/types"
)

// Generator генерирует документацию для тестов
type Generator struct {
	config *types.Config
}

// New создает новый генератор документации
func New(config *types.Config) *Generator {
	if config == nil {
		config = types.DefaultConfig()
	}
	return &Generator{
		config: config,
	}
}

// GenerateMarkdown генерирует Markdown документацию из результата парсинга
func (g *Generator) GenerateMarkdown(result *types.ParseResult) string {
	var sb strings.Builder

	// Заголовок документа
	sb.WriteString(fmt.Sprintf("# %s\n\n", g.config.Title))
	sb.WriteString(fmt.Sprintf("**Автор:** %s  \n", g.config.Author))
	sb.WriteString(fmt.Sprintf("**Версия:** %s  \n", g.config.Version))
	sb.WriteString(fmt.Sprintf("**Дата генерации:** %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// Оглавление
	sb.WriteString("## Оглавление\n\n")
	if g.config.GroupByPackage {
		g.generateTOCByPackage(&sb, result.Packages)
	} else if g.config.GroupByType {
		g.generateTOCByType(&sb, result.Packages)
	} else {
		g.generateSimpleTOC(&sb, result.Packages)
	}

	// Статистика
	sb.WriteString("## Статистика тестов\n\n")
	g.generateStatistics(&sb, &result.Stats)

	// Основной контент
	if g.config.GroupByPackage {
		g.generateContentByPackage(&sb, result.Packages)
	} else if g.config.GroupByType {
		g.generateContentByType(&sb, result.Packages)
	} else {
		g.generateSimpleContent(&sb, result.Packages)
	}

	return sb.String()
}

// generateTOCByPackage генерирует оглавление по пакетам
func (g *Generator) generateTOCByPackage(sb *strings.Builder, packages map[string]*types.PackageInfo) {
	var packageNames []string
	for name := range packages {
		packageNames = append(packageNames, name)
	}
	sort.Strings(packageNames)

	for _, packageName := range packageNames {
		pkg := packages[packageName]
		sb.WriteString(fmt.Sprintf("- [Пакет %s](#пакет-%s)\n", packageName, strings.ToLower(packageName)))

		for _, test := range pkg.Tests {
			anchor := strings.ToLower(strings.ReplaceAll(test.Name, "_", "-"))
			sb.WriteString(fmt.Sprintf("  - [%s](#%s)\n", test.Name, anchor))
		}
	}
	sb.WriteString("\n")
}

// generateTOCByType генерирует оглавление по типам тестов
func (g *Generator) generateTOCByType(sb *strings.Builder, packages map[string]*types.PackageInfo) {
	typeGroups := g.groupTestsByType(packages)

	var testTypes []string
	for testType := range typeGroups {
		testTypes = append(testTypes, string(testType))
	}
	sort.Strings(testTypes)

	for _, testType := range testTypes {
		tests := typeGroups[types.TestType(testType)]
		sb.WriteString(fmt.Sprintf("- [%s тесты](#%s-тесты)\n",
			g.getTestTypeDisplayName(types.TestType(testType)),
			strings.ToLower(testType)))

		for _, test := range tests {
			anchor := strings.ToLower(strings.ReplaceAll(test.Name, "_", "-"))
			sb.WriteString(fmt.Sprintf("  - [%s](#%s)\n", test.Name, anchor))
		}
	}
	sb.WriteString("\n")
}

// generateSimpleTOC генерирует простое оглавление
func (g *Generator) generateSimpleTOC(sb *strings.Builder, packages map[string]*types.PackageInfo) {
	var allTests []types.TestInfo
	for _, pkg := range packages {
		allTests = append(allTests, pkg.Tests...)
	}

	sort.Slice(allTests, func(i, j int) bool {
		return allTests[i].Name < allTests[j].Name
	})

	for _, test := range allTests {
		anchor := strings.ToLower(strings.ReplaceAll(test.Name, "_", "-"))
		sb.WriteString(fmt.Sprintf("- [%s](#%s)\n", test.Name, anchor))
	}
	sb.WriteString("\n")
}

// generateStatistics генерирует статистику тестов
func (g *Generator) generateStatistics(sb *strings.Builder, stats *types.Statistics) {
	sb.WriteString(fmt.Sprintf("- **Всего тестов:** %d\n", stats.TotalTests))
	sb.WriteString(fmt.Sprintf("- **Активных тестов:** %d\n", stats.ActiveTests))
	sb.WriteString(fmt.Sprintf("- **Пропущенных тестов:** %d\n", stats.SkippedTests))
	sb.WriteString(fmt.Sprintf("- **Пакетов:** %d\n\n", stats.PackageCount))

	sb.WriteString("### Распределение по типам\n\n")
	var testTypeNames []string
	for testType := range stats.TypeDistribution {
		testTypeNames = append(testTypeNames, string(testType))
	}
	sort.Strings(testTypeNames)

	for _, testTypeName := range testTypeNames {
		testType := types.TestType(testTypeName)
		count := stats.TypeDistribution[testType]
		percentage := float64(count) / float64(stats.TotalTests) * 100
		sb.WriteString(fmt.Sprintf("- **%s:** %d (%.1f%%)\n",
			g.getTestTypeDisplayName(testType), count, percentage))
	}
	sb.WriteString("\n")
}

// generateContentByPackage генерирует контент, сгруппированный по пакетам
func (g *Generator) generateContentByPackage(sb *strings.Builder, packages map[string]*types.PackageInfo) {
	var packageNames []string
	for name := range packages {
		packageNames = append(packageNames, name)
	}
	sort.Strings(packageNames)

	for _, packageName := range packageNames {
		pkg := packages[packageName]
		sb.WriteString(fmt.Sprintf("## Пакет %s\n\n", packageName))

		if pkg.Description != "" {
			sb.WriteString(fmt.Sprintf("%s\n\n", pkg.Description))
		}

		sb.WriteString(fmt.Sprintf("**Путь:** `%s`\n\n", pkg.Path))

		for _, test := range pkg.Tests {
			g.generateTestSection(sb, test)
		}
	}
}

// generateContentByType генерирует контент, сгруппированный по типам тестов
func (g *Generator) generateContentByType(sb *strings.Builder, packages map[string]*types.PackageInfo) {
	typeGroups := g.groupTestsByType(packages)

	var testTypes []string
	for testType := range typeGroups {
		testTypes = append(testTypes, string(testType))
	}
	sort.Strings(testTypes)

	for _, testType := range testTypes {
		tests := typeGroups[types.TestType(testType)]
		sb.WriteString(fmt.Sprintf("## %s тесты\n\n", g.getTestTypeDisplayName(types.TestType(testType))))

		for _, test := range tests {
			g.generateTestSection(sb, test)
		}
	}
}

// generateSimpleContent генерирует простой контент
func (g *Generator) generateSimpleContent(sb *strings.Builder, packages map[string]*types.PackageInfo) {
	var allTests []types.TestInfo
	for _, pkg := range packages {
		allTests = append(allTests, pkg.Tests...)
	}

	sort.Slice(allTests, func(i, j int) bool {
		return allTests[i].Name < allTests[j].Name
	})

	sb.WriteString("## Тесты\n\n")
	for _, test := range allTests {
		g.generateTestSection(sb, test)
	}
}

// generateTestSection генерирует секцию для отдельного теста
func (g *Generator) generateTestSection(sb *strings.Builder, test types.TestInfo) {
	sb.WriteString(fmt.Sprintf("### %s\n\n", test.Name))

	// Базовая информация
	sb.WriteString("| Параметр | Значение |\n")
	sb.WriteString("|----------|----------|\n")
	sb.WriteString(fmt.Sprintf("| **Тип** | %s |\n", g.getTestTypeDisplayName(test.Type)))
	sb.WriteString(fmt.Sprintf("| **Пакет** | `%s` |\n", test.Package))
	sb.WriteString(fmt.Sprintf("| **Файл** | `%s:%d` |\n", test.File, test.Line))

	if test.Skipped {
		sb.WriteString("| **Статус** | ⏭️ Пропущен |\n")
		if test.SkipReason != "" {
			sb.WriteString(fmt.Sprintf("| **Причина пропуска** | %s |\n", test.SkipReason))
		}
	} else {
		sb.WriteString("| **Статус** | ✅ Активен |\n")
	}

	if test.Author != "" {
		sb.WriteString(fmt.Sprintf("| **Автор** | %s |\n", test.Author))
	}

	if !test.Created.IsZero() {
		sb.WriteString(fmt.Sprintf("| **Создан** | %s |\n", test.Created.Format("2006-01-02")))
	}

	if !test.Updated.IsZero() {
		sb.WriteString(fmt.Sprintf("| **Обновлен** | %s |\n", test.Updated.Format("2006-01-02")))
	}

	if len(test.Tags) > 0 {
		tags := make([]string, len(test.Tags))
		for i, tag := range test.Tags {
			tags[i] = fmt.Sprintf("`%s`", tag)
		}
		sb.WriteString(fmt.Sprintf("| **Теги** | %s |\n", strings.Join(tags, ", ")))
	}

	sb.WriteString("\n")

	// Описание
	if test.Description != "" {
		sb.WriteString("**Описание:**\n\n")
		sb.WriteString(fmt.Sprintf("%s\n\n", test.Description))
	}

	// Тест-кейсы
	if len(test.TestCases) > 0 {
		sb.WriteString("#### Тест-кейсы\n\n")
		for i, testCase := range test.TestCases {
			sb.WriteString(fmt.Sprintf("**%d. %s**\n\n", i+1, testCase.Name))

			if testCase.Description != "" {
				sb.WriteString(fmt.Sprintf("%s\n\n", testCase.Description))
			}

			if testCase.Input != "" {
				sb.WriteString(fmt.Sprintf("- **Входные данные:** %s\n", testCase.Input))
			}

			if testCase.Expected != "" {
				sb.WriteString(fmt.Sprintf("- **Ожидаемый результат:** %s\n", testCase.Expected))
			}

			if len(testCase.Steps) > 0 {
				sb.WriteString("\n**Шаги:**\n\n")
				for j, step := range testCase.Steps {
					sb.WriteString(fmt.Sprintf("%d. %s", j+1, step.Action))
					if step.Expected != "" {
						sb.WriteString(fmt.Sprintf(" → %s", step.Expected))
					}
					sb.WriteString("\n")
				}
			}

			sb.WriteString("\n")
		}
	}

	// Дополнительные метаданные
	if len(test.Metadata) > 0 {
		sb.WriteString("#### Дополнительная информация\n\n")
		for key, value := range test.Metadata {
			sb.WriteString(fmt.Sprintf("- **%s:** %s\n", strings.Title(key), value))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
}

// groupTestsByType группирует тесты по типам
func (g *Generator) groupTestsByType(packages map[string]*types.PackageInfo) map[types.TestType][]types.TestInfo {
	typeGroups := make(map[types.TestType][]types.TestInfo)

	for _, pkg := range packages {
		for _, test := range pkg.Tests {
			typeGroups[test.Type] = append(typeGroups[test.Type], test)
		}
	}

	// Сортируем тесты в каждой группе
	for testType := range typeGroups {
		sort.Slice(typeGroups[testType], func(i, j int) bool {
			return typeGroups[testType][i].Name < typeGroups[testType][j].Name
		})
	}

	return typeGroups
}

// getTestTypeDisplayName возвращает отображаемое имя типа теста
func (g *Generator) getTestTypeDisplayName(testType types.TestType) string {
	switch testType {
	case types.UnitTest:
		return "Модульные"
	case types.IntegrationTest:
		return "Интеграционные"
	case types.FunctionalTest:
		return "Функциональные"
	case types.E2ETest:
		return "E2E"
	case types.PerformanceTest:
		return "Производительности"
	case types.SecurityTest:
		return "Безопасности"
	case types.RegressionTest:
		return "Регрессионные"
	case types.SmokeTest:
		return "Дымовые"
	default:
		return string(testType)
	}
}
