package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient представляет HTTP клиент для внешних запросов
type HTTPClient struct {
	client  *http.Client
	baseURL string
}

// NewHTTPClient создает новый HTTP клиент
func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

// Request представляет структуру запроса
type Request struct {
	Message string `json:"message"`
	UserID  int64  `json:"user_id"`
}

// Response представляет структуру ответа
type Response struct {
	Answer string `json:"answer"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// SendRequest отправляет запрос к внешнему API
func (c *HTTPClient) SendRequest(message string, userID int64) (*Response, error) {
	req := Request{
		Message: message,
		UserID:  userID,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга JSON: %v", err)
	}

	// Если baseURL не установлен, используем встроенный агент
	if c.baseURL == "" {
		return c.processWithBuiltinAgent(message, userID)
	}

	resp, err := c.client.Post(c.baseURL+"/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ошибка HTTP запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON ответа: %v", err)
	}

	return &response, nil
}

// processWithBuiltinAgent обрабатывает запрос с помощью встроенного агента
func (c *HTTPClient) processWithBuiltinAgent(message string, userID int64) (*Response, error) {
	// Создаем агента с YandexGPT если доступен
	agent := createAgent()
	answer, err := agent.ProcessMessage(message, userID)
	if err != nil {
		return &Response{
			Answer: "Извините, произошла ошибка при обработке вашего запроса.",
			Status: "error",
			Error:  err.Error(),
		}, nil
	}

	return &Response{
		Answer: answer,
		Status: "success",
	}, nil
}

