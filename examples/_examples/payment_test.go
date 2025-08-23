package examples

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// @type: integration
// @author: Анна Петрова
// @created: 2024-01-21
// @tags: payment, integration, external_api
// @testcase: Обработка платежа - тестирует обработку платежа через внешний API
// @testcase: Возврат платежа - тестирует процесс возврата средств
// @step: Создать платежный запрос - запрос должен быть валидным
// @step: Отправить запрос в платежную систему - должен вернуть успешный ответ
// @step: Проверить статус платежа - статус должен быть "завершен"
// @step: Сохранить информацию о платеже - данные должны быть записаны в БД
// TestPaymentProcessing тестирует обработку платежей
func TestPaymentProcessing(t *testing.T) {
	if testing.Short() {
		t.Skip("Пропускаем интеграционный тест в режиме быстрых тестов")
	}

	paymentService := NewPaymentService()

	payment := &Payment{
		Amount:      10000, // 100.00 в копейках
		Currency:    "RUB",
		CardNumber:  "4111111111111111",
		ExpiryMonth: 12,
		ExpiryYear:  2025,
		CVV:         "123",
		Description: "Тестовый платеж",
	}

	// Обработка платежа
	result, err := paymentService.ProcessPayment(payment)
	if err != nil {
		t.Fatalf("Ошибка при обработке платежа: %v", err)
	}
	if result.Status != "completed" {
		t.Errorf("Ожидали статус 'completed', получили '%s'", result.Status)
	}
	if result.TransactionID == "" {
		t.Error("TransactionID не должен быть пустым")
	}

	// Проверка сохранения в БД
	savedPayment, err := paymentService.GetPayment(result.TransactionID)
	if err != nil {
		t.Fatalf("Ошибка при получении сохраненного платежа: %v", err)
	}
	if payment.Amount != savedPayment.Amount {
		t.Errorf("Сумма не совпадает: ожидали %d, получили %d", payment.Amount, savedPayment.Amount)
	}
	if savedPayment.Status != "completed" {
		t.Errorf("Ожидали статус 'completed', получили '%s'", savedPayment.Status)
	}
}

// @type: security
// @author: Сергей Безопасности
// @created: 2024-01-22
// @tags: payment, security, validation
// @testcase: Защита от SQL-инъекций - проверяет защищенность от SQL-инъекций
// @testcase: Валидация номера карты - проверяет корректность номера карты
// @testcase: Маскирование данных карты - проверяет маскирование чувствительных данных
// @step: Попытаться передать SQL-инъекцию в параметрах - должна быть отклонена
// @step: Передать некорректный номер карты - должна вернуться ошибка валидации
// @step: Проверить маскирование номера карты в логах - номер должен быть замаскирован
// TestPaymentSecurity тестирует безопасность платежной системы
func TestPaymentSecurity(t *testing.T) {
	paymentService := NewPaymentService()

	// Тест SQL-инъекции
	maliciousPayment := &Payment{
		Description: "'; DROP TABLE payments; --",
		Amount:      5000,
		Currency:    "RUB",
	}

	_, err := paymentService.ProcessPayment(maliciousPayment)
	if err == nil {
		t.Error("SQL-инъекция должна быть отклонена")
	}

	// Тест валидации номера карты
	invalidCardPayment := &Payment{
		CardNumber: "1234567890123456", // Неверный номер карты
		Amount:     5000,
		Currency:   "RUB",
	}

	_, err = paymentService.ProcessPayment(invalidCardPayment)
	if err == nil {
		t.Error("Неверный номер карты должен быть отклонен")
	}

	// Тест маскирования
	validPayment := &Payment{
		CardNumber: "4111111111111111",
		Amount:     5000,
		Currency:   "RUB",
	}

	logEntry := paymentService.GetLogEntry(validPayment)
	if !strings.Contains(logEntry, "4111********1111") {
		t.Error("Номер карты должен быть замаскирован в логах")
	}
	if strings.Contains(logEntry, "4111111111111111") {
		t.Error("Полный номер карты не должен появляться в логах")
	}
}

// @type: e2e
// @author: Команда QA
// @created: 2024-01-23
// @tags: payment, e2e, user_journey
// @testcase: Полный процесс покупки - тестирует весь путь пользователя от выбора товара до получения чека
// @step: Добавить товар в корзину - товар должен появиться в корзине
// @step: Перейти к оформлению заказа - должна открыться страница оформления
// @step: Заполнить данные карты - форма должна принять данные
// @step: Подтвердить платеж - платеж должен быть обработан
// @step: Получить подтверждение - должен быть отображен чек
// @skip_reason: Требует настройки тестовой среды E2E
// TestFullPurchaseJourney тестирует полный процесс покупки (пропущен)
func TestFullPurchaseJourney(t *testing.T) {
	t.Skip("Требует настройки тестовой среды E2E")

	// Этот тест был бы полным E2E тестом
	// с использованием браузера и реального UI
}

// @type: performance
// @author: Команда Performance
// @created: 2024-01-24
// @tags: payment, performance, load
// @testcase: Нагрузочное тестирование - проверяет производительность при высокой нагрузке
// @performance_threshold: 500ms
// BenchmarkPaymentThroughput тестирует пропускную способность платежной системы
func BenchmarkPaymentThroughput(b *testing.B) {
	paymentService := NewPaymentService()

	payment := &Payment{
		Amount:      1000,
		Currency:    "RUB",
		CardNumber:  "4111111111111111",
		ExpiryMonth: 12,
		ExpiryYear:  2025,
		CVV:         "123",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := paymentService.ProcessPayment(payment)
			if err != nil {
				b.Fatalf("Ошибка при обработке платежа: %v", err)
			}
		}
	})
}

// @type: regression
// @author: Команда QA
// @created: 2024-01-25
// @updated: 2024-01-26
// @tags: payment, regression, bug_fix
// @testcase: Исправление ошибки двойного списания - проверяет исправление бага #1234
// @bug_id: #1234
// @step: Отправить одинаковые платежи подряд - второй должен быть отклонен
// @step: Проверить, что деньги списаны только один раз - баланс должен уменьшиться только на сумму одного платежа
// TestDoubleChargeRegression проверяет исправление ошибки двойного списания
func TestDoubleChargeRegression(t *testing.T) {
	paymentService := NewPaymentService()

	payment := &Payment{
		Amount:         5000,
		Currency:       "RUB",
		CardNumber:     "4111111111111111",
		ExpiryMonth:    12,
		ExpiryYear:     2025,
		CVV:            "123",
		IdempotencyKey: "unique-key-123",
	}

	// Первый платеж должен пройти
	result1, err := paymentService.ProcessPayment(payment)
	if err != nil {
		t.Fatalf("Первый платеж должен быть успешным: %v", err)
	}
	if result1.Status != "completed" {
		t.Errorf("Ожидали статус 'completed', получили '%s'", result1.Status)
	}

	// Второй идентичный платеж должен быть отклонен
	_, err = paymentService.ProcessPayment(payment)
	if err == nil {
		t.Error("Повторный платеж должен быть отклонен")
	}
	if !strings.Contains(err.Error(), "duplicate") {
		t.Error("Ошибка должна содержать информацию о дублировании")
	}
}

// Вспомогательные структуры и функции

type Payment struct {
	Amount         int64     `json:"amount"`
	Currency       string    `json:"currency"`
	CardNumber     string    `json:"card_number"`
	ExpiryMonth    int       `json:"expiry_month"`
	ExpiryYear     int       `json:"expiry_year"`
	CVV            string    `json:"cvv"`
	Description    string    `json:"description"`
	IdempotencyKey string    `json:"idempotency_key"`
	Status         string    `json:"status"`
	TransactionID  string    `json:"transaction_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type PaymentResult struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

type PaymentService struct {
	processedPayments map[string]*Payment
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		processedPayments: make(map[string]*Payment),
	}
}

func (ps *PaymentService) ProcessPayment(payment *Payment) (*PaymentResult, error) {
	// Проверка на дублирование
	if payment.IdempotencyKey != "" {
		if _, exists := ps.processedPayments[payment.IdempotencyKey]; exists {
			return nil, fmt.Errorf("duplicate payment detected")
		}
	}

	// Базовая валидация
	if payment.CardNumber == "1234567890123456" {
		return nil, fmt.Errorf("invalid card number")
	}

	if strings.Contains(payment.Description, "DROP TABLE") {
		return nil, fmt.Errorf("malicious input detected")
	}

	// Симуляция обработки платежа
	transactionID := fmt.Sprintf("txn_%d", time.Now().Unix())

	payment.TransactionID = transactionID
	payment.Status = "completed"
	payment.CreatedAt = time.Now()

	if payment.IdempotencyKey != "" {
		ps.processedPayments[payment.IdempotencyKey] = payment
	}

	return &PaymentResult{
		TransactionID: transactionID,
		Status:        "completed",
		Message:       "Payment processed successfully",
	}, nil
}

func (ps *PaymentService) GetPayment(transactionID string) (*Payment, error) {
	// Поиск платежа по ID (заглушка)
	return &Payment{
		TransactionID: transactionID,
		Amount:        10000,
		Status:        "completed",
	}, nil
}

func (ps *PaymentService) GetLogEntry(payment *Payment) string {
	// Маскирование номера карты для логов
	maskedCard := payment.CardNumber[:4] + "********" + payment.CardNumber[len(payment.CardNumber)-4:]
	return fmt.Sprintf("Processing payment with card: %s", maskedCard)
}
