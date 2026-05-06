package core

import "errors"

// Sentinel errors for the core package
var (
	ErrNoProviderConfigured = errors.New("no AI provider configured: set OPENAI_API_KEY, ANTHROPIC_API_KEY, or OLLAMA_URL")
	ErrPatternNotFound      = errors.New("pattern not found")
	ErrEmptyInput           = errors.New("input cannot be empty")
	ErrInvalidModel         = errors.New("invalid or unsupported model")
	// ErrContextTooLong is returned when the input exceeds the model's context window
	ErrContextTooLong = errors.New("input exceeds model context window limit")
	// ErrSessionNotFound is returned when a named session cannot be located
	ErrSessionNotFound = errors.New("session not found")
)
