// Package main демонстрирует продвинутое использование testdoc библиотеки
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/testdoc-org/testdoc"
	"github.com/testdoc-org/testdoc/pkg/types"
)

func main() {
	// Загружаем конфигурацию из файла
	config, err := testdoc.LoadConfig("../../config.yaml")
	if err != nil {
		// Если конфиг не найден, используем по умолчанию
		config = testdoc.DefaultConfig()
		config.Title = "Продвинутый пример документации тестов"
		config.Author = "Демонстрация фильтрации"
	}

	// Парсим тесты
	result, err := testdoc.ParseDirectory("../../examples/_examples", config)
	if err != nil {
		log.Fatalf("Ошибка парсинга тестов: %v", err)
	}

	// Показываем общую статистику
	fmt.Println("📊 Общая статистика:")
	printStats(result)

	// Создаем фильтры
	filter := testdoc.NewFilter()
	stats := testdoc.NewStatistics()

	// Фильтруем только unit тесты
	unitTests := filter.ByType(result, types.UnitTest)
	fmt.Println("\n🔧 Unit тесты:")
	printStats(unitTests)

	// Фильтруем только интеграционные тесты
	integrationTests := filter.ByType(result, types.IntegrationTest)
	fmt.Println("\n🔗 Интеграционные тесты:")
	printStats(integrationTests)

	// Фильтруем тесты по тегам
	paymentTests := filter.ByTags(result, []string{"payment"})
	fmt.Println("\n💳 Тесты платежей:")
	printStats(paymentTests)

	// Показываем наиболее часто встречающийся тип тестов
	mostCommonType, count := stats.GetMostCommonTestType(result)
	fmt.Printf("\n🏆 Наиболее частый тип тестов: %s (%d тестов)\n", mostCommonType, count)

	// Показываем покрытие по типам
	coverage := stats.CalculateTestCoverage(result)
	fmt.Println("\n📈 Покрытие по типам тестов:")
	for testType, percent := range coverage {
		fmt.Printf("   - %s: %.1f%%\n", testType, percent)
	}

	// Генерируем документацию для всех тестов
	allDocsMarkdown := testdoc.GenerateMarkdown(result, config)
	err = testdoc.WriteToFile(allDocsMarkdown, "all-tests-documentation.md")
	if err != nil {
		log.Fatalf("Ошибка записи файла: %v", err)
	}

	// Генерируем отдельную документацию только для unit тестов
	config.Title = "Документация Unit тестов"
	unitDocsMarkdown := testdoc.GenerateMarkdown(unitTests, config)
	err = testdoc.WriteToFile(unitDocsMarkdown, "unit-tests-documentation.md")
	if err != nil {
		log.Fatalf("Ошибка записи файла: %v", err)
	}

	// Генерируем отдельную документацию для тестов безопасности
	securityTests := filter.ByType(result, types.SecurityTest)
	if securityTests.Stats.TotalTests > 0 {
		config.Title = "Документация тестов безопасности"
		securityDocsMarkdown := testdoc.GenerateMarkdown(securityTests, config)
		err = testdoc.WriteToFile(securityDocsMarkdown, "security-tests-documentation.md")
		if err != nil {
			log.Fatalf("Ошибка записи файла: %v", err)
		}
		fmt.Println("✅ Документация тестов безопасности: security-tests-documentation.md")
	}

	fmt.Println("\n✅ Файлы сгенерированы:")
	fmt.Println("   - all-tests-documentation.md (все тесты)")
	fmt.Println("   - unit-tests-documentation.md (только unit тесты)")

	// Демонстрируем поиск тестов по автору
	authors := getUniqueAuthors(result)
	if len(authors) > 0 {
		fmt.Printf("\n👥 Найденные авторы: %s\n", strings.Join(authors, ", "))

		// Генерируем документацию для первого автора
		firstAuthor := authors[0]
		authorTests := filter.ByAuthor(result, firstAuthor)
		if authorTests.Stats.TotalTests > 0 {
			config.Title = fmt.Sprintf("Тесты автора: %s", firstAuthor)
			authorDocsMarkdown := testdoc.GenerateMarkdown(authorTests, config)
			authorFileName := fmt.Sprintf("tests-by-%s.md", strings.ReplaceAll(strings.ToLower(firstAuthor), " ", "-"))
			err = testdoc.WriteToFile(authorDocsMarkdown, authorFileName)
			if err == nil {
				fmt.Printf("   - %s (тесты автора %s)\n", authorFileName, firstAuthor)
			}
		}
	}
}

func printStats(result *types.ParseResult) {
	fmt.Printf("   - Всего тестов: %d\n", result.Stats.TotalTests)
	fmt.Printf("   - Активных: %d\n", result.Stats.ActiveTests)
	fmt.Printf("   - Пропущенных: %d\n", result.Stats.SkippedTests)
	fmt.Printf("   - Пакетов: %d\n", result.Stats.PackageCount)

	if len(result.Stats.TypeDistribution) > 0 {
		fmt.Printf("   - Типы: ")
		var types []string
		for testType, count := range result.Stats.TypeDistribution {
			types = append(types, fmt.Sprintf("%s(%d)", testType, count))
		}
		fmt.Printf("%s\n", strings.Join(types, ", "))
	}
}

func getUniqueAuthors(result *types.ParseResult) []string {
	authorsMap := make(map[string]bool)

	for _, pkg := range result.Packages {
		for _, test := range pkg.Tests {
			if test.Author != "" {
				authorsMap[test.Author] = true
			}
		}
	}

	var authors []string
	for author := range authorsMap {
		authors = append(authors, author)
	}

	return authors
}
