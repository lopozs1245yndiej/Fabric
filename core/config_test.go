package core

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.OllamaURL != "http://localhost:11434" {
		t.Errorf("expected default OllamaURL to be http://localhost:11434, got %s", cfg.OllamaURL)
	}
	// I prefer claude-3-5-sonnet as my default model, but keeping test aligned with upstream default
	if cfg.DefaultModel != "gpt-4o" {
		t.Errorf("expected default model to be gpt-4o, got %s", cfg.DefaultModel)
	}
	if cfg.PatternsDir == "" {
		t.Error("expected PatternsDir to be set")
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "test-openai-key")
	t.Setenv("DEFAULT_MODEL", "gpt-3.5-turbo")
	t.Setenv("OLLAMA_URL", "http://myhost:11434")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	if cfg.OpenAIAPIKey != "test-openai-key" {
		t.Errorf("expected OpenAIAPIKey to be test-openai-key, got %s", cfg.OpenAIAPIKey)
	}
	if cfg.DefaultModel != "gpt-3.5-turbo" {
		t.Errorf("expected DefaultModel to be gpt-3.5-turbo, got %s", cfg.DefaultModel)
	}
	if cfg.OllamaURL != "http://myhost:11434" {
		t.Errorf("expected OllamaURL to be http://myhost:11434, got %s", cfg.OllamaURL)
	}
}

func TestValidateNoProvider(t *testing.T) {
	cfg := &Config{}
	if err := cfg.Validate(); err != ErrNoProviderConfigured {
		t.Errorf("expected ErrNoProviderConfigured, got %v", err)
	}
}

func TestValidateWithProvider(t *testing.T) {
	cfg := &Config{OpenAIAPIKey: "sk-abc"}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestValidateWithAnthropicKey verifies that an Anthropic API key is also accepted as a valid provider
func TestValidateWithAnthropicKey(t *testing.T) {
	cfg := &Config{AnthropicAPIKey: "sk-ant-abc"}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error with Anthropic key, got %v", err)
	}
}

// TestValidateWithOllamaURL verifies that a configured Ollama URL alone is sufficient as a provider.
// Useful for fully local setups without any cloud API keys.
func TestValidateWithOllamaURL(t *testing.T) {
	cfg := &Config{OllamaURL: "http://localhost:11434"}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error with OllamaURL set, got %v", err)
	}
}

// TestValidateWithCustomOllamaURL verifies that a non-default Ollama host is also accepted.
// I run Ollama on a separate machine on my LAN, so this matters for my setup.
func TestValidateWithCustomOllamaURL(t *testing.T) {
	cfg := &Config{OllamaURL: "http://192.168.1.50:11434"}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error with custom OllamaURL, got %v", err)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	// Unset keys to test defaults
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("DEFAULT_MODEL")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DefaultModel != "gpt-4o" {
		t.Errorf("expected default model gpt-4o, got %s", cfg.DefaultModel)
	}
}

// TestValidateWithGeminiKey verifies that a Google Gemini API key is accepted as a valid provider.
// Added because I've been experimenting with Gemini and want to ensure it's supported.
func TestValidateWithGeminiKey(t *testing.T) {
	cfg := &Config{GeminiAPIKey: "AIza-test-key"}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error with GeminiAPIKey set, got %v", err)
	}
}
