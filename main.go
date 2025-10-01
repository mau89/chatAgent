package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем переменные окружения")
	}

	// Получаем порт для HTTP сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Запускаем HTTP сервер в отдельной горутине
	go func() {
		httpServer := NewHTTPServer(port)
		if err := httpServer.Start(); err != nil {
			log.Printf("Ошибка запуска HTTP сервера: %v", err)
		}
	}()

	// Получаем токен бота
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не установлен")
	}

	// Проверяем настройки Yandex GPT
	useYandexGPT := os.Getenv("USE_YANDEX_GPT")
	if useYandexGPT == "true" {
		apiKey := os.Getenv("YANDEX_GPT_API_KEY")
		folderID := os.Getenv("YANDEX_GPT_FOLDER_ID")
		
		if apiKey == "" || folderID == "" {
			log.Println("⚠️  Yandex GPT включен, но API ключ или Folder ID не установлены")
			log.Println("   Установите YANDEX_GPT_API_KEY и YANDEX_GPT_FOLDER_ID в .env файле")
			log.Println("   Или установите USE_YANDEX_GPT=false для использования встроенного агента")
		} else {
			log.Println("✅ Yandex GPT настроен и готов к работе")
		}
	} else {
		log.Println("ℹ️  Используется встроенный агент (Yandex GPT отключен)")
	}

	// Создаем и запускаем бота
	bot := NewTelegramBot(botToken)
	if err := bot.Start(); err != nil {
		log.Fatal("Ошибка запуска бота:", err)
	}
}
