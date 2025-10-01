package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// Agent представляет интеллектуального агента
type Agent struct {
	conversationHistory map[int64][]ConversationEntry
	tools               map[string]Tool
	yandexGPT           *YandexGPTClient
	useYandexGPT        bool
}

// ConversationEntry представляет запись в истории разговора
type ConversationEntry struct {
	Message   string
	Response  string
	Timestamp time.Time
}

// Tool представляет инструмент, который может использовать агент
type Tool struct {
	Name        string
	Description string
	Handler     func(string, int64) (string, error)
}

// NewAgent создает новый экземпляр агента
func NewAgent() *Agent {
	agent := &Agent{
		conversationHistory: make(map[int64][]ConversationEntry),
		tools:               make(map[string]Tool),
		useYandexGPT:        false,
	}

	// Регистрируем доступные инструменты
	agent.registerTools()

	return agent
}

// NewAgentWithYandexGPT создает агента с поддержкой Yandex GPT
func NewAgentWithYandexGPT(apiKey, folderID string) *Agent {
	agent := &Agent{
		conversationHistory: make(map[int64][]ConversationEntry),
		tools:               make(map[string]Tool),
		yandexGPT:           NewYandexGPTClient(apiKey, folderID),
		useYandexGPT:        true,
	}

	// Регистрируем доступные инструменты
	agent.registerTools()

	return agent
}

// registerTools регистрирует доступные инструменты агента
func (a *Agent) registerTools() {
	a.tools["weather"] = Tool{
		Name:        "weather",
		Description: "Получить информацию о погоде",
		Handler:     a.handleWeatherRequest,
	}

	a.tools["time"] = Tool{
		Name:        "time",
		Description: "Получить текущее время",
		Handler:     a.handleTimeRequest,
	}

	a.tools["calculate"] = Tool{
		Name:        "calculate",
		Description: "Выполнить математические вычисления",
		Handler:     a.handleCalculateRequest,
	}

	a.tools["help"] = Tool{
		Name:        "help",
		Description: "Показать доступные команды",
		Handler:     a.handleHelpRequest,
	}
}

// ProcessMessage обрабатывает входящее сообщение
func (a *Agent) ProcessMessage(message string, userID int64) (string, error) {
	log.Printf("Обработка сообщения от пользователя %d: %s", userID, message)

	// Добавляем сообщение в историю
	a.addToHistory(userID, message, "")

	// Определяем, какой инструмент использовать
	toolName := a.determineTool(message)
	
	var response string
	var err error

	// Если включен Yandex GPT и это не специальная команда, используем его
	if a.useYandexGPT && a.yandexGPT != nil && toolName == "general" {
		response, err = a.yandexGPT.GenerateResponse(message, userID)
		if err != nil {
			log.Printf("Ошибка Yandex GPT, переключаемся на встроенные инструменты: %v", err)
			// Fallback на встроенные инструменты
			response = a.generateGeneralResponse(message)
			err = nil // Сбрасываем ошибку, так как мы обработали её
		}
	} else if tool, exists := a.tools[toolName]; exists {
		response, err = tool.Handler(message, userID)
	} else {
		// Если инструмент не найден, используем общий ответ
		response = a.generateGeneralResponse(message)
	}

	if err != nil {
		log.Printf("Ошибка при обработке сообщения: %v", err)
		return "Извините, произошла ошибка при обработке вашего запроса.", nil
	}

	// Обновляем историю с ответом
	a.updateLastResponse(userID, response)

	return response, nil
}

// determineTool определяет, какой инструмент использовать на основе сообщения
func (a *Agent) determineTool(message string) string {
	message = strings.ToLower(message)

	// Ключевые слова для определения инструментов
	if strings.Contains(message, "погода") || strings.Contains(message, "weather") {
		return "weather"
	}
	if strings.Contains(message, "время") || strings.Contains(message, "time") || 
	   strings.Contains(message, "сколько времени") {
		return "time"
	}
	if strings.Contains(message, "вычисли") || strings.Contains(message, "calculate") || 
	   strings.Contains(message, "сложи") || strings.Contains(message, "умножь") {
		return "calculate"
	}
	if strings.Contains(message, "помощь") || strings.Contains(message, "help") || 
	   strings.Contains(message, "команды") || strings.HasPrefix(message, "/help") {
		return "help"
	}

	return "general"
}

// generateGeneralResponse генерирует общий ответ
func (a *Agent) generateGeneralResponse(message string) string {
	responses := []string{
		"Интересный вопрос! Можете уточнить, что именно вас интересует?",
		"Я понимаю ваш вопрос. Могу помочь с информацией о погоде, времени, вычислениями или другими задачами.",
		"Хорошо! Для более точного ответа используйте команды: /help - список команд, /weather - погода, /time - время.",
		"Я готов помочь! Попробуйте спросить о погоде, времени или попросите выполнить вычисления.",
	}

	// Простая логика выбора ответа на основе длины сообщения
	index := len(message) % len(responses)
	return responses[index]
}

// addToHistory добавляет сообщение в историю разговора
func (a *Agent) addToHistory(userID int64, message, response string) {
	if a.conversationHistory[userID] == nil {
		a.conversationHistory[userID] = make([]ConversationEntry, 0)
	}

	a.conversationHistory[userID] = append(a.conversationHistory[userID], ConversationEntry{
		Message:   message,
		Response:  response,
		Timestamp: time.Now(),
	})

	// Ограничиваем историю последними 10 сообщениями
	if len(a.conversationHistory[userID]) > 10 {
		a.conversationHistory[userID] = a.conversationHistory[userID][len(a.conversationHistory[userID])-10:]
	}
}

// updateLastResponse обновляет последний ответ в истории
func (a *Agent) updateLastResponse(userID int64, response string) {
	if history, exists := a.conversationHistory[userID]; exists && len(history) > 0 {
		history[len(history)-1].Response = response
	}
}

// Обработчики инструментов

func (a *Agent) handleWeatherRequest(message string, userID int64) (string, error) {
	return "🌤️ К сожалению, я пока не подключен к сервису погоды. Но могу сказать, что сегодня отличный день для прогулки!", nil
}

func (a *Agent) handleTimeRequest(message string, userID int64) (string, error) {
	currentTime := time.Now().Format("15:04:05")
	currentDate := time.Now().Format("02.01.2006")
	return fmt.Sprintf("🕐 Текущее время: %s\n📅 Дата: %s", currentTime, currentDate), nil
}

func (a *Agent) handleCalculateRequest(message string, userID int64) (string, error) {
	// Простой калькулятор для базовых операций
	message = strings.ToLower(message)
	
	// Ищем числа и операции
	if strings.Contains(message, "+") {
		return "➕ Для сложения используйте формат: '2 + 3'. Я пока учусь математике! 😊", nil
	}
	if strings.Contains(message, "-") {
		return "➖ Для вычитания используйте формат: '5 - 2'. Я пока учусь математике! 😊", nil
	}
	if strings.Contains(message, "*") || strings.Contains(message, "×") {
		return "✖️ Для умножения используйте формат: '3 * 4'. Я пока учусь математике! 😊", nil
	}
	if strings.Contains(message, "/") || strings.Contains(message, "÷") {
		return "➗ Для деления используйте формат: '8 / 2'. Я пока учусь математике! 😊", nil
	}
	
	return "🧮 Для вычислений используйте команду /calculate или напишите 'вычисли 2+2'", nil
}

func (a *Agent) handleHelpRequest(message string, userID int64) (string, error) {
	helpText := `🤖 *Доступные команды:*

/start - Начать работу с ботом
/help - Показать это сообщение
/weather - Информация о погоде
/time - Текущее время
/calculate - Математические вычисления

*Примеры вопросов:*
• "Какая погода?"
• "Сколько времени?"
• "Вычисли 2+2"
• "Помощь"

Я готов помочь вам! 😊`

	return helpText, nil
}

