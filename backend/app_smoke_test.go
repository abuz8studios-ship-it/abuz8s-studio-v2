package main

import (
	"path/filepath"
	"testing"

	"github.com/abuz8s/studio/internal/agents"
	"github.com/abuz8s/studio/internal/config"
)

// Fresh-install lifecycle: no provider enabled must NOT panic; agents, tasks,
// and database must all come up. (Regression: initAgents called .Name() on a
// nil provider and crashed every fresh start.)
func TestFreshStartupNoProvider(t *testing.T) {
	dir := t.TempDir()

	app := NewApp()
	app.config = config.DefaultConfig()
	app.config.Storage.DataDir = dir
	app.config.DatabasePath = filepath.Join(dir, "abuz8s.db")

	if err := app.initDatabase(); err != nil {
		t.Fatalf("initDatabase failed: %v", err)
	}
	defer app.db.Close()

	app.initLLMClient()
	if app.llmClient != nil {
		t.Fatal("expected nil active provider with all providers disabled")
	}
	if app.llm == nil {
		t.Fatal("expected non-nil multi-provider client")
	}

	app.initAgents()
	if len(app.agents) != 10 {
		t.Fatalf("expected 10 agents, got %d", len(app.agents))
	}

	app.initScheduler()
	defer app.scheduler.Stop()
	tasks := app.scheduler.List()
	if len(tasks) != 10 {
		t.Fatalf("expected 10 scheduled tasks, got %d", len(tasks))
	}
	for _, task := range tasks {
		if !task.Enabled {
			t.Fatalf("task %s not enabled", task.ID)
		}
	}

	// Generation without a provider must be a clean error, not a panic.
	if err := app.executeAgent(agents.AgentScriptWriter); err == nil {
		t.Fatal("expected clean no-provider error, got nil")
	}
	if _, err := app.GenerateContent(GenerateRequest{AgentType: string(agents.AgentXPostGenerator)}); err == nil {
		t.Fatal("expected clean no-provider error from GenerateContent, got nil")
	}

	// Content store round-trip on the fresh database.
	if _, err := app.GetContentList("", "", 10, 0); err != nil {
		t.Fatalf("GetContentList failed: %v", err)
	}
	stats := app.GetContentStats()
	if stats["total"] != 0 {
		t.Fatalf("fresh db total = %v, want 0", stats["total"])
	}
}
