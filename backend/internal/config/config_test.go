package config

import (
	"testing"
)

func TestDefaultConfigProvidersDisabled(t *testing.T) {
	cfg := DefaultConfig()
	p := cfg.Providers
	flags := []bool{
		p.OpenRouter.Enabled, p.OpenCode.Enabled, p.OpenAI.Enabled, p.Codex.Enabled,
		p.Gemini.Enabled, p.Nous.Enabled, p.Kimi.Enabled, p.MiniMax.Enabled,
		p.GLM.Enabled, p.Groq.Enabled, p.Together.Enabled,
		p.Ollama.Enabled, p.LMStudio.Enabled,
	}
	for i, f := range flags {
		if f {
			t.Fatalf("provider %d enabled by default, want all disabled", i)
		}
	}
	if p.Ollama.BaseURL != "http://localhost:11434" {
		t.Fatalf("ollama baseurl = %q", p.Ollama.BaseURL)
	}
	if cfg.DatabasePath == "" {
		t.Fatal("DatabasePath empty")
	}
}

// Partial nested updates must not wipe sibling provider configs.
func TestUpdateDeepMergesProviders(t *testing.T) {
	cfg := DefaultConfig()
	err := cfg.Update(map[string]interface{}{
		"providers": map[string]interface{}{
			"openai": map[string]interface{}{"enabled": true, "apiKey": "sk-test"},
		},
		"niche": map[string]interface{}{"name": "AI Tech"},
	})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if !cfg.Providers.OpenAI.Enabled || cfg.Providers.OpenAI.APIKey != "sk-test" {
		t.Fatalf("openai not updated: %+v", cfg.Providers.OpenAI)
	}
	if cfg.Providers.OpenAI.BaseURL != "https://api.openai.com/v1" {
		t.Fatalf("openai baseurl wiped: %q", cfg.Providers.OpenAI.BaseURL)
	}
	if cfg.Providers.Ollama.BaseURL != "http://localhost:11434" {
		t.Fatalf("ollama baseurl wiped: %q", cfg.Providers.Ollama.BaseURL)
	}
	if cfg.Niche.Name != "AI Tech" {
		t.Fatalf("niche name = %q", cfg.Niche.Name)
	}
}

func TestGetReturnsEqualCopy(t *testing.T) {
	cfg := DefaultConfig()
	got := cfg.Get()
	if got.ID != cfg.ID || got.Version != cfg.Version {
		t.Fatal("Get() diverged from config")
	}
	if got.Providers.Ollama.BaseURL != cfg.Providers.Ollama.BaseURL {
		t.Fatal("Get() providers diverged")
	}
}
