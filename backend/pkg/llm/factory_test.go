package llm

import (
	"testing"
)

func TestNewProviderAllNames(t *testing.T) {
	if len(ProviderNames) != 13 {
		t.Fatalf("expected 13 providers, got %d", len(ProviderNames))
	}
	for _, name := range ProviderNames {
		p, err := NewProvider(name, ProviderConfig{})
		if err != nil {
			t.Fatalf("NewProvider(%s) failed: %v", name, err)
		}
		if p.Name() != name {
			t.Fatalf("provider Name() = %q, want %q", p.Name(), name)
		}
	}
}

func TestNewProviderUnknown(t *testing.T) {
	if _, err := NewProvider("nope", ProviderConfig{}); err == nil {
		t.Fatal("expected error for unknown provider, got nil")
	}
}
