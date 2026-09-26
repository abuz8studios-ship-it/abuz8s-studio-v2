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
)

// ==================== NOUS RESEARCH ====================

type NousProvider struct {
	config ProviderConfig
	client *http.Client
}

func NewNousProvider(cfg ProviderConfig) *NousProvider {
	return &NousProvider{
		config: cfg,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *NousProvider) Name() string { return "nous" }

func (p *NousProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	url := "https://api.nousresearch.com/v1/chat/completions"
	
	payload := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
	}
	
	jsonData, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	
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
		return nil, fmt.Errorf("no response from Nous")
	}
	
	return &CompletionResponse{
		Content:      result.Choices[0].Message.Content,
		Model:        result.Model,
		TokensUsed:   result.Usage.TotalTokens,
		FinishReason: result.Choices[0].FinishReason,
	}, nil
}

func (p *NousProvider) Stream(ctx context.Context, req CompletionRequest, callback func(Chunk)) error {
	return fmt.Errorf("streaming not implemented for Nous")
}

func (p *NousProvider) ListModels() ([]Model, error) {
	return []Model{
		{ID: "nous-hermes-2", Name: "Nous Hermes 2", Provider: "nous", MaxTokens: 4096, IsLocal: false},
	}, nil
}

func (p *NousProvider) Test() (bool, error) {
	_, err := p.ListModels()
	return err == nil, err
}

// ==================== KIMI (Moonshot AI) ====================

type KimiProvider struct {
	config ProviderConfig
	client *http.Client
}

func NewKimiProvider(cfg ProviderConfig) *KimiProvider {
	return &KimiProvider{
		config: cfg,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *KimiProvider) Name() string { return "kimi" }

func (p *KimiProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	url := "https://api.moonshot.cn/v1/chat/completions"
	
	payload := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
	}
	
	jsonData, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	
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
		return nil, fmt.Errorf("no response from Kimi")
	}
	
	return &CompletionResponse{
		Content:      result.Choices[0].Message.Content,
		Model:        result.Model,
		TokensUsed:   result.Usage.TotalTokens,
		FinishReason: result.Choices[0].FinishReason,
	}, nil
}

func (p *KimiProvider) Stream(ctx context.Context, req CompletionRequest, callback func(Chunk)) error {
	return fmt.Errorf("streaming not implemented for Kimi")
}

func (p *KimiProvider) ListModels() ([]Model, error) {
	return []Model{
		{ID: "moonshot-v1-8k", Name: "Moonshot v1 8K", Provider: "kimi", MaxTokens: 8192, IsLocal: false},
		{ID: "moonshot-v1-32k", Name: "Moonshot v1 32K", Provider: "kimi", MaxTokens: 32768, IsLocal: false},
	}, nil
}

func (p *KimiProvider) Test() (bool, error) {
	_, err := p.ListModels()
	return err == nil, err
}

// ==================== MINIMAX ====================

type MiniMaxProvider struct {
	config ProviderConfig
	client *http.Client
}

func NewMiniMaxProvider(cfg ProviderConfig) *MiniMaxProvider {
	return &MiniMaxProvider{
		config: cfg,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *MiniMaxProvider) Name() string { return "minimax" }

func (p *MiniMaxProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	url := "https://api.minimax.chat/v1/text/chatcompletion_v2"
	
	payload := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
	}
	
	jsonData, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
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
	
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	
	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no response from MiniMax")
	}
	
	return &CompletionResponse{
		Content:      result.Choices[0].Message.Content,
		Model:        result.Model,
		TokensUsed:   result.Usage.TotalTokens,
		FinishReason: result.Choices[0].FinishReason,
	}, nil
}

func (p *MiniMaxProvider) Stream(ctx context.Context, req CompletionRequest, callback func(Chunk)) error {
	return fmt.Errorf("streaming not implemented for MiniMax")
}

func (p *MiniMaxProvider) ListModels() ([]Model, error) {
	return []Model{
		{ID: "abab6.5s-chat", Name: "MiniMax abab6.5s", Provider: "minimax", MaxTokens: 8192, IsLocal: false},
	}, nil
}

func (p *MiniMaxProvider) Test() (bool, error) {
	_, err := p.ListModels()
	return err == nil, err
}

// ==================== GLM (Zhipu AI) ====================

type GLMProvider struct {
	config ProviderConfig
	client *http.Client
}

func NewGLMProvider(cfg ProviderConfig) *GLMProvider {
	return &GLMProvider{
		config: cfg,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *GLMProvider) Name() string { return "glm" }

func (p *GLMProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	url := "https://open.bigmodel.cn/api/paas/v4/chat/completions"
	
	payload := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
	}
	
	jsonData, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Authorization", p.config.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	
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
		return nil, fmt.Errorf("no response from GLM")
	}
	
	return &CompletionResponse{
		Content:      result.Choices[0].Message.Content,
		Model:        result.Model,
		TokensUsed:   result.Usage.TotalTokens,
		FinishReason: result.Choices[0].FinishReason,
	}, nil
}

func (p *GLMProvider) Stream(ctx context.Context, req CompletionRequest, callback func(Chunk)) error {
	return fmt.Errorf("streaming not implemented for GLM")
}

func (p *GLMProvider) ListModels() ([]Model, error) {
	return []Model{
		{ID: "glm-4", Name: "GLM-4", Provider: "glm", MaxTokens: 8192, IsLocal: false},
		{ID: "glm-4-flash", Name: "GLM-4 Flash", Provider: "glm", MaxTokens: 8192, IsLocal: false},
	}, nil
}

func (p *GLMProvider) Test() (bool, error) {
	_, err := p.ListModels()
	return err == nil, err
}

// ==================== LM STUDIO ====================

type LMStudioProvider struct {
	config ProviderConfig
	client *http.Client
}

func NewLMStudioProvider(cfg ProviderConfig) *LMStudioProvider {
	return &LMStudioProvider{
		config: cfg,
		client: &http.Client{Timeout: 300 * time.Second},
	}
}

func (p *LMStudioProvider) Name() string { return "lmstudio" }

func (p *LMStudioProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	url := p.config.BaseURL + "/chat/completions"
	
	payload := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
		"stream":    false,
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
		return nil, fmt.Errorf("no response from LM Studio")
	}
	
	return &CompletionResponse{
		Content:      result.Choices[0].Message.Content,
		Model:        result.Model,
		TokensUsed:   result.Usage.TotalTokens,
		FinishReason: result.Choices[0].FinishReason,
	}, nil
}

func (p *LMStudioProvider) Stream(ctx context.Context, req CompletionRequest, callback func(Chunk)) error {
	url := p.config.BaseURL + "/chat/completions"
	
	payload := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
		"stream":   true,
	}
	
	jsonData, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return handleSSEStream(resp.Body, callback)
}

func (p *LMStudioProvider) ListModels() ([]Model, error) {
	url := p.config.BaseURL + "/models"
	
	resp, err := p.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	var models []Model
	for _, m := range result.Data {
		models = append(models, Model{
			ID:       m.ID,
			Name:     m.ID,
			Provider: "lmstudio",
			IsLocal:  true,
		})
	}
	
	return models, nil
}

func (p *LMStudioProvider) Test() (bool, error) {
	_, err := p.ListModels()
	return err == nil, err
}

// ==================== GROQ ====================

type GroqProvider struct {
	config ProviderConfig
	client *http.Client
}

func NewGroqProvider(cfg ProviderConfig) *GroqProvider {
	return &GroqProvider{
		config: cfg,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *GroqProvider) Name() string { return "groq" }

func (p *GroqProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"
	
	payload := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
	}
	
	jsonData, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	
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
		return nil, fmt.Errorf("no response from Groq")
	}
	
	return &CompletionResponse{
		Content:      result.Choices[0].Message.Content,
		Model:        result.Model,
		TokensUsed:   result.Usage.TotalTokens,
		FinishReason: result.Choices[0].FinishReason,
	}, nil
}

func (p *GroqProvider) Stream(ctx context.Context, req CompletionRequest, callback func(Chunk)) error {
	return fmt.Errorf("streaming not implemented for Groq")
}

func (p *GroqProvider) ListModels() ([]Model, error) {
	return []Model{
		{ID: "llama3-8b-8192", Name: "Llama 3 8B", Provider: "groq", MaxTokens: 8192, IsLocal: false},
		{ID: "llama3-70b-8192", Name: "Llama 3 70B", Provider: "groq", MaxTokens: 8192, IsLocal: false},
		{ID: "mixtral-8x7b-32768", Name: "Mixtral 8x7B", Provider: "groq", MaxTokens: 32768, IsLocal: false},
	}, nil
}

func (p *GroqProvider) Test() (bool, error) {
	_, err := p.ListModels()
	return err == nil, err
}

// ==================== TOGETHER AI ====================

type TogetherProvider struct {
	config ProviderConfig
	client *http.Client
}

func NewTogetherProvider(cfg ProviderConfig) *TogetherProvider {
	return &TogetherProvider{
		config: cfg,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *TogetherProvider) Name() string { return "together" }

func (p *TogetherProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	url := "https://api.together.xyz/v1/chat/completions"
	
	payload := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
	}
	
	jsonData, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	
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
		return nil, fmt.Errorf("no response from Together AI")
	}
	
	return &CompletionResponse{
		Content:      result.Choices[0].Message.Content,
		Model:        result.Model,
		TokensUsed:   result.Usage.TotalTokens,
		FinishReason: result.Choices[0].FinishReason,
	}, nil
}

func (p *TogetherProvider) Stream(ctx context.Context, req CompletionRequest, callback func(Chunk)) error {
	return fmt.Errorf("streaming not implemented for Together AI")
}

func (p *TogetherProvider) ListModels() ([]Model, error) {
	return []Model{
		{ID: "meta-llama/Llama-3-70b-chat-hf", Name: "Llama 3 70B", Provider: "together", MaxTokens: 8192, IsLocal: false},
	}, nil
}

func (p *TogetherProvider) Test() (bool, error) {
	_, err := p.ListModels()
	return err == nil, err
}

// ==================== GEMINI (Google) ====================

type GeminiProvider struct {
	config ProviderConfig
	client *http.Client
}

func NewGeminiProvider(cfg ProviderConfig) *GeminiProvider {
	return &GeminiProvider{
		config: cfg,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *GeminiProvider) Name() string { return "gemini" }

func (p *GeminiProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	model := req.Model
	if model == "" {
		model = "gemini-1.5-flash"
	}
	
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, p.config.APIKey)
	
	// Convert messages to Gemini format
	var contents []map[string]interface{}
	for _, msg := range req.Messages {
		role := "user"
		if msg.Role == "assistant" {
			role = "model"
		}
		contents = append(contents, map[string]interface{}{
			"role": role,
			"parts": []map[string]string{
				{"text": msg.Content},
			},
		})
	}
	
	payload := map[string]interface{}{
		"contents": contents,
		"generationConfig": map[string]interface{}{
			"temperature":     req.Temperature,
			"maxOutputTokens": req.MaxTokens,
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
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
			FinishReason string `json:"finishReason"`
		} `json:"candidates"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	if len(result.Candidates) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}
	
	content := ""
	for _, part := range result.Candidates[0].Content.Parts {
		content += part.Text
	}
	
	return &CompletionResponse{
		Content:      content,
		Model:        model,
		FinishReason: result.Candidates[0].FinishReason,
	}, nil
}

func (p *GeminiProvider) Stream(ctx context.Context, req CompletionRequest, callback func(Chunk)) error {
	return fmt.Errorf("streaming not implemented for Gemini")
}

func (p *GeminiProvider) ListModels() ([]Model, error) {
	return []Model{
		{ID: "gemini-1.5-flash", Name: "Gemini 1.5 Flash", Provider: "gemini", MaxTokens: 8192, IsLocal: false},
		{ID: "gemini-1.5-pro", Name: "Gemini 1.5 Pro", Provider: "gemini", MaxTokens: 8192, IsLocal: false},
	}, nil
}

func (p *GeminiProvider) Test() (bool, error) {
	_, err := p.ListModels()
	return err == nil, err
}

// ==================== CODEX (OpenAI) ====================

type CodexProvider struct {
	config ProviderConfig
	client *http.Client
}

func NewCodexProvider(cfg ProviderConfig) *CodexProvider {
	return &CodexProvider{
		config: cfg,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *CodexProvider) Name() string { return "codex" }

func (p *CodexProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	// Codex uses OpenAI-compatible API
	url := "https://api.openai.com/v1/chat/completions"
	
	payload := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
	}
	
	jsonData, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	
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
		return nil, fmt.Errorf("no response from Codex")
	}
	
	return &CompletionResponse{
		Content:      result.Choices[0].Message.Content,
		Model:        result.Model,
		TokensUsed:   result.Usage.TotalTokens,
		FinishReason: result.Choices[0].FinishReason,
	}, nil
}

func (p *CodexProvider) Stream(ctx context.Context, req CompletionRequest, callback func(Chunk)) error {
	return fmt.Errorf("streaming not implemented for Codex")
}

func (p *CodexProvider) ListModels() ([]Model, error) {
	return []Model{
		{ID: "codex-latest", Name: "Codex", Provider: "codex", MaxTokens: 8192, IsLocal: false},
	}, nil
}

func (p *CodexProvider) Test() (bool, error) {
	_, err := p.ListModels()
	return err == nil, err
}

// ==================== OPENCODE ====================

type OpenCodeProvider struct {
	config ProviderConfig
	client *http.Client
}

func NewOpenCodeProvider(cfg ProviderConfig) *OpenCodeProvider {
	return &OpenCodeProvider{
		config: cfg,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *OpenCodeProvider) Name() string { return "opencode" }

func (p *OpenCodeProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	url := p.config.BaseURL + "/chat/completions"
	if url == "" {
		url = "https://api.opencode.ai/v1/chat/completions"
	}
	
	payload := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
	}
	
	jsonData, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	
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
		return nil, fmt.Errorf("no response from OpenCode")
	}
	
	return &CompletionResponse{
		Content:      result.Choices[0].Message.Content,
		Model:        result.Model,
		TokensUsed:   result.Usage.TotalTokens,
		FinishReason: result.Choices[0].FinishReason,
	}, nil
}

func (p *OpenCodeProvider) Stream(ctx context.Context, req CompletionRequest, callback func(Chunk)) error {
	return fmt.Errorf("streaming not implemented for OpenCode")
}

func (p *OpenCodeProvider) ListModels() ([]Model, error) {
	return []Model{
		{ID: "opencode-1", Name: "OpenCode 1", Provider: "opencode", MaxTokens: 8192, IsLocal: false},
	}, nil
}

func (p *OpenCodeProvider) Test() (bool, error) {
	_, err := p.ListModels()
	return err == nil, err
}

// handleSSEStream handles Server-Sent Events streaming
func handleSSEStream(body io.ReadCloser, callback func(Chunk)) error {
	defer body.Close()
	
	reader := bufio.NewReader(body)
	
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
