package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/faysk/whatsapp-bot/config"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
}

// AskChatGPT envia o prompt para a API da OpenAI e retorna a resposta
func AskChatGPT(prompt string) (string, error) {
	if config.AppConfig.OpenAIKey == "" {
		return "", fmt.Errorf("❌ OPENAI_API_KEY está vazia — verifique o .env ou config.Load()")
	}

	reqBody := ChatRequest{
		Model: config.AppConfig.OpenAIModel,
		Messages: []Message{
			{Role: "system", Content: "Você é um assistente tradutor de notícias cripto para português."},
			{Role: "user", Content: prompt},
		},
		MaxTokens:   config.AppConfig.MaxTokens,
		Temperature: float32(config.AppConfig.Temperature),
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("❌ Erro ao gerar JSON da requisição: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("❌ Erro ao criar requisição HTTP: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.AppConfig.OpenAIKey)
	req.Header.Set("User-Agent", "FayskBot/1.0")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("❌ Erro ao enviar requisição à OpenAI: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Retorno inesperado ou erro HTTP
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("❌ Erro HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result ChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("❌ Erro ao decodificar resposta da OpenAI: %w", err)
	}

	if result.Error.Message != "" {
		return "", fmt.Errorf("❌ OpenAI erro: %s (%s)", result.Error.Message, result.Error.Code)
	}

	if len(result.Choices) == 0 || result.Choices[0].Message.Content == "" {
		return "🤖 A IA não respondeu nada útil.", nil
	}

	return result.Choices[0].Message.Content, nil
}
