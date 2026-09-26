package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/google/uuid"
)

// Config holds all application configuration
type Config struct {
	mu sync.RWMutex

	ID            string            `json:"id"`
	Version       string            `json:"version"`
	SetupComplete bool              `json:"setupComplete"`
	Niche         NicheConfig       `json:"niche"`
	Voice         VoiceConfig       `json:"voice"`
	Content       ContentConfig     `json:"content"`
	Providers     ProvidersConfig   `json:"providers"`
	Schedule      ScheduleConfig    `json:"schedule"`
	Storage       StorageConfig     `json:"storage"`
	DatabasePath  string            `json:"-"`
}

// NicheConfig holds niche and brand settings
type NicheConfig struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	TargetAudience string   `json:"targetAudience"`
	Keywords       []string `json:"keywords"`
	Competitors    []string `json:"competitors"`
}

// VoiceConfig holds voice profile settings
type VoiceConfig struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Style       string   `json:"style"`
	Tone        string   `json:"tone"`
	Patterns    []string `json:"patterns"`
	Phrases     []string `json:"phrases"`
	Avoid       []string `json:"avoid"`
	SampleText  string   `json:"sampleText"`
}

// ContentConfig holds content generation settings
type ContentConfig struct {
	YouTube    YouTubeConfig    `json:"youtube"`
	XPosts     XPostsConfig     `json:"xPosts"`
	Blog       BlogConfig       `json:"blog"`
	Newsletter NewsletterConfig `json:"newsletter"`
	Thumbnails ThumbnailsConfig `json:"thumbnails"`
	Clips      ClipsConfig      `json:"clips"`
}

type YouTubeConfig struct {
	Enabled       bool     `json:"enabled"`
	ScriptsPerDay int      `json:"scriptsPerDay"`
	ScriptLength  int      `json:"scriptLength"`
	Topics        []string `json:"topics"`
}

type XPostsConfig struct {
	Enabled     bool     `json:"enabled"`
	PostsPerDay int      `json:"postsPerDay"`
	Style       string   `json:"style"`
}

type BlogConfig struct {
	Enabled      bool     `json:"enabled"`
	PostsPerWeek int      `json:"postsPerWeek"`
	SEOEnabled   bool     `json:"seoEnabled"`
}

type NewsletterConfig struct {
	Enabled   bool     `json:"enabled"`
	Frequency string   `json:"frequency"`
	Segments  []string `json:"segments"`
}

type ThumbnailsConfig struct {
	Enabled          bool     `json:"enabled"`
	VariantsPerVideo int      `json:"variantsPerVideo"`
	Style            string   `json:"style"`
}

type ClipsConfig struct {
	Enabled     bool     `json:"enabled"`
	ClipsPerDay int      `json:"clipsPerDay"`
	Platforms   []string `json:"platforms"`
}

// ProvidersConfig holds all LLM provider configurations
type ProvidersConfig struct {
	Priority     []string          `json:"priority"`
	DefaultModel string            `json:"defaultModel"`
	
	// Cloud Providers
	OpenRouter   ProviderConfig    `json:"openrouter"`
	OpenCode     ProviderConfig    `json:"opencode"`
	OpenAI       ProviderConfig    `json:"openai"`
	Codex        ProviderConfig    `json:"codex"`
	Gemini       ProviderConfig    `json:"gemini"`
	Nous         ProviderConfig    `json:"nous"`
	Kimi         ProviderConfig    `json:"kimi"`
	MiniMax      ProviderConfig    `json:"minimax"`
	GLM          ProviderConfig    `json:"glm"`
	Groq         ProviderConfig    `json:"groq"`
	Together     ProviderConfig    `json:"together"`
	
	// Local Providers
	Ollama       ProviderConfig    `json:"ollama"`
	LMStudio     ProviderConfig    `json:"lmstudio"`
}

// ProviderConfig for individual providers
type ProviderConfig struct {
	Enabled bool   `json:"enabled"`
	APIKey  string `json:"apiKey"`
	BaseURL string `json:"baseUrl"`
	Model   string `json:"model"`
	Timeout int    `json:"timeout"`
}

// ScheduleConfig holds scheduling settings
type ScheduleConfig struct {
	Timezone string `json:"timezone"`
	Times    struct {
		Intel      string `json:"intel"`
		Scripts    string `json:"scripts"`
		Posts      string `json:"posts"`
		Thumbnails string `json:"thumbnails"`
		Blog       string `json:"blog"`
		Outreach   string `json:"outreach"`
		Newsletter string `json:"newsletter"`
		Clips      string `json:"clips"`
	} `json:"times"`
}

// StorageConfig holds storage settings
type StorageConfig struct {
	DataDir       string `json:"dataDir"`
	OutputDir     string `json:"outputDir"`
	CloudSync     bool   `json:"cloudSync"`
	CloudProvider string `json:"cloudProvider"`
}

// DefaultConfig returns a new default configuration
func DefaultConfig() *Config {
	cfg := &Config{
		ID:            uuid.New().String(),
		Version:       "2.0.0",
		SetupComplete: false,
		Niche: NicheConfig{
			Keywords:    []string{},
			Competitors: []string{},
		},
		Voice: VoiceConfig{
			ID:       uuid.New().String(),
			Name:     "Default",
			Style:    "conversational",
			Tone:     "professional",
			Patterns: []string{},
			Phrases:  []string{},
			Avoid:    []string{},
		},
		Content: ContentConfig{
			YouTube: YouTubeConfig{
				Enabled:       true,
				ScriptsPerDay: 5,
				ScriptLength:  1800,
				Topics:        []string{},
			},
			XPosts: XPostsConfig{
				Enabled:     true,
				PostsPerDay: 5,
				Style:       "hook-first",
			},
			Blog: BlogConfig{
				Enabled:      true,
				PostsPerWeek: 7,
				SEOEnabled:   true,
			},
			Newsletter: NewsletterConfig{
				Enabled:   true,
				Frequency: "daily",
				Segments:  []string{"main"},
			},
			Thumbnails: ThumbnailsConfig{
				Enabled:          true,
				VariantsPerVideo: 3,
				Style:            "high-contrast",
			},
			Clips: ClipsConfig{
				Enabled:     true,
				ClipsPerDay: 3,
				Platforms:   []string{"tiktok", "instagram", "youtube-shorts"},
			},
		},
		Providers: ProvidersConfig{
			Priority:     []string{"ollama", "lmstudio", "openrouter", "openai", "anthropic"},
			DefaultModel: "gpt-4-turbo-preview",
			OpenRouter: ProviderConfig{
				Enabled: false,
				BaseURL: "https://openrouter.ai/api/v1",
			},
			OpenCode: ProviderConfig{
				Enabled: false,
				BaseURL: "https://api.opencode.ai/v1",
			},
			OpenAI: ProviderConfig{
				Enabled: false,
				BaseURL: "https://api.openai.com/v1",
			},
			Codex: ProviderConfig{
				Enabled: false,
				BaseURL: "https://api.openai.com/v1",
			},
			Gemini: ProviderConfig{
				Enabled: false,
				BaseURL: "https://generativelanguage.googleapis.com/v1",
			},
			Nous: ProviderConfig{
				Enabled: false,
				BaseURL: "https://api.nousresearch.com/v1",
			},
			Kimi: ProviderConfig{
				Enabled: false,
				BaseURL: "https://api.moonshot.cn/v1",
			},
			MiniMax: ProviderConfig{
				Enabled: false,
				BaseURL: "https://api.minimax.chat/v1",
			},
			GLM: ProviderConfig{
				Enabled: false,
				BaseURL: "https://open.bigmodel.cn/api/paas/v4",
			},
			Groq: ProviderConfig{
				Enabled: false,
				BaseURL: "https://api.groq.com/openai/v1",
			},
			Together: ProviderConfig{
				Enabled: false,
				BaseURL: "https://api.together.xyz/v1",
			},
			Ollama: ProviderConfig{
				Enabled: false,
				BaseURL: "http://localhost:11434",
				Model:   "llama3",
			},
			LMStudio: ProviderConfig{
				Enabled: false,
				BaseURL: "http://localhost:1234/v1",
			},
		},
		Schedule: ScheduleConfig{
			Timezone: "America/New_York",
			Times: struct {
				Intel      string `json:"intel"`
				Scripts    string `json:"scripts"`
				Posts      string `json:"posts"`
				Thumbnails string `json:"thumbnails"`
				Blog       string `json:"blog"`
				Outreach   string `json:"outreach"`
				Newsletter string `json:"newsletter"`
				Clips      string `json:"clips"`
			}{
				Intel:      "0 5 * * *",
				Scripts:    "0 6 * * *",
				Posts:      "0 6 * * *",
				Thumbnails: "0 7 * * *",
				Blog:       "0 7 * * *",
				Outreach:   "0 7 * * *",
				Newsletter: "0 8 * * *",
				Clips:      "0 */4 * * *",
			},
		},
		Storage: StorageConfig{
			DataDir:   getDefaultDataDir(),
			OutputDir: getDefaultOutputDir(),
			CloudSync: false,
		},
	}
	
	cfg.DatabasePath = filepath.Join(cfg.Storage.DataDir, "abuz8s.db")
	
	return cfg
}

// Load loads configuration from file
func Load() (*Config, error) {
	configPath := getConfigPath()
	
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := DefaultConfig()
		if err := cfg.Save(); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	
	cfg.DatabasePath = filepath.Join(cfg.Storage.DataDir, "abuz8s.db")
	
	return &cfg, nil
}

// Save saves configuration to file
func (c *Config) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	configPath := getConfigPath()
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(configPath, data, 0600)
}

// Update updates configuration with new values
func (c *Config) Update(updates map[string]interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	data, _ := json.Marshal(c)
	var current map[string]interface{}
	json.Unmarshal(data, &current)
	
	for k, v := range updates {
		current[k] = v
	}
	
	newData, _ := json.Marshal(current)
	return json.Unmarshal(newData, c)
}

// GetPublic returns a safe version of config (no API keys)
func (c *Config) GetPublic() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return map[string]interface{}{
		"id":            c.ID,
		"version":       c.Version,
		"setupComplete": c.SetupComplete,
		"niche":         c.Niche,
		"voice":         c.Voice,
		"content":       c.Content,
		"schedule":      c.Schedule,
		"providers": map[string]interface{}{
			"openrouter": c.Providers.OpenRouter.Enabled,
			"opencode":   c.Providers.OpenCode.Enabled,
			"openai":     c.Providers.OpenAI.Enabled,
			"codex":      c.Providers.Codex.Enabled,
			"gemini":     c.Providers.Gemini.Enabled,
			"nous":       c.Providers.Nous.Enabled,
			"kimi":       c.Providers.Kimi.Enabled,
			"minimax":    c.Providers.MiniMax.Enabled,
			"glm":        c.Providers.GLM.Enabled,
			"groq":       c.Providers.Groq.Enabled,
			"together":   c.Providers.Together.Enabled,
			"ollama":     c.Providers.Ollama.Enabled,
			"lmstudio":   c.Providers.LMStudio.Enabled,
		},
	}
}

// Get returns a copy of the config (field-by-field: the mutex must not be copied)
func (c *Config) Get() Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Config{
		ID:            c.ID,
		Version:       c.Version,
		SetupComplete: c.SetupComplete,
		Niche:         c.Niche,
		Voice:         c.Voice,
		Content:       c.Content,
		Providers:     c.Providers,
		Schedule:      c.Schedule,
		Storage:       c.Storage,
		DatabasePath:  c.DatabasePath,
	}
}

func getConfigPath() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "ABUZ8sStudio", "config.json")
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "ABUZ8sStudio", "config.json")
	default:
		return filepath.Join(os.Getenv("HOME"), ".config", "abuz8s-studio", "config.json")
	}
}

func getDefaultDataDir() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "ABUZ8sStudio", "data")
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "ABUZ8sStudio", "data")
	default:
		return filepath.Join(os.Getenv("HOME"), ".local", "share", "abuz8s-studio")
	}
}

func getDefaultOutputDir() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("USERPROFILE"), "Documents", "ABUZ8sStudio")
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Documents", "ABUZ8sStudio")
	default:
		return filepath.Join(os.Getenv("HOME"), "Documents", "ABUZ8sStudio")
	}
}
