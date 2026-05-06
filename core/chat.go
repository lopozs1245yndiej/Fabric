package core

import (
	"context"
	"fmt"
	"strings"
)

// Message represents a single message in a chat conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest holds the parameters for a chat completion request.
type ChatRequest struct {
	// Messages is the conversation history.
	Messages []Message `json:"messages"`
	// Model is the LLM model identifier to use.
	Model string `json:"model"`
	// Temperature controls randomness (0.0 - 2.0).
	Temperature float64 `json:"temperature"`
	// MaxTokens limits the response length. Zero means provider default.
	MaxTokens int `json:"max_tokens,omitempty"`
	// SystemPrompt is an optional system-level instruction prepended to the conversation.
	SystemPrompt string `json:"-"`
}

// ChatResponse contains the result of a chat completion request.
type ChatResponse struct {
	// Content is the generated text from the model.
	Content string `json:"content"`
	// Model is the model that produced the response.
	Model string `json:"model"`
	// FinishReason indicates why the model stopped generating.
	FinishReason string `json:"finish_reason"`
}

// ChatProvider defines the interface that any LLM backend must implement.
type ChatProvider interface {
	// Chat sends a request and returns a response.
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
	// Name returns the provider identifier string.
	Name() string
}

// ChatSession manages a stateful conversation with a provider.
type ChatSession struct {
	provider ChatProvider
	config   *Config
	history  []Message
}

// NewChatSession creates a new ChatSession using the given provider and config.
func NewChatSession(provider ChatProvider, cfg *Config) *ChatSession {
	return &ChatSession{
		provider: provider,
		config:   cfg,
		history:  make([]Message, 0),
	}
}

// Send appends the user message to history, calls the provider, and records
// the assistant reply. It returns the assistant's response text.
func (s *ChatSession) Send(ctx context.Context, userMessage string) (string, error) {
	// Trim the message before the empty-check so whitespace-only strings are caught.
	userMessage = strings.TrimSpace(userMessage)
	if userMessage == "" {
		return "", fmt.Errorf("%w: user message must not be empty", ErrInvalidInput)
	}

	s.history = append(s.history, Message{
		Role:    "user",
		Content: userMessage,
	})

	req := ChatRequest{
		Messages:    s.history,
		Model:       s.config.Model,
		Temperature: s.config.Temperature,
	}

	resp, err := s.provider.Chat(ctx, req)
	if err != nil {
		// Remove the last user message so the caller can retry.
		s.history = s.history[:len(s.history)-1]
		return "", fmt.Errorf("chat provider %q: %w", s.provider.Name(), err)
	}

	s.history = append(s.history, Message{
		Role:    "assistant",
		Content: resp.Content,
	})

	return resp.Content, nil
}

// Reset clears the conversation history, starting a fresh session.
func (s *ChatSession) Reset() {
	s.history = make([]Message, 0)
}

// History returns a read-only copy of the current conversation history.
func (s *ChatSession) History() []Message {
	copy := make([]Message, len(s.history))
	for i, m := range s.history {
		copy[i] = m
	}
	return copy
}
