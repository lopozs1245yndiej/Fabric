package core

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	OpenAIAPIKey    string
	AnthropicAPIKey string
	OllamaURL       string
	DefaultModel    string
	PatternsDir     string
	OutputDir       string
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		OllamaURL:    "http://localhost:11434",
		DefaultModel: "gpt-4o-mini", // prefer mini by default — cheaper and fast enough for most tasks
		PatternsDir:  filepath.Join(homeDir, ".config", "fabric", "patterns"),
		OutputDir:    filepath.Join(homeDir, ".config", "fabric", "output"),
	}
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig() (*Config, error) {
	// Attempt to load .env file, ignore error if not present
	_ = godotenv.Load()

	cfg := DefaultConfig()

	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		cfg.OpenAIAPIKey = v
	}
	if v := os.Getenv("ANTHROPIC_API_KEY"); v != "" {
		cfg.AnthropicAPIKey = v
	}
	if v := os.Getenv("OLLAMA_URL"); v != "" {
		cfg.OllamaURL = v
	}
	if v := os.Getenv("DEFAULT_MODEL"); v != "" {
		cfg.DefaultModel = v
	}
	if v := os.Getenv("PATTERNS_DIR"); v != "" {
		cfg.PatternsDir = v
	}
	if v := os.Getenv("OUTPUT_DIR"); v != "" {
		cfg.OutputDir = v
	}

	return cfg, nil
}

// Validate checks that required configuration values are present
func (c *Config) Validate() error {
	if c.OpenAIAPIKey == "" && c.AnthropicAPIKey == "" && c.OllamaURL == "" {
		return ErrNoProviderConfigured
	}
	return nil
}
