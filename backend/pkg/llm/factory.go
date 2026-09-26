package llm

import (
	"fmt"
)

// ProviderNames lists every supported provider id.
var ProviderNames = []string{
	"openrouter", "opencode", "openai", "codex", "gemini",
	"nous", "kimi", "minimax", "glm", "groq", "together",
	"ollama", "lmstudio",
}

// NewProvider builds a provider by name. Unknown names error.
func NewProvider(name string, cfg ProviderConfig) (Provider, error) {
	switch name {
	case "openrouter":
		return NewOpenRouterProvider(cfg), nil
	case "opencode":
		return NewOpenCodeProvider(cfg), nil
	case "openai":
		return NewOpenAIProvider(cfg), nil
	case "codex":
		return NewCodexProvider(cfg), nil
	case "gemini":
		return NewGeminiProvider(cfg), nil
	case "nous":
		return NewNousProvider(cfg), nil
	case "kimi":
		return NewKimiProvider(cfg), nil
	case "minimax":
		return NewMiniMaxProvider(cfg), nil
	case "glm":
		return NewGLMProvider(cfg), nil
	case "groq":
		return NewGroqProvider(cfg), nil
	case "together":
		return NewTogetherProvider(cfg), nil
	case "ollama":
		return NewOllamaProvider(cfg), nil
	case "lmstudio":
		return NewLMStudioProvider(cfg), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}
