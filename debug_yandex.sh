#!/bin/bash

# Скрипт для диагностики проблем с Yandex Cloud

echo "🔍 Диагностика проблем с Yandex Cloud"
echo "======================================"

# Проверяем наличие .env файла
if [ ! -f .env ]; then
    echo "❌ Файл .env не найден"
    echo "📝 Создайте файл .env на основе config.env.example"
    exit 1
fi

# Загружаем переменные окружения
source .env

echo ""
echo "📋 Проверка конфигурации:"
echo "------------------------"

# Проверяем токен бота
if [ "$TELEGRAM_BOT_TOKEN" = "your_telegram_bot_token_here" ] || [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo "❌ TELEGRAM_BOT_TOKEN не настроен"
else
    echo "✅ TELEGRAM_BOT_TOKEN настроен"
fi

# Проверяем API ключ
if [ "$YANDEX_GPT_API_KEY" = "your_yandex_api_key_here" ] || [ -z "$YANDEX_GPT_API_KEY" ]; then
    echo "❌ YANDEX_GPT_API_KEY не настроен"
    echo "   Текущее значение: $YANDEX_GPT_API_KEY"
else
    echo "✅ YANDEX_GPT_API_KEY настроен"
    echo "   Формат ключа: ${YANDEX_GPT_API_KEY:0:10}..."
fi

# Проверяем Folder ID
if [ "$YANDEX_GPT_FOLDER_ID" = "your_folder_id_here" ] || [ -z "$YANDEX_GPT_FOLDER_ID" ]; then
    echo "❌ YANDEX_GPT_FOLDER_ID не настроен"
    echo "   Текущее значение: $YANDEX_GPT_FOLDER_ID"
else
    echo "✅ YANDEX_GPT_FOLDER_ID настроен"
    echo "   Folder ID: $YANDEX_GPT_FOLDER_ID"
fi

echo ""
echo "🔧 Возможные решения:"
echo "--------------------"

if [ "$YANDEX_GPT_API_KEY" = "your_yandex_api_key_here" ] || [ -z "$YANDEX_GPT_API_KEY" ]; then
    echo "1. 📖 Следуйте инструкции в YANDEX_SETUP.md"
    echo "2. 🔑 Попробуйте создать IAM токен вместо API ключа:"
    echo "   - Перейдите в https://console.cloud.yandex.ru/iam/tokens"
    echo "   - Нажмите 'Получить IAM токен'"
    echo "   - Скопируйте токен (начинается с t1.)"
    echo "   - Замените YANDEX_GPT_API_KEY в .env файле"
    echo ""
    echo "3. 🛠️ Или используйте YC CLI:"
    echo "   curl -sSL https://storage.yandexcloud.net/yandexcloud-yc/install.sh | bash"
    echo "   yc init"
    echo "   yc iam api-key create --service-account-name ai-bot-service-account"
fi

if [ "$YANDEX_GPT_FOLDER_ID" = "your_folder_id_here" ] || [ -z "$YANDEX_GPT_FOLDER_ID" ]; then
    echo "4. 📁 Получите Folder ID:"
    echo "   - Перейдите в https://console.cloud.yandex.ru/"
    echo "   - В левом меню выберите 'Облако'"
    echo "   - Скопируйте ID папки (начинается с b1g...)"
fi

echo ""
echo "🚨 Частые проблемы:"
echo "------------------"
echo "• Аккаунт не активирован - нужна банковская карта"
echo "• Неправильная роль у сервисного аккаунта - должна быть ai.languageModels.user"
echo "• Yandex GPT не активирован - перейдите в раздел Foundation Models"
echo "• Региональные ограничения - попробуйте другой регион"

echo ""
echo "📞 Поддержка:"
echo "-------------"
echo "• Документация: https://cloud.yandex.ru/docs"
echo "• Поддержка: https://cloud.yandex.ru/support"
echo "• Telegram: @yandexcloud_support"

echo ""
echo "✅ После исправления проблем запустите:"
echo "   ./test_yandex.sh"
