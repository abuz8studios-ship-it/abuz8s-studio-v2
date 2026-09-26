package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// OpenRouterProvider implements the OpenRouter API
type OpenRouterProvider struct {
	config ProviderConfig
	client *http.Client
}

// ProviderConfig holds common provider configuration
type ProviderConfig struct {
	Enabled bool   `json:"enabled"`
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
	Model   string `json:"model"`
}

// NewOpenRouterProvider creates a new OpenRouter provider
func NewOpenRouterProvider(cfg ProviderConfig) *OpenRouterProvider {
	return &OpenRouterProvider{
		config: cfg,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *OpenRouterProvider) Name() string {
	return "openrouter"
}

func (p *OpenRouterProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	payload := map[string]interface{}{
		"model":       req.Model,
		"messages":    req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
		"top_p":       req.TopP,
	}

	if req.System != "" {
		payload["messages"] = append([]Message{{Role: "system", Content: req.System}}, req.Messages...)
	}

	jsonData, _ := json.Marshal(payload)

	httpReq, _ := http.NewRequestWithContext(ctx, "POST",
		"https://openrouter.ai/api/v1/chat/completions",
		bytes.NewBuffer(jsonData))

	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("HTTP-Referer", "https://abuz8s.studio")
	httpReq.Header.Set("X-Title", "ABUZ8s Studio")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Model string `json:"model"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenRouter")
	}

	return &CompletionResponse{
		Content:      result.Choices[0].Message.Content,
		Model:        result.Model,
		TokensUsed:   result.Usage.TotalTokens,
		FinishReason: result.Choices[0].FinishReason,
	}, nil
}

func (p *OpenRouterProvider) Stream(ctx context.Context, req CompletionRequest, callback func(Chunk)) error {
	payload := map[string]interface{}{
		"model":       req.Model,
		"messages":    req.Messages,
		"stream":      true,
		"temperature": req.Temperature,
	}

	jsonData, _ := json.Marshal(payload)

	httpReq, _ := http.NewRequestWithContext(ctx, "POST",
		"https://openrouter.ai/api/v1/chat/completions",
		bytes.NewBuffer(jsonData))

	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	httpReq.Header.Set("HTTP-Referer", "https://abuz8s.studio")
	httpReq.Header.Set("X-Title", "ABUZ8s Studio")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			callback(Chunk{Done: true})
			break
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}

		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) > 0 {
			callback(Chunk{Content: chunk.Choices[0].Delta.Content})
		}
	}

	return nil
}

func (p *OpenRouterProvider) ListModels() ([]Model, error) {
	req, _ := http.NewRequest("GET", "https://openrouter.ai/api/v1/models", nil)
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Context     int    `json:"context_length"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var models []Model
	for _, m := range result.Data {
		models = append(models, Model{
			ID:          m.ID,
			Name:        m.Name,
			Provider:    "openrouter",
			MaxTokens:   m.Context,
			Description: m.Description,
			IsLocal:     false,
		})
	}

	return models, nil
}

func (p *OpenRouterProvider) Test() (bool, error) {
	_, err := p.ListModels()
	return err == nil, err
}

// OpenAIProvider implements OpenAI API
type OpenAIProvider struct {
	config ProviderConfig
	client *openai.Client
}

func NewOpenAIProvider(cfg ProviderConfig) *OpenAIProvider {
	return &OpenAIProvider{
		config: cfg,
		client: openai.NewClient(cfg.APIKey),
	}
}

func (p *OpenAIProvider) Name() string {
	return "openai"
}

func (p *OpenAIProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	messages := make([]openai.ChatCompletionMessage, len(req.Messages))
	for i, m := range req.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		TopP:        req.TopP,
	})

	if err != nil {
		return nil, err
	}

	return &CompletionResponse{
		Content:      resp.Choices[0].Message.Content,
		Model:        resp.Model,
		TokensUsed:   resp.Usage.TotalTokens,
		FinishReason: string(resp.Choices[0].FinishReason),
	}, nil
}

func (p *OpenAIProvider) Stream(ctx context.Context, req CompletionRequest, callback func(Chunk)) error {
	messages := make([]openai.ChatCompletionMessage, len(req.Messages))
	for i, m := range req.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	stream, err := p.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stream:      true,
	})

	if err != nil {
		return err
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			callback(Chunk{Done: true})
			break
		}
		if err != nil {
			return err
		}

		if len(response.Choices) > 0 {
			callback(Chunk{Content: response.Choices[0].Delta.Content})
		}
	}

	return nil
}

func (p *OpenAIProvider) ListModels() ([]Model, error) {
	return []Model{
		{ID: "gpt-4-turbo", Name: "GPT-4 Turbo", Provider: "openai", MaxTokens: 128000, IsLocal: false},
		{ID: "gpt-4", Name: "GPT-4", Provider: "openai", MaxTokens: 8192, IsLocal: false},
		{ID: "gpt-3.5-turbo", Name: "GPT-3.5 Turbo", Provider: "openai", MaxTokens: 16385, IsLocal: false},
	}, nil
}

func (p *OpenAIProvider) Test() (bool, error) {
	_, err := p.client.ListModels(context.Background())
	return err == nil, err
}

// OllamaProvider implements local Ollama API
type OllamaProvider struct {
	config ProviderConfig
	client *http.Client
}

func NewOllamaProvider(cfg ProviderConfig) *OllamaProvider {
	return &OllamaProvider{
		config: cfg,
		client: &http.Client{Timeout: 300 * time.Second},
	}
}

func (p *OllamaProvider) Name() string {
	return "ollama"
}

func (p *OllamaProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	url := fmt.Sprintf("%s/api/generate", p.config.BaseURL)

	prompt := ""
	for _, m := range req.Messages {
		prompt += fmt.Sprintf("%s: %s\n", m.Role, m.Content)
	}

	payload := map[string]interface{}{
		"model":  req.Model,
		"prompt": prompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": req.Temperature,
			"num_predict": req.MaxTokens,
			"top_p":       req.TopP,
		},
	}

	jsonData, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Response string `json:"response"`
		Model    string `json:"model"`
		Done     bool   `json:"done"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &CompletionResponse{
		Content: result.Response,
		Model:   result.Model,
	}, nil
}

func (p *OllamaProvider) Stream(ctx context.Context, req CompletionRequest, callback func(Chunk)) error {
	url := fmt.Sprintf("%s/api/generate", p.config.BaseURL)

	prompt := ""
	for _, m := range req.Messages {
		prompt += fmt.Sprintf("%s: %s\n", m.Role, m.Content)
	}

	payload := map[string]interface{}{
		"model":  req.Model,
		"prompt": prompt,
		"stream": true,
	}

	jsonData, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		var chunk struct {
			Response string `json:"response"`
			Done     bool   `json:"done"`
		}

		if err := json.Unmarshal(line, &chunk); err != nil {
			continue
		}

		callback(Chunk{Content: chunk.Response, Done: chunk.Done})

		if chunk.Done {
			break
		}
	}

	return nil
}

func (p *OllamaProvider) ListModels() ([]Model, error) {
	url := fmt.Sprintf("%s/api/tags", p.config.BaseURL)

	resp, err := p.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Models []struct {
			Name       string `json:"name"`
			Model      string `json:"model"`
			Size       int64  `json:"size"`
			ModifiedAt string `json:"modified_at"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var models []Model
	for _, m := range result.Models {
		models = append(models, Model{
			ID:       m.Name,
			Name:     m.Name,
			Provider: "ollama",
			IsLocal:  true,
		})
	}

	return models, nil
}

func (p *OllamaProvider) Test() (bool, error) {
	_, err := p.ListModels()
	return err == nil, err
}
