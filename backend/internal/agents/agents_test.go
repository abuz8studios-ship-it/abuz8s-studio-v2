package agents

import (
	"context"
	"strings"
	"testing"
)

type fakeLLM struct{ reply string }

func (f *fakeLLM) Complete(_ context.Context, prompt string, _ RequestOptions) (string, error) {
	return f.reply, nil
}

func TestAgentExecuteFillsTemplate(t *testing.T) {
	a := New(AgentScriptWriter, map[string]interface{}{"niche": "AI Tech"})
	res, err := a.Execute(context.Background(), &fakeLLM{reply: "My Title\nBody text"}, map[string]interface{}{
		"niche": "AI Tech", "topic": "agents", "length": "60",
		"audience": "devs", "tone": "pro", "style": "casual",
	})
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if res.Error != "" {
		t.Fatalf("result error: %s", res.Error)
	}
	if res.Title != "My Title" {
		t.Fatalf("title = %q", res.Title)
	}
	if a.Status != AgentStatusIdle {
		t.Fatalf("status = %q, want idle", a.Status)
	}
}

func TestAllAgentPromptsSubstantive(t *testing.T) {
	types := []AgentType{
		AgentIntelCollector, AgentScriptWriter, AgentXPostGenerator,
		AgentThumbnailForge, AgentBlogWriter, AgentOutreachEngine,
		AgentNewsletter, AgentClipFactory, AgentPerformanceEval, AgentWeeklyDigest,
	}
	for _, at := range types {
		a := New(at, nil)
		if len(strings.TrimSpace(a.Prompt)) < 100 {
			t.Fatalf("agent %s prompt too short (%d chars)", at, len(a.Prompt))
		}
		if a.Name == string(at) {
			t.Fatalf("agent %s missing display name", at)
		}
	}
}
