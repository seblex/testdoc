// Package parser предоставляет функциональность для анализа Go тест-файлов
package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/seblex/testdoc/pkg/types"
)

// Parser анализирует Go файлы с тестами
type Parser struct {
	fileSet *token.FileSet
}

// New создает новый парсер тестов
func New() *Parser {
	return &Parser{
		fileSet: token.NewFileSet(),
	}
}

// ParseFile анализирует один тест-файл и возвращает информацию о тестах
func (p *Parser) ParseFile(filename string) ([]types.TestInfo, error) {
	src, err := parser.ParseFile(p.fileSet, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var tests []types.TestInfo

	for _, decl := range src.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if p.isTestFunction(fn.Name.Name) {
				testInfo := p.parseTestFunction(fn, src, filename)
				tests = append(tests, testInfo)
			}
		}
	}

	return tests, nil
}

// ParseDirectory рекурсивно анализирует директорию и возвращает результат парсинга
func (p *Parser) ParseDirectory(rootPath string, config *types.Config) (*types.ParseResult, error) {
	packages := make(map[string]*types.PackageInfo)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !p.isTestFile(path) {
			return nil
		}

		// Проверяем паттерны исключения и включения
		if p.shouldExcludeFile(path, config) {
			return nil
		}

		tests, err := p.ParseFile(path)
		if err != nil {
			// Логируем предупреждение, но продолжаем
			return nil
		}

		for _, test := range tests {
			// Применяем значения по умолчанию
			if test.Type == "" {
				test.Type = types.UnitTest
			}

			// Пропускаем пропущенные тесты, если настроено
			if test.Skipped && !config.IncludeSkipped {
				continue
			}

			packageKey := test.Package
			if packages[packageKey] == nil {
				packages[packageKey] = &types.PackageInfo{
					Name:      test.Package,
					Path:      filepath.Dir(path),
					Tests:     []types.TestInfo{},
					TestTypes: []types.TestType{},
				}
			}

			packages[packageKey].Tests = append(packages[packageKey].Tests, test)

			// Добавляем тип теста в список типов пакета
			found := false
			for _, t := range packages[packageKey].TestTypes {
				if t == test.Type {
					found = true
					break
				}
			}
			if !found {
				packages[packageKey].TestTypes = append(packages[packageKey].TestTypes, test.Type)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	result := &types.ParseResult{
		Packages: packages,
	}
	result.CalculateStats()

	return result, nil
}

// isTestFunction проверяет, является ли функция тест-функцией
func (p *Parser) isTestFunction(name string) bool {
	return strings.HasPrefix(name, "Test") ||
		strings.HasPrefix(name, "Benchmark") ||
		strings.HasPrefix(name, "Example")
}

// isTestFile проверяет, является ли файл тест-файлом
func (p *Parser) isTestFile(filename string) bool {
	return strings.HasSuffix(filename, "_test.go")
}

// shouldExcludeFile проверяет, нужно ли исключить файл
func (p *Parser) shouldExcludeFile(path string, config *types.Config) bool {
	// Проверяем паттерны включения
	if len(config.IncludePatterns) > 0 {
		included := false
		for _, pattern := range config.IncludePatterns {
			if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
				included = true
				break
			}
		}
		if !included {
			return true
		}
	}

	// Проверяем паттерны исключения
	for _, pattern := range config.ExcludePatterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
	}

	return false
}

// parseTestFunction извлекает информацию о тест-функции
func (p *Parser) parseTestFunction(fn *ast.FuncDecl, file *ast.File, filename string) types.TestInfo {
	position := p.fileSet.Position(fn.Pos())

	testInfo := types.TestInfo{
		Name:     fn.Name.Name,
		File:     filepath.Base(filename),
		Line:     position.Line,
		Package:  file.Name.Name,
		Tags:     []string{},
		Metadata: make(map[string]string),
	}

	// Анализируем комментарии функции
	if fn.Doc != nil {
		p.parseDocComments(fn.Doc, &testInfo)
	}

	// Анализируем тело функции для поиска skip-ов
	if fn.Body != nil {
		p.analyzeTestBody(fn.Body, &testInfo)
	}

	return testInfo
}

// parseDocComments анализирует doc-комментарии функции
func (p *Parser) parseDocComments(docGroup *ast.CommentGroup, testInfo *types.TestInfo) {
	fullComment := ""

	for _, comment := range docGroup.List {
		line := strings.TrimPrefix(comment.Text, "//")
		line = strings.TrimPrefix(line, "/*")
		line = strings.TrimSuffix(line, "*/")
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		// Проверяем на аннотации
		if strings.HasPrefix(line, "@") {
			p.parseAnnotation(line, testInfo)
		} else {
			if fullComment != "" {
				fullComment += " "
			}
			fullComment += line
		}
	}

	testInfo.Description = fullComment
}

// parseAnnotation парсит аннотации в комментариях
func (p *Parser) parseAnnotation(line string, testInfo *types.TestInfo) {
	// @type: unit|integration|functional|e2e|performance|security|regression|smoke
	if strings.HasPrefix(line, "@type:") {
		typeStr := strings.TrimSpace(strings.TrimPrefix(line, "@type:"))
		testInfo.Type = types.TestType(typeStr)
		return
	}

	// @author: имя автора
	if strings.HasPrefix(line, "@author:") {
		testInfo.Author = strings.TrimSpace(strings.TrimPrefix(line, "@author:"))
		return
	}

	// @tags: tag1,tag2,tag3
	if strings.HasPrefix(line, "@tags:") {
		tagsStr := strings.TrimSpace(strings.TrimPrefix(line, "@tags:"))
		tags := strings.Split(tagsStr, ",")
		for _, tag := range tags {
			testInfo.Tags = append(testInfo.Tags, strings.TrimSpace(tag))
		}
		return
	}

	// @testcase: название кейса - описание
	if strings.HasPrefix(line, "@testcase:") {
		caseStr := strings.TrimSpace(strings.TrimPrefix(line, "@testcase:"))
		parts := strings.SplitN(caseStr, "-", 2)

		testCase := types.TestCase{
			Name: strings.TrimSpace(parts[0]),
		}

		if len(parts) > 1 {
			testCase.Description = strings.TrimSpace(parts[1])
		}

		testInfo.TestCases = append(testInfo.TestCases, testCase)
		return
	}

	// @step: действие - ожидаемый результат
	if strings.HasPrefix(line, "@step:") {
		stepStr := strings.TrimSpace(strings.TrimPrefix(line, "@step:"))
		parts := strings.SplitN(stepStr, "-", 2)

		step := types.Step{
			Action: strings.TrimSpace(parts[0]),
		}

		if len(parts) > 1 {
			step.Expected = strings.TrimSpace(parts[1])
		}

		// Добавляем к последнему тест-кейсу
		if len(testInfo.TestCases) > 0 {
			lastIndex := len(testInfo.TestCases) - 1
			testInfo.TestCases[lastIndex].Steps = append(testInfo.TestCases[lastIndex].Steps, step)
		}
		return
	}

	// @created: дата создания
	if strings.HasPrefix(line, "@created:") {
		dateStr := strings.TrimSpace(strings.TrimPrefix(line, "@created:"))
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			testInfo.Created = date
		}
		return
	}

	// @updated: дата обновления
	if strings.HasPrefix(line, "@updated:") {
		dateStr := strings.TrimSpace(strings.TrimPrefix(line, "@updated:"))
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			testInfo.Updated = date
		}
		return
	}

	// Произвольные метаданные @key: value
	if strings.Contains(line, ":") {
		parts := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(strings.TrimPrefix(parts[0], "@"))
		value := strings.TrimSpace(parts[1])
		testInfo.Metadata[key] = value
	}
}

// analyzeTestBody анализирует тело тест-функции
func (p *Parser) analyzeTestBody(body *ast.BlockStmt, testInfo *types.TestInfo) {
	ast.Inspect(body, func(n ast.Node) bool {
		// Ищем вызовы t.Skip(), t.Skipf(), t.SkipNow()
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "Skip" || sel.Sel.Name == "Skipf" || sel.Sel.Name == "SkipNow" {
					testInfo.Skipped = true

					// Пытаемся извлечь причину пропуска из аргументов
					if len(call.Args) > 0 {
						if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
							reason, _ := strconv.Unquote(lit.Value)
							testInfo.SkipReason = reason
						}
					}
				}
			}
		}
		return true
	})
}
