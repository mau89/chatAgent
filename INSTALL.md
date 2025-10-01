# Установка и настройка

## Установка Go

### macOS (через Homebrew)
```bash
brew install go
```

### macOS (ручная установка)
1. Скачайте Go с [golang.org](https://golang.org/dl/)
2. Установите пакет
3. Добавьте в PATH:
```bash
export PATH=$PATH:/usr/local/go/bin
```

### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install golang-go
```

### Windows
1. Скачайте установщик с [golang.org](https://golang.org/dl/)
2. Запустите установщик
3. Перезапустите терминал

## Проверка установки

```bash
go version
```

Должно показать версию Go 1.21 или выше.

## Настройка проекта

1. **Клонируйте или скачайте проект**

2. **Перейдите в директорию проекта**
```bash
cd chatAgent
```

3. **Установите зависимости**
```bash
go mod tidy
```

4. **Настройте переменные окружения**
```bash
cp config.env.example .env
```

5. **Отредактируйте .env файл**
```env
TELEGRAM_BOT_TOKEN=your_actual_bot_token_here
PORT=8080
EXTERNAL_API_URL=
```

6. **Получите токен бота**
   - Найдите [@BotFather](https://t.me/BotFather) в Telegram
   - Отправьте `/newbot`
   - Следуйте инструкциям
   - Скопируйте токен в .env файл

## Запуск

### Быстрый запуск
```bash
./run.sh
```

### Ручной запуск
```bash
go run .
```

## Тестирование

### Тест API (в отдельном терминале)
```bash
./test_api.sh
```

### Тест в браузере
Откройте http://localhost:8080

## Структура после установки

```
chatAgent/
├── main.go
├── agent.go
├── http_client.go
├── telegram_bot.go
├── http_server.go
├── go.mod
├── go.sum
├── .env
├── config.env.example
├── run.sh
├── test_api.sh
├── README.md
├── INSTALL.md
└── .gitignore
```

## Возможные проблемы

### Ошибка "command not found: go"
- Убедитесь, что Go установлен
- Проверьте PATH: `echo $PATH`
- Перезапустите терминал

### Ошибка "module not found"
- Выполните `go mod tidy`
- Проверьте интернет соединение

### Ошибка "TELEGRAM_BOT_TOKEN не установлен"
- Проверьте файл .env
- Убедитесь, что токен скопирован правильно
- Убедитесь, что нет пробелов в токене

### Бот не отвечает
- Проверьте токен бота
- Убедитесь, что бот запущен
- Проверьте логи в консоли

