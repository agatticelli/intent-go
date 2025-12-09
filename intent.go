package intent

import "context"

// Processor defines the interface for NLP intent processing
type Processor interface {
	// ParseCommand processes natural language input and returns normalized command
	ParseCommand(ctx context.Context, input string) (*NormalizedCommand, error)

	// Name returns the processor name (e.g., "witai", "openai")
	Name() string

	// SupportedLanguages returns list of supported language codes
	SupportedLanguages() []string
}
