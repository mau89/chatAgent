#!/bin/bash

# Скрипт для тестирования Yandex GPT

echo "🧪 Тестирование Yandex GPT..."

# Проверяем наличие .env файла
if [ ! -f .env ]; then
    echo "❌ Файл .env не найден"
    echo "📝 Скопируйте config.env.example в .env и настройте ключи"
    exit 1
fi

# Загружаем переменные окружения
source .env

# Проверяем наличие ключей
if [ "$YANDEX_GPT_API_KEY" = "your_yandex_api_key_here" ] || [ -z "$YANDEX_GPT_API_KEY" ]; then
    echo "❌ YANDEX_GPT_API_KEY не настроен"
    echo "📖 См. инструкцию: YANDEX_SETUP.md"
    exit 1
fi

if [ "$YANDEX_GPT_FOLDER_ID" = "your_folder_id_here" ] || [ -z "$YANDEX_GPT_FOLDER_ID" ]; then
    echo "❌ YANDEX_GPT_FOLDER_ID не настроен"
    echo "📖 См. инструкцию: YANDEX_SETUP.md"
    exit 1
fi

echo "✅ Ключи найдены"
echo "🔑 API Key: ${YANDEX_GPT_API_KEY:0:20}..."
echo "📁 Folder ID: $YANDEX_GPT_FOLDER_ID"

# Тестируем API
echo ""
echo "🚀 Тестируем API..."

# Запускаем бота в фоне
go run . &
BOT_PID=$!

# Ждем запуска
sleep 3

# Тестовые запросы
echo ""
echo "📝 Тестовые запросы:"

echo "1. Общий вопрос:"
curl -s -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Привет! Как дела?", "user_id": 12345}' | jq -r '.answer'

echo ""
echo "2. Вопрос о погоде:"
curl -s -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Какая погода?", "user_id": 12345}' | jq -r '.answer'

echo ""
echo "3. Время:"
curl -s -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Сколько времени?", "user_id": 12345}' | jq -r '.answer'

echo ""
echo "4. Помощь:"
curl -s -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Помощь", "user_id": 12345}' | jq -r '.answer'

# Останавливаем бота
kill $BOT_PID 2>/dev/null

echo ""
echo "✅ Тестирование завершено!"
echo "📖 Если видите fallback ответы, проверьте настройку ключей в YANDEX_SETUP.md"
