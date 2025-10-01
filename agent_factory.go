package main

import (
	"log"
	"os"
)

// createAgent создает агента на основе конфигурации
func createAgent() *Agent {
	// Проверяем, нужно ли использовать Yandex GPT
	useYandexGPT := os.Getenv("USE_YANDEX_GPT")
	if useYandexGPT == "true" {
		apiKey := os.Getenv("YANDEX_GPT_API_KEY")
		folderID := os.Getenv("YANDEX_GPT_FOLDER_ID")
		
		if apiKey != "" && folderID != "" {
			log.Printf("Создаем агента с Yandex GPT")
			return NewAgentWithYandexGPT(apiKey, folderID)
		} else {
			log.Printf("Yandex GPT включен, но API ключ или Folder ID не установлены. Используем встроенного агента.")
		}
	}
	
	log.Printf("Создаем встроенного агента")
	return NewAgent()
}
