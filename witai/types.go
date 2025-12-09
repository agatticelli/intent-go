package witai

// WitAIResponse represents the response from Wit.ai API
type WitAIResponse struct {
	Text     string                   `json:"text"`
	Intents  []WitAIIntent            `json:"intents"`
	Entities map[string][]WitAIEntity `json:"entities"`
	Traits   map[string][]interface{} `json:"traits"`
}

// WitAIIntent represents an intent from Wit.ai
type WitAIIntent struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

// WitAIEntity represents an entity from Wit.ai
type WitAIEntity struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Role       string  `json:"role"`
	Start      int     `json:"start"`
	End        int     `json:"end"`
	Body       string  `json:"body"`
	Value      string  `json:"value"`
	Confidence float64 `json:"confidence"`
	Type       string  `json:"type"`
}
