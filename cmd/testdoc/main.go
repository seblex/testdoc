// Package main предоставляет CLI для генерации документации тестов
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/seblex/testdoc"
	"github.com/seblex/testdoc/pkg/types"
)

const version = "1.0.0"

func main() {
	var (
		outputFile   = flag.String("output", "test-documentation.md", "Файл для вывода документации")
		configFile   = flag.String("config", "", "Файл конфигурации YAML (опционально)")
		showVersion  = flag.Bool("version", false, "Показать версию")
		showHelp     = flag.Bool("help", false, "Показать справку")
		filterType   = flag.String("type", "", "Фильтр по типу тестов (unit, integration, functional, e2e, performance, security, regression, smoke)")
		filterAuthor = flag.String("author", "", "Фильтр по автору")
		filterTags   = flag.String("tags", "", "Фильтр по тегам (через запятую)")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "TestDoc v%s - Генератор документации для Go тестов\n\n", version)
		fmt.Fprintf(os.Stderr, "Использование: %s [опции] [путь]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "ВАЖНО: Все опции должны указываться ДО пути к директории!\n\n")
		fmt.Fprintf(os.Stderr, "Опции:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nПримеры:\n")
		fmt.Fprintf(os.Stderr, "  %s                                    # Анализ текущей директории\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s ./examples                         # Анализ директории examples\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -output docs.md ./tests            # С указанием выходного файла\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -config config.yaml                # С конфигурацией\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -type unit -author 'John Doe'      # С фильтрами\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -tags 'api,database' ./internal    # Фильтр по тегам\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nПоддерживаемые типы тестов:\n")
		for _, testType := range testdoc.GetSupportedTestTypes() {
			fmt.Fprintf(os.Stderr, "  - %s\n", testType)
		}
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("TestDoc v%s\n", version)
		os.Exit(0)
	}

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	path := "."
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	// Загружаем конфигурацию
	var config *types.Config
	var err error

	if *configFile != "" {
		config, err = testdoc.LoadConfig(*configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка загрузки конфигурации: %v\n", err)
			os.Exit(1)
		}
	} else {
		config = testdoc.DefaultConfig()
	}

	// Валидируем конфигурацию
	err = testdoc.ValidateConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка валидации конфигурации: %v\n", err)
		os.Exit(1)
	}

	// Парсим тесты
	result, err := testdoc.ParseDirectory(path, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка анализа тестов: %v\n", err)
		os.Exit(1)
	}

	// Применяем фильтры
	if *filterType != "" {
		testType := types.TestType(*filterType)
		if !testdoc.IsValidTestType(testType) {
			fmt.Fprintf(os.Stderr, "Неизвестный тип теста: %s\n", *filterType)
			fmt.Fprintf(os.Stderr, "Поддерживаемые типы: %v\n", testdoc.GetSupportedTestTypes())
			os.Exit(1)
		}
		filter := testdoc.NewFilter()
		result = filter.ByType(result, testType)
	}

	if *filterAuthor != "" {
		filter := testdoc.NewFilter()
		result = filter.ByAuthor(result, *filterAuthor)
	}

	if *filterTags != "" {
		tags := strings.Split(*filterTags, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		filter := testdoc.NewFilter()
		result = filter.ByTags(result, tags)
	}

	// Проверяем, что найдены тесты
	if result.Stats.TotalTests == 0 {
		fmt.Fprintf(os.Stderr, "Не найдено тестов в директории: %s\n", path)
		if *filterType != "" || *filterAuthor != "" || *filterTags != "" {
			fmt.Fprintf(os.Stderr, "Попробуйте изменить фильтры или проверить директорию.\n")
		}
		os.Exit(1)
	}

	// Генерируем документацию
	markdown := testdoc.GenerateMarkdown(result, config)

	// Сохраняем в файл
	err = testdoc.WriteToFile(markdown, *outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка записи файла: %v\n", err)
		os.Exit(1)
	}

	// Выводим статистику
	fmt.Printf("✅ Документация успешно сгенерирована: %s\n", *outputFile)
	fmt.Printf("📊 Статистика:\n")
	fmt.Printf("   - Всего тестов: %d\n", result.Stats.TotalTests)
	fmt.Printf("   - Активных: %d\n", result.Stats.ActiveTests)
	fmt.Printf("   - Пропущенных: %d\n", result.Stats.SkippedTests)
	fmt.Printf("   - Пакетов: %d\n", result.Stats.PackageCount)

	if len(result.Stats.TypeDistribution) > 0 {
		fmt.Printf("   - Распределение по типам:\n")
		for testType, count := range result.Stats.TypeDistribution {
			percentage := float64(count) / float64(result.Stats.TotalTests) * 100
			fmt.Printf("     * %s: %d (%.1f%%)\n", testType, count, percentage)
		}
	}

	// Успешное завершение
	os.Exit(0)
}
