// Package main демонстрирует базовое использование testdoc библиотеки
package main

import (
	"fmt"
	"log"

	"github.com/seblex/testdoc"
)

func main() {
	// Создаем конфигурацию
	config := testdoc.DefaultConfig()
	config.Title = "Документация тестов примера"
	config.Author = "Демонстрационный пример"
	config.GroupByType = true

	// Генерируем документацию для примеров
	doc, err := testdoc.GenerateFromDirectory("../../examples/_examples", config)
	if err != nil {
		log.Fatalf("Ошибка генерации документации: %v", err)
	}

	// Сохраняем в файл
	outputFile := "example-documentation.md"
	err = testdoc.WriteToFile(doc, outputFile)
	if err != nil {
		log.Fatalf("Ошибка записи файла: %v", err)
	}

	fmt.Printf("✅ Документация сгенерирована: %s\n", outputFile)
}
