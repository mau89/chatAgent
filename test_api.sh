#!/bin/bash

# Скрипт для тестирования API

API_URL="http://localhost:8080"

echo "🧪 Тестирование Chat Agent API..."

# Тест 1: Проверка здоровья сервиса
echo "1. Проверка здоровья сервиса..."
curl -s "$API_URL/health" | jq .

echo -e "\n"

# Тест 2: Отправка сообщения
echo "2. Отправка сообщения агенту..."
curl -s -X POST "$API_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Привет! Как дела?",
    "user_id": 12345
  }' | jq .

echo -e "\n"

# Тест 3: Запрос времени
echo "3. Запрос времени..."
curl -s -X POST "$API_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Сколько времени?",
    "user_id": 12345
  }' | jq .

echo -e "\n"

# Тест 4: Запрос погоды
echo "4. Запрос погоды..."
curl -s -X POST "$API_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Какая погода?",
    "user_id": 12345
  }' | jq .

echo -e "\n"

# Тест 5: Запрос помощи
echo "5. Запрос помощи..."
curl -s -X POST "$API_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Помощь",
    "user_id": 12345
  }' | jq .

echo -e "\n✅ Тестирование завершено!"

