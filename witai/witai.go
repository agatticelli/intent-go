package witai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/agatticelli/intent-go"
	"github.com/agatticelli/intent-go/validators"
)

// Processor implements intent.Processor for Wit.ai
type Processor struct {
	token  string
	client *http.Client
}

// New creates a new Wit.ai NLP processor
func New(token string) (*Processor, error) {
	if token == "" {
		return nil, fmt.Errorf("wit.ai token is required")
	}

	return &Processor{
		token:  token,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Name returns the processor name
func (p *Processor) Name() string {
	return "witai"
}

// SupportedLanguages returns list of supported language codes
func (p *Processor) SupportedLanguages() []string {
	return []string{"en", "es"}
}

// ParseCommand processes natural language input and returns normalized command
func (p *Processor) ParseCommand(ctx context.Context, input string) (*intent.NormalizedCommand, error) {
	// Call Wit.ai API
	witResp, err := p.callWitAI(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("wit.ai call failed: %w", err)
	}

	// Transform Wit.ai response to NormalizedCommand
	cmd := transformWitResponse(witResp, input)

	// Validate the command
	validators.ValidateCommand(cmd)

	return cmd, nil
}

// callWitAI makes HTTP request to Wit.ai API
func (p *Processor) callWitAI(ctx context.Context, input string) (*WitAIResponse, error) {
	apiURL := "https://api.wit.ai/message"
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("v", "20240304")
	q.Add("q", input)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+p.token)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wit.ai returned status %d", resp.StatusCode)
	}

	var witResp WitAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&witResp); err != nil {
		return nil, err
	}

	return &witResp, nil
}
