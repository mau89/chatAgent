package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramBot представляет Telegram бота
type TelegramBot struct {
	bot        *tgbotapi.BotAPI
	httpClient *HTTPClient
}

// NewTelegramBot создает новый экземпляр Telegram бота
func NewTelegramBot(token string) *TelegramBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Ошибка создания бота:", err)
	}

	// Получаем URL для внешнего API
	apiURL := os.Getenv("EXTERNAL_API_URL")
	httpClient := NewHTTPClient(apiURL)

	return &TelegramBot{
		bot:        bot,
		httpClient: httpClient,
	}
}

// Start запускает бота
func (tb *TelegramBot) Start() error {
	log.Printf("Бот @%s запущен", tb.bot.Self.UserName)

	// Настраиваем обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tb.bot.GetUpdatesChan(u)

	// Обрабатываем обновления
	for update := range updates {
		if update.Message != nil {
			go tb.handleMessage(update.Message)
		}
	}

	return nil
}

// handleMessage обрабатывает входящие сообщения
func (tb *TelegramBot) handleMessage(message *tgbotapi.Message) {
	log.Printf("Получено сообщение от %s (%d): %s", 
		message.From.UserName, message.From.ID, message.Text)

	// Обрабатываем команды
	if message.IsCommand() {
		tb.handleCommand(message)
		return
	}

	// Обрабатываем обычные сообщения
	tb.handleTextMessage(message)
}

// handleCommand обрабатывает команды
func (tb *TelegramBot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		tb.sendMessage(message.Chat.ID, 
			"🤖 Привет! Я ваш персональный агент-помощник.\n\n"+
			"Я могу помочь с:\n"+
			"• Информацией о погоде\n"+
			"• Текущим временем\n"+
			"• Математическими вычислениями\n"+
			"• И многим другим!\n\n"+
			"Просто напишите мне вопрос или используйте /help для списка команд.")
		
	case "help":
		response, _ := tb.httpClient.SendRequest("/help", message.From.ID)
		tb.sendMessage(message.Chat.ID, response.Answer)
		
	case "weather":
		response, _ := tb.httpClient.SendRequest("погода", message.From.ID)
		tb.sendMessage(message.Chat.ID, response.Answer)
		
	case "time":
		response, _ := tb.httpClient.SendRequest("время", message.From.ID)
		tb.sendMessage(message.Chat.ID, response.Answer)
		
	case "calculate":
		tb.sendMessage(message.Chat.ID, 
			"🧮 Для вычислений напишите выражение, например:\n"+
			"• 2 + 3\n"+
			"• 10 - 5\n"+
			"• 4 * 6\n"+
			"• 15 / 3")
		
	default:
		tb.sendMessage(message.Chat.ID, 
			"❓ Неизвестная команда. Используйте /help для списка доступных команд.")
	}
}

// handleTextMessage обрабатывает текстовые сообщения
func (tb *TelegramBot) handleTextMessage(message *tgbotapi.Message) {
	// Показываем, что бот печатает
	tb.sendTypingAction(message.Chat.ID)

	// Отправляем запрос через HTTP клиент
	response, err := tb.httpClient.SendRequest(message.Text, message.From.ID)
	if err != nil {
		log.Printf("Ошибка HTTP запроса: %v", err)
		tb.sendMessage(message.Chat.ID, 
			"❌ Извините, произошла ошибка при обработке вашего запроса.")
		return
	}

	// Отправляем ответ пользователю
	tb.sendMessage(message.Chat.ID, response.Answer)
}

// sendMessage отправляет сообщение пользователю
func (tb *TelegramBot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	
	if _, err := tb.bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

// sendTypingAction показывает индикатор "печатает"
func (tb *TelegramBot) sendTypingAction(chatID int64) {
	action := tgbotapi.NewChatAction(chatID, "typing")
	tb.bot.Send(action)
}

// sendMessageWithKeyboard отправляет сообщение с клавиатурой
func (tb *TelegramBot) sendMessageWithKeyboard(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	
	// Создаем клавиатуру с быстрыми командами
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🌤️ Погода"),
			tgbotapi.NewKeyboardButton("🕐 Время"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🧮 Вычисления"),
			tgbotapi.NewKeyboardButton("❓ Помощь"),
		),
	)
	keyboard.OneTimeKeyboard = true
	msg.ReplyMarkup = keyboard
	
	if _, err := tb.bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения с клавиатурой: %v", err)
	}
}

