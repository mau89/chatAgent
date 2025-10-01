#!/bin/bash

# Скрипт для запуска Telegram Chat Agent

echo "🤖 Запуск Telegram Chat Agent..."

# Проверяем наличие .env файла
if [ ! -f .env ]; then
    echo "⚠️  Файл .env не найден."
    echo "📝 Выберите вариант настройки:"
    echo "   1. Быстрый старт с DeepSeek (бесплатно!)"
    echo "   2. Встроенный агент (без внешних API)"
    echo "   3. Полная настройка"
    echo ""
    read -p "Введите номер (1-3): " choice
    
    case $choice in
        1)
            echo "🚀 Копируем конфигурацию для быстрого старта с DeepSeek..."
            cp quick_start.env .env
            echo "📝 Осталось добавить:"
            echo "   - TELEGRAM_BOT_TOKEN (получить у @BotFather)"
            echo "   - DEEPSEEK_API_KEY (получить на https://platform.deepseek.com/)"
            ;;
        2)
            echo "🔧 Копируем конфигурацию для встроенного агента..."
            cp config.env.example .env
            echo "📝 Добавьте TELEGRAM_BOT_TOKEN в .env файл"
            echo "   Установите USE_DEEPSEEK=false"
            ;;
        3)
            echo "⚙️  Копируем полную конфигурацию..."
            cp config.env.example .env
            echo "📝 Настройте DeepSeek в .env файле"
            echo "   См. DEEPSEEK_SETUP.md"
            ;;
        *)
            echo "❌ Неверный выбор. Копируем базовую конфигурацию..."
            cp config.env.example .env
            ;;
    esac
    
    echo ""
    echo "📖 Документация:"
    echo "   - README.md - общая информация"
    echo "   - DEEPSEEK_SETUP.md - настройка DeepSeek"
    echo ""
    echo "🔧 После настройки .env файла запустите: ./run.sh"
    exit 1
fi

# Проверяем наличие токена
if ! grep -q "TELEGRAM_BOT_TOKEN=your_bot_token_here" .env; then
    echo "✅ Конфигурация найдена"
else
    echo "❌ Пожалуйста, установите TELEGRAM_BOT_TOKEN в файле .env"
    exit 1
fi

# Устанавливаем зависимости
echo "📦 Устанавливаем зависимости..."
go mod tidy

# Запускаем приложение
echo "🚀 Запускаем бота..."
go run .

