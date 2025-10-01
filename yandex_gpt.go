package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// YandexGPTClient представляет клиент для работы с Yandex GPT API
type YandexGPTClient struct {
	apiKey    string
	folderID  string
	httpClient *http.Client
	baseURL   string
}

// YandexGPTRequest представляет запрос к Yandex GPT API
type YandexGPTRequest struct {
	ModelURI string `json:"modelUri"`
	CompletionOptions struct {
		Stream    bool    `json:"stream"`
		Temperature float64 `json:"temperature"`
		MaxTokens  int    `json:"maxTokens"`
	} `json:"completionOptions"`
	Messages []YandexGPTMessage `json:"messages"`
}

// YandexGPTMessage представляет сообщение в Yandex GPT
type YandexGPTMessage struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

// YandexGPTResponse представляет ответ от Yandex GPT API
type YandexGPTResponse struct {
	Result struct {
		Alternatives []struct {
			Message struct {
				Role string `json:"role"`
				Text string `json:"text"`
			} `json:"message"`
			Status string `json:"status"`
		} `json:"alternatives"`
		Usage struct {
			InputTextTokens  string `json:"inputTextTokens"`
			CompletionTokens string `json:"completionTokens"`
			TotalTokens      string `json:"totalTokens"`
		} `json:"usage"`
		ModelVersion string `json:"modelVersion"`
	} `json:"result"`
}

// NewYandexGPTClient создает новый клиент Yandex GPT
func NewYandexGPTClient(apiKey, folderID string) *YandexGPTClient {
	return &YandexGPTClient{
		apiKey:   apiKey,
		folderID: folderID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://llm.api.cloud.yandex.net/foundationModels/v1/completion",
	}
}

// GenerateResponse генерирует ответ с помощью Yandex GPT
func (c *YandexGPTClient) GenerateResponse(message string, userID int64) (string, error) {
	// Формируем запрос
	request := YandexGPTRequest{
		ModelURI: fmt.Sprintf("gpt://%s/yandexgpt-lite", c.folderID),
		CompletionOptions: struct {
			Stream    bool    `json:"stream"`
			Temperature float64 `json:"temperature"`
			MaxTokens  int    `json:"maxTokens"`
		}{
			Stream:     false,
			Temperature: 0.6,
			MaxTokens:  2000,
		},
		Messages: []YandexGPTMessage{
			{
				Role: "user",
				Text: message,
			},
		},
	}

	// Конвертируем в JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("ошибка маршалинга JSON: %v", err)
	}

	// Создаем HTTP запрос
	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %v", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Api-Key "+c.apiKey)

	// Отправляем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка HTTP запроса: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API вернул ошибку %d: %s", resp.StatusCode, string(body))
	}

	// Парсим ответ
	var response YandexGPTResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("ошибка парсинга JSON ответа: %v", err)
	}

	// Проверяем, что есть ответ
	if len(response.Result.Alternatives) == 0 {
		return "", fmt.Errorf("пустой ответ от API")
	}

	// Возвращаем текст ответа
	generatedText := response.Result.Alternatives[0].Message.Text
	if generatedText == "" {
		return "", fmt.Errorf("пустой текст ответа")
	}

	return generatedText, nil
}

// IsAvailable проверяет доступность Yandex GPT API
func (c *YandexGPTClient) IsAvailable() bool {
	// Простой тестовый запрос
	_, err := c.GenerateResponse("Привет", 0)
	return err == nil
}
