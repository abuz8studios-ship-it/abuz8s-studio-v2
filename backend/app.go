package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	
	"github.com/abuz8s/studio/internal/agents"
	"github.com/abuz8s/studio/internal/config"
	"github.com/abuz8s/studio/internal/models"
	"github.com/abuz8s/studio/internal/scheduler"
	"github.com/abuz8s/studio/pkg/llm"
	_ "modernc.org/sqlite"
)

// App struct
type App struct {
	ctx       context.Context
	config    *config.Config
	db        *sql.DB
	scheduler *scheduler.Scheduler
	agents    map[agents.AgentType]*agents.Agent
	llmClient llm.Provider
	llm       *llm.Client
	
	// Runtime state
	isRunning bool
	queueSize int
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		agents: make(map[agents.AgentType]*agents.Agent),
	}
}

// Startup is called when the app starts
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		cfg = config.DefaultConfig()
	}
	a.config = cfg
	
	// Initialize database
	if err := a.initDatabase(); err != nil {
		log.Printf("Failed to init database: %v", err)
	}
	
	// Initialize LLM client
	a.initLLMClient()
	
	// Initialize agents
	a.initAgents()
	
	// Initialize scheduler
	a.initScheduler()
	
	log.Println("[App] Startup complete")
}

// Shutdown is called when the app shuts down
func (a *App) Shutdown(ctx context.Context) {
	if a.scheduler != nil {
		a.scheduler.Stop()
	}
	if a.db != nil {
		a.db.Close()
	}
	log.Println("[App] Shutdown complete")
}

func (a *App) initDatabase() error {
	dbPath := a.config.DatabasePath
	
	// Ensure directory exists
	dir := dbPath[:len(dbPath)-len("abuz8s.db")]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	
	a.db = db
	
	// Run migrations
	if err := models.Execute(db); err != nil {
		return err
	}
	
	return nil
}

func (a *App) initLLMClient() {
	// Multi-provider client is never nil; the active single provider may be
	// nil when nothing is enabled — callers must handle that with a clean error.
	a.llm = llm.NewClient(a.config)
	a.llmClient = nil
	for _, name := range a.config.Providers.Priority {
		if p, err := a.llm.GetProvider(name); err == nil {
			a.llmClient = p
			break
		}
	}
	if a.llmClient == nil {
		for _, name := range llm.ProviderNames {
			if p, err := a.llm.GetProvider(name); err == nil {
				a.llmClient = p
				break
			}
		}
	}
	if a.llmClient != nil {
		log.Printf("[App] Active LLM provider: %s", a.llmClient.Name())
	} else {
		log.Println("[App] No LLM provider enabled — generation returns a clean error until one is configured")
	}
}

func (a *App) initAgents() {
	// Create all content agents
	agentTypes := []agents.AgentType{
		agents.AgentIntelCollector,
		agents.AgentScriptWriter,
		agents.AgentXPostGenerator,
		agents.AgentThumbnailForge,
		agents.AgentBlogWriter,
		agents.AgentOutreachEngine,
		agents.AgentNewsletter,
		agents.AgentClipFactory,
		agents.AgentPerformanceEval,
		agents.AgentWeeklyDigest,
	}
	
	providerName := "none"
	if a.llmClient != nil {
		providerName = a.llmClient.Name()
	}
	for _, at := range agentTypes {
		agent := agents.New(at, map[string]interface{}{
			"niche":     a.config.Niche.Name,
			"tone":      a.config.Voice.Tone,
			"style":     a.config.Voice.Style,
			"provider":  providerName,
		})
		
		if a.config.Providers.DefaultModel != "" {
			agent.Model = a.config.Providers.DefaultModel
		}
		
		a.agents[at] = agent
	}
}

func (a *App) initScheduler() {
	a.scheduler = scheduler.New()
	
	// Register default tasks
	defaultTasks := scheduler.GetDefaultTasks()
	
	for _, task := range defaultTasks {
		// Find corresponding agent
		var agentType agents.AgentType
		switch task.ID {
		case scheduler.TaskIntelCollector:
			agentType = agents.AgentIntelCollector
		case scheduler.TaskScriptWriter:
			agentType = agents.AgentScriptWriter
		case scheduler.TaskXPostGenerator:
			agentType = agents.AgentXPostGenerator
		case scheduler.TaskThumbnailForge:
			agentType = agents.AgentThumbnailForge
		case scheduler.TaskBlogWriter:
			agentType = agents.AgentBlogWriter
		case scheduler.TaskOutreachEngine:
			agentType = agents.AgentOutreachEngine
		case scheduler.TaskNewsletter:
			agentType = agents.AgentNewsletter
		case scheduler.TaskClipFactory:
			agentType = agents.AgentClipFactory
		case scheduler.TaskPerformanceEval:
			agentType = agents.AgentPerformanceEval
		case scheduler.TaskWeeklyDigest:
			agentType = agents.AgentWeeklyDigest
		}
		
		// Set handler
		task.Handler = func() error {
			return a.executeAgent(agentType)
		}
		
		if err := a.scheduler.Register(task); err != nil {
			log.Printf("Failed to register task %s: %v", task.ID, err)
		}
	}
	
	a.scheduler.Start()
}

func (a *App) executeAgent(agentType agents.AgentType) error {
	agent, exists := a.agents[agentType]
	if !exists {
		return fmt.Errorf("agent not found: %s", agentType)
	}
	if a.llmClient == nil {
		return fmt.Errorf("no LLM provider configured — enable one in Settings")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	
	// Build input context
	input := map[string]interface{}{
		"niche":     a.config.Niche.Name,
		"audience":  a.config.Niche.TargetAudience,
		"tone":      a.config.Voice.Tone,
		"style":     a.config.Voice.Style,
		"date":      time.Now().Format("2006-01-02"),
	}
	
	result, err := agent.Execute(ctx, &llmAdapter{a.llmClient}, input)
	if err != nil {
		return err
	}
	
	// Save result to database
	if result.Error == "" {
		content := models.NewContent(agentTypeString(agentType), result.Title, result.Body)
		content.Provider = a.llmClient.Name()
		content.Model = agent.Model
		
		store := models.NewContentStore(a.db)
		if err := store.Create(content); err != nil {
			return err
		}
		
		// Emit event to frontend
		runtime.EventsEmit(a.ctx, "content:created", content)
	}
	
	return nil
}

// llmAdapter adapts our LLM provider to the agent interface
type llmAdapter struct {
	provider llm.Provider
}

func (a *llmAdapter) Complete(ctx context.Context, prompt string, opts agents.RequestOptions) (string, error) {
	req := llm.CompletionRequest{
		Model:       opts.Model,
		Temperature: float32(opts.Temperature),
		MaxTokens:   opts.MaxTokens,
		Messages: []llm.Message{
			{Role: "user", Content: prompt},
		},
	}
	
	resp, err := a.provider.Complete(ctx, req)
	if err != nil {
		return "", err
	}
	
	return resp.Content, nil
}

func agentTypeString(at agents.AgentType) string {
	switch at {
	case agents.AgentIntelCollector:
		return models.ContentTypeTrendReport
	case agents.AgentScriptWriter:
		return models.ContentTypeScript
	case agents.AgentXPostGenerator:
		return models.ContentTypeXPost
	case agents.AgentThumbnailForge:
		return models.ContentTypeThumbnail
	case agents.AgentBlogWriter:
		return models.ContentTypeBlogPost
	case agents.AgentOutreachEngine:
		return models.ContentTypeOutreach
	case agents.AgentNewsletter:
		return models.ContentTypeNewsletter
	case agents.AgentClipFactory:
		return models.ContentTypeClipScript
	default:
		return models.ContentTypeIdea
	}
}

// ==================== FRONTEND API ====================

// GetConfig returns the current configuration (safe version without API keys)
func (a *App) GetConfig() map[string]interface{} {
	return a.config.GetPublic()
}

// GetFullConfig returns the full config (for settings page)
func (a *App) GetFullConfig() config.Config {
	return a.config.Get()
}

// UpdateConfig updates the configuration
func (a *App) UpdateConfig(updates map[string]interface{}) error {
	if err := a.config.Update(updates); err != nil {
		return err
	}
	if err := a.config.Save(); err != nil {
		return err
	}
	// Re-pick the active provider so Settings changes apply without restart.
	a.initLLMClient()
	return nil
}

// GetSetupComplete returns whether setup is complete
func (a *App) GetSetupComplete() bool {
	return a.config.SetupComplete
}

// CompleteSetup marks setup as complete
func (a *App) CompleteSetup() error {
	a.config.SetupComplete = true
	return a.config.Save()
}

// GetStatus returns the current app status
func (a *App) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"isRunning":  a.isRunning,
		"queueSize":  a.queueSize,
		"setupComplete": a.config.SetupComplete,
		"agentCount": len(a.agents),
	}
}

// StartAgents starts the content generation pipeline
func (a *App) StartAgents() error {
	a.isRunning = true
	
	// Trigger all agents once
	for _, agent := range a.agents {
		go func(at agents.AgentType) {
			if err := a.executeAgent(at); err != nil {
				log.Printf("Agent %s failed: %v", at, err)
			}
		}(agent.Type)
	}
	
	return nil
}

// StopAgents stops the content generation pipeline
func (a *App) StopAgents() error {
	a.isRunning = false
	return nil
}

// GetAgents returns all agent information
func (a *App) GetAgents() []*agents.Agent {
	var list []*agents.Agent
	for _, agent := range a.agents {
		list = append(list, agent)
	}
	return list
}

// RunAgent manually triggers a specific agent
func (a *App) RunAgent(agentType string) error {
	at := agents.AgentType(agentType)
	return a.executeAgent(at)
}

// GetSchedulerTasks returns all scheduled tasks
func (a *App) GetSchedulerTasks() []*scheduler.Task {
	return a.scheduler.List()
}

// UpdateTaskSchedule updates a task's schedule
func (a *App) UpdateTaskSchedule(taskID, schedule string) error {
	return a.scheduler.UpdateSchedule(taskID, schedule)
}

// EnableTask enables a scheduled task
func (a *App) EnableTask(taskID string) error {
	return a.scheduler.Enable(taskID)
}

// DisableTask disables a scheduled task
func (a *App) DisableTask(taskID string) error {
	return a.scheduler.Disable(taskID)
}

// GetContentList returns a list of content items
func (a *App) GetContentList(contentType, status string, limit, offset int) ([]models.Content, error) {
	store := models.NewContentStore(a.db)
	return store.List(contentType, status, limit, offset)
}

// GetContent returns a single content item
func (a *App) GetContent(id string) (*models.Content, error) {
	store := models.NewContentStore(a.db)
	return store.Get(id)
}

// DeleteContent deletes a content item
func (a *App) DeleteContent(id string) error {
	store := models.NewContentStore(a.db)
	return store.Delete(id)
}

// SearchContent searches content
func (a *App) SearchContent(query string, limit int) ([]models.Content, error) {
	store := models.NewContentStore(a.db)
	return store.Search(query, limit)
}

// GetContentStats returns content statistics
func (a *App) GetContentStats() map[string]interface{} {
	store := models.NewContentStore(a.db)
	
	stats := map[string]interface{}{
		"types": map[string]int{},
		"status": map[string]int{},
	}
	
	// Count by type
	for _, t := range []string{
		models.ContentTypeScript,
		models.ContentTypeThumbnail,
		models.ContentTypeXPost,
		models.ContentTypeBlogPost,
		models.ContentTypeNewsletter,
		models.ContentTypeOutreach,
	} {
		count, _ := store.Count(t, "")
		stats["types"].(map[string]int)[t] = count
	}
	
	// Count by status
	for _, s := range []string{
		models.ContentStatusDraft,
		models.ContentStatusScheduled,
		models.ContentStatusPublished,
	} {
		count, _ := store.Count("", s)
		stats["status"].(map[string]int)[s] = count
	}
	
	total, _ := store.Count("", "")
	stats["total"] = total
	
	return stats
}

// GenerateContent generates content on demand
func (a *App) GenerateContent(req GenerateRequest) (*models.Content, error) {
	agentType := agents.AgentType(req.AgentType)
	
	agent, exists := a.agents[agentType]
	if !exists {
		return nil, fmt.Errorf("agent not found: %s", agentType)
	}
	if a.llmClient == nil {
		return nil, fmt.Errorf("no LLM provider configured — enable one in Settings")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	
	// Build input with request parameters
	input := map[string]interface{}{
		"niche":     a.config.Niche.Name,
		"audience":  a.config.Niche.TargetAudience,
		"tone":      a.config.Voice.Tone,
		"style":     a.config.Voice.Style,
	}
	
	// Add request parameters
	for k, v := range req.Parameters {
		input[k] = v
	}
	
	result, err := agent.Execute(ctx, &llmAdapter{a.llmClient}, input)
	if err != nil {
		return nil, err
	}
	
	content := models.NewContent(req.ContentType, result.Title, result.Body)
	content.Provider = a.llmClient.Name()
	content.Model = agent.Model
	
	store := models.NewContentStore(a.db)
	if err := store.Create(content); err != nil {
		return nil, err
	}
	
	return content, nil
}

// TestProvider tests a provider connection using stored config
func (a *App) TestProvider(provider string) (bool, error) {
	cfg := a.providerConfigByName(provider)
	if cfg == nil {
		return false, fmt.Errorf("unknown provider: %s", provider)
	}
	p, err := llm.NewProvider(provider, llm.ProviderConfig{
		Enabled: true,
		APIKey:  cfg.APIKey,
		BaseURL: cfg.BaseURL,
		Model:   cfg.Model,
	})
	if err != nil {
		return false, err
	}
	return p.Test()
}

func (a *App) providerConfigByName(name string) *config.ProviderConfig {
	p := &a.config.Providers
	switch name {
	case "openrouter":
		return &p.OpenRouter
	case "opencode":
		return &p.OpenCode
	case "openai":
		return &p.OpenAI
	case "codex":
		return &p.Codex
	case "gemini":
		return &p.Gemini
	case "nous":
		return &p.Nous
	case "kimi":
		return &p.Kimi
	case "minimax":
		return &p.MiniMax
	case "glm":
		return &p.GLM
	case "groq":
		return &p.Groq
	case "together":
		return &p.Together
	case "ollama":
		return &p.Ollama
	case "lmstudio":
		return &p.LMStudio
	default:
		return nil
	}
}

// GetProviders returns enabled provider status
func (a *App) GetProviders() map[string]bool {
	return map[string]bool{
		"openrouter": a.config.Providers.OpenRouter.Enabled,
		"opencode":   a.config.Providers.OpenCode.Enabled,
		"openai":     a.config.Providers.OpenAI.Enabled,
		"codex":      a.config.Providers.Codex.Enabled,
		"gemini":     a.config.Providers.Gemini.Enabled,
		"nous":       a.config.Providers.Nous.Enabled,
		"kimi":       a.config.Providers.Kimi.Enabled,
		"minimax":    a.config.Providers.MiniMax.Enabled,
		"glm":        a.config.Providers.GLM.Enabled,
		"groq":       a.config.Providers.Groq.Enabled,
		"together":   a.config.Providers.Together.Enabled,
		"ollama":     a.config.Providers.Ollama.Enabled,
		"lmstudio":   a.config.Providers.LMStudio.Enabled,
	}
}

// GenerateRequest for on-demand content generation
type GenerateRequest struct {
	AgentType  string                 `json:"agentType"`
	ContentType string                `json:"contentType"`
	Parameters map[string]interface{} `json:"parameters"`
}
