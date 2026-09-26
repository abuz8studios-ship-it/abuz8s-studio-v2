package llm

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/abuz8s/studio/internal/config"
)

// Provider interface for all LLM providers
type Provider interface {
	Name() string
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
	Stream(ctx context.Context, req CompletionRequest, callback func(Chunk)) error
	ListModels() ([]Model, error)
	Test() (bool, error)
}

// CompletionRequest for all providers
type CompletionRequest struct {
	Model       string
	Messages    []Message
	Temperature float32
	MaxTokens   int
	TopP        float32
	System      string
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionResponse from providers
type CompletionResponse struct {
	Content      string `json:"content"`
	Model        string `json:"model"`
	TokensUsed   int    `json:"tokens_used"`
	FinishReason string `json:"finish_reason"`
}

// Chunk for streaming
type Chunk struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

// Model represents an available model
type Model struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	MaxTokens   int    `json:"max_tokens"`
	Description string `json:"description"`
	IsLocal     bool   `json:"is_local"`
}

// Client manages all LLM providers
type Client struct {
	config    *config.Config
	providers map[string]Provider
	http      *http.Client
}

// NewClient creates a new LLM client with all providers
func NewClient(cfg *config.Config) *Client {
	client := &Client{
		config:    cfg,
		providers: make(map[string]Provider),
		http:      &http.Client{Timeout: 120 * time.Second},
	}

	// Initialize all providers
	client.initProviders()
	
	return client
}

// toProviderConfig converts config.ProviderConfig to llm.ProviderConfig
func toProviderConfig(cfg config.ProviderConfig) ProviderConfig {
	return ProviderConfig{
		Enabled: cfg.Enabled,
		APIKey:  cfg.APIKey,
		BaseURL: cfg.BaseURL,
		Model:   cfg.Model,
	}
}

func (c *Client) initProviders() {
	// OpenRouter - Universal API for 100+ models
	if c.config.Providers.OpenRouter.Enabled {
		c.providers["openrouter"] = NewOpenRouterProvider(toProviderConfig(c.config.Providers.OpenRouter))
	}

	// OpenCode
	if c.config.Providers.OpenCode.Enabled {
		c.providers["opencode"] = NewOpenCodeProvider(toProviderConfig(c.config.Providers.OpenCode))
	}

	// ChatGPT / OpenAI
	if c.config.Providers.OpenAI.Enabled {
		c.providers["openai"] = NewOpenAIProvider(toProviderConfig(c.config.Providers.OpenAI))
	}

	// Codex
	if c.config.Providers.Codex.Enabled {
		c.providers["codex"] = NewCodexProvider(toProviderConfig(c.config.Providers.Codex))
	}

	// Google Gemini
	if c.config.Providers.Gemini.Enabled {
		c.providers["gemini"] = NewGeminiProvider(toProviderConfig(c.config.Providers.Gemini))
	}

	// Nous Research
	if c.config.Providers.Nous.Enabled {
		c.providers["nous"] = NewNousProvider(toProviderConfig(c.config.Providers.Nous))
	}

	// Kimi (Moonshot AI)
	if c.config.Providers.Kimi.Enabled {
		c.providers["kimi"] = NewKimiProvider(toProviderConfig(c.config.Providers.Kimi))
	}

	// MiniMax
	if c.config.Providers.MiniMax.Enabled {
		c.providers["minimax"] = NewMiniMaxProvider(toProviderConfig(c.config.Providers.MiniMax))
	}

	// GLM (Zhipu AI)
	if c.config.Providers.GLM.Enabled {
		c.providers["glm"] = NewGLMProvider(toProviderConfig(c.config.Providers.GLM))
	}

	// Local: Ollama
	if c.config.Providers.Ollama.Enabled {
		c.providers["ollama"] = NewOllamaProvider(toProviderConfig(c.config.Providers.Ollama))
	}

	// Local: LM Studio
	if c.config.Providers.LMStudio.Enabled {
		c.providers["lmstudio"] = NewLMStudioProvider(toProviderConfig(c.config.Providers.LMStudio))
	}

	// Groq
	if c.config.Providers.Groq.Enabled {
		c.providers["groq"] = NewGroqProvider(toProviderConfig(c.config.Providers.Groq))
	}

	// Together AI
	if c.config.Providers.Together.Enabled {
		c.providers["together"] = NewTogetherProvider(toProviderConfig(c.config.Providers.Together))
	}
}

// Complete sends a completion request to the best available provider
func (c *Client) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	provider := c.getBestProvider()
	if provider == nil {
		return nil, fmt.Errorf("no LLM provider available")
	}
	return provider.Complete(ctx, req)
}

// Stream streams completion from the best available provider
func (c *Client) Stream(ctx context.Context, req CompletionRequest, callback func(Chunk)) error {
	provider := c.getBestProvider()
	if provider == nil {
		return fmt.Errorf("no LLM provider available")
	}
	return provider.Stream(ctx, req, callback)
}

// getBestProvider returns the best available provider based on preferences
func (c *Client) getBestProvider() Provider {
	// Priority order based on config
	priority := c.config.Providers.Priority
	
	for _, name := range priority {
		if p, ok := c.providers[name]; ok {
			if ok, _ := p.Test(); ok {
				return p
			}
		}
	}
	
	// Fallback to any working provider
	for _, p := range c.providers {
		if ok, _ := p.Test(); ok {
			return p
		}
	}
	
	return nil
}

// GetProvider returns a specific provider
func (c *Client) GetProvider(name string) (Provider, error) {
	p, ok := c.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider not found: %s", name)
	}
	return p, nil
}

// GetProviders returns all configured providers
func (c *Client) GetProviders() []map[string]interface{} {
	var result []map[string]interface{}
	
	for name, provider := range c.providers {
		models, _ := provider.ListModels()
		result = append(result, map[string]interface{}{
			"name":     name,
			"models":   len(models),
			"is_local": name == "ollama" || name == "lmstudio",
		})
	}
	
	return result
}

// GetAvailableProviders returns names of working providers
func (c *Client) GetAvailableProviders() []string {
	var available []string
	
	for name, provider := range c.providers {
		if ok, _ := provider.Test(); ok {
			available = append(available, name)
		}
	}
	
	return available
}

// ListModels returns all available models from all providers
func (c *Client) ListModels() ([]Model, error) {
	var allModels []Model
	
	for _, provider := range c.providers {
		models, err := provider.ListModels()
		if err != nil {
			continue
		}
		allModels = append(allModels, models...)
	}
	
	return allModels, nil
}

// ListLocalModels returns only local models
func (c *Client) ListLocalModels() ([]Model, error) {
	var localModels []Model
	
	if p, ok := c.providers["ollama"]; ok {
		models, _ := p.ListModels()
		localModels = append(localModels, models...)
	}
	
	if p, ok := c.providers["lmstudio"]; ok {
		models, _ := p.ListModels()
		localModels = append(localModels, models...)
	}
	
	return localModels, nil
}

// TestProvider tests a specific provider
func (c *Client) TestProvider(name string) (bool, error) {
	p, ok := c.providers[name]
	if !ok {
		return false, fmt.Errorf("provider not found: %s", name)
	}
	return p.Test()
}

// UpdateProvider updates provider configuration
func (c *Client) UpdateProvider(name string, cfg map[string]interface{}) error {
	// Implementation depends on provider type
	return nil
}

// GenerateContent is a high-level function for content generation
func (c *Client) GenerateContent(ctx context.Context, prompt string, contentType string, voiceProfile string) (string, error) {
	req := CompletionRequest{
		Model:       c.config.Providers.DefaultModel,
		Temperature: 0.8,
		MaxTokens:   4000,
		Messages: []Message{
			{Role: "system", Content: c.buildSystemPrompt(contentType, voiceProfile)},
			{Role: "user", Content: prompt},
		},
	}
	
	resp, err := c.Complete(ctx, req)
	if err != nil {
		return "", err
	}
	
	return resp.Content, nil
}

func (c *Client) buildSystemPrompt(contentType, voiceProfile string) string {
	base := "You are ABUZ8s Studio, an expert content creation AI. "
	
	switch contentType {
	case "youtube_script":
		base += "Generate engaging YouTube video scripts with hooks, structure, and calls-to-action. "
	case "x_post":
		base += "Write viral X/Twitter posts with strong hooks and engagement. "
	case "blog":
		base += "Create SEO-optimized blog posts with proper structure. "
	case "newsletter":
		base += "Write compelling newsletter content that drives opens and clicks. "
	}
	
	if voiceProfile != "" {
		base += fmt.Sprintf("Match the voice profile: %s. ", voiceProfile)
	}
	
	return base
}
