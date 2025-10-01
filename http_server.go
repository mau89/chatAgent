package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// HTTPServer представляет HTTP сервер для API
type HTTPServer struct {
	port       string
	httpClient *HTTPClient
}

// NewHTTPServer создает новый HTTP сервер
func NewHTTPServer(port string) *HTTPServer {
	return &HTTPServer{
		port:       port,
		httpClient: NewHTTPClient(""), // Используем агента с Yandex GPT
	}
}

// Start запускает HTTP сервер
func (s *HTTPServer) Start() error {
	http.HandleFunc("/chat", s.handleChat)
	http.HandleFunc("/health", s.handleHealth)
	http.HandleFunc("/", s.handleRoot)

	log.Printf("HTTP сервер запущен на порту %s", s.port)
	return http.ListenAndServe(":"+s.port, nil)
}

// handleChat обрабатывает запросы к /chat
func (s *HTTPServer) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	// Обрабатываем запрос через HTTP клиент
	response, err := s.httpClient.SendRequest(req.Message, req.UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка обработки: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleHealth обрабатывает запросы к /health
func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"service": "chat-agent-api",
	})
}

// handleRoot обрабатывает запросы к корневому пути
func (s *HTTPServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Chat Agent API</title>
    <meta charset="utf-8">
</head>
<body>
    <h1>🤖 Chat Agent API</h1>
    <p>API для интеллектуального агента</p>
    
    <h2>Доступные эндпоинты:</h2>
    <ul>
        <li><strong>POST /chat</strong> - Отправить сообщение агенту</li>
        <li><strong>GET /health</strong> - Проверка состояния сервиса</li>
    </ul>
    
    <h2>Пример запроса к /chat:</h2>
    <pre>
{
    "message": "Привет!",
    "user_id": 12345
}
    </pre>
    
    <h2>Пример ответа:</h2>
    <pre>
{
    "answer": "Привет! Как дела?",
    "status": "success"
}
    </pre>
</body>
</html>
    `)
}

