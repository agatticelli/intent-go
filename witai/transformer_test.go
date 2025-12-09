package witai

import (
	"reflect"
	"testing"

	"github.com/agatticelli/intent-go"
	"github.com/agatticelli/trading-common-types"
)

func TestNormalizeSymbol(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Crypto names
		{"Bitcoin lowercase", "bitcoin", "BTC-USDT"},
		{"Bitcoin uppercase", "BITCOIN", "BTC-USDT"},
		{"Ethereum lowercase", "ethereum", "ETH-USDT"},
		{"Solana", "solana", "SOL-USDT"},
		{"Cardano", "cardano", "ADA-USDT"},
		{"Dogecoin", "dogecoin", "DOGE-USDT"},

		// Ticker symbols
		{"BTC lowercase", "btc", "BTC-USDT"},
		{"BTC uppercase", "BTC", "BTC-USDT"},
		{"ETH lowercase", "eth", "ETH-USDT"},
		{"ETH uppercase", "ETH", "ETH-USDT"},
		{"SOL", "sol", "SOL-USDT"},
		{"BNB", "bnb", "BNB-USDT"},
		{"XRP", "xrp", "XRP-USDT"},
		{"ADA", "ada", "ADA-USDT"},
		{"DOGE", "doge", "DOGE-USDT"},

		// Already formatted
		{"Already formatted", "BTC-USDT", "BTC-USDT"},
		{"Lowercase formatted", "btc-usdt", "BTC-USDT"},

		// Unknown symbols
		{"Unknown symbol", "UNKNOWN", "UNKNOWN-USDT"},
		{"Another unknown", "XYZ", "XYZ-USDT"},

		// With whitespace
		{"With spaces", "  btc  ", "BTC-USDT"},
		{"With tabs", "\teth\t", "ETH-USDT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeSymbol(tt.input); got != tt.want {
				t.Errorf("normalizeSymbol(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeSide(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  types.Side
	}{
		// English - LONG
		{"buy", "buy", types.SideLong},
		{"Buy uppercase", "BUY", types.SideLong},
		{"long", "long", types.SideLong},
		{"Long uppercase", "LONG", types.SideLong},
		{"bullish", "bullish", types.SideLong},
		{"Bullish uppercase", "BULLISH", types.SideLong},

		// Spanish - LONG
		{"comprar", "comprar", types.SideLong},
		{"Comprar uppercase", "COMPRAR", types.SideLong},
		{"largo", "largo", types.SideLong},
		{"Largo uppercase", "LARGO", types.SideLong},
		{"alcista", "alcista", types.SideLong},
		{"Alcista uppercase", "ALCISTA", types.SideLong},

		// English - SHORT
		{"sell", "sell", types.SideShort},
		{"Sell uppercase", "SELL", types.SideShort},
		{"short", "short", types.SideShort},
		{"Short uppercase", "SHORT", types.SideShort},
		{"bearish", "bearish", types.SideShort},
		{"Bearish uppercase", "BEARISH", types.SideShort},

		// Spanish - SHORT
		{"vender", "vender", types.SideShort},
		{"Vender uppercase", "VENDER", types.SideShort},
		{"corto", "corto", types.SideShort},
		{"Corto uppercase", "CORTO", types.SideShort},
		{"bajista", "bajista", types.SideShort},
		{"Bajista uppercase", "BAJISTA", types.SideShort},

		// With whitespace
		{"With spaces", "  buy  ", types.SideLong},
		{"With tabs", "\tsell\t", types.SideShort},

		// Unknown defaults to LONG
		{"Unknown", "unknown", types.SideLong},
		{"Empty", "", types.SideLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeSide(tt.input); got != tt.want {
				t.Errorf("normalizeSide(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestMapWitIntent(t *testing.T) {
	tests := []struct {
		name      string
		witIntent string
		want      intent.Intent
	}{
		{"open_position", "open_position", intent.IntentOpenPosition},
		{"close_position", "close_position", intent.IntentClosePosition},
		{"view_positions", "view_positions", intent.IntentViewPositions},
		{"view_orders", "view_orders", intent.IntentViewOrders},
		{"cancel_orders", "cancel_orders", intent.IntentCancelOrders},
		{"check_balance", "check_balance", intent.IntentCheckBalance},
		{"break_even", "break_even", intent.IntentBreakEven},
		{"trailing_stop", "trailing_stop", intent.IntentTrailingStop},
		{"unknown", "unknown_intent", intent.IntentUnknown},
		{"empty", "", intent.IntentUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapWitIntent(tt.witIntent); got != tt.want {
				t.Errorf("mapWitIntent(%q) = %v, want %v", tt.witIntent, got, tt.want)
			}
		})
	}
}

func TestParseTPLevels(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []types.TPLevel
	}{
		{
			name:  "Single TP level",
			input: "46000:100",
			want: []types.TPLevel{
				{Price: 46000.0, Percentage: 100.0},
			},
		},
		{
			name:  "Two TP levels",
			input: "46000:50,47000:50",
			want: []types.TPLevel{
				{Price: 46000.0, Percentage: 50.0},
				{Price: 47000.0, Percentage: 50.0},
			},
		},
		{
			name:  "Three TP levels",
			input: "46000:30,47000:40,48000:30",
			want: []types.TPLevel{
				{Price: 46000.0, Percentage: 30.0},
				{Price: 47000.0, Percentage: 40.0},
				{Price: 48000.0, Percentage: 30.0},
			},
		},
		{
			name:  "With decimal percentages",
			input: "46000:33.33,47000:33.33,48000:33.34",
			want: []types.TPLevel{
				{Price: 46000.0, Percentage: 33.33},
				{Price: 47000.0, Percentage: 33.33},
				{Price: 48000.0, Percentage: 33.34},
			},
		},
		{
			name:  "With whitespace",
			input: " 46000:50 , 47000:50 ",
			want: []types.TPLevel{
				{Price: 46000.0, Percentage: 50.0},
				{Price: 47000.0, Percentage: 50.0},
			},
		},
		{
			name:  "Invalid format - missing colon",
			input: "46000",
			want:  []types.TPLevel{},
		},
		{
			name:  "Invalid format - non-numeric",
			input: "abc:def",
			want:  []types.TPLevel{},
		},
		{
			name:  "Partial invalid",
			input: "46000:50,invalid,47000:50",
			want: []types.TPLevel{
				{Price: 46000.0, Percentage: 50.0},
				{Price: 47000.0, Percentage: 50.0},
			},
		},
		{
			name:  "Empty string",
			input: "",
			want:  []types.TPLevel{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTPLevels(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("parseTPLevels(%q) returned %d levels, want %d", tt.input, len(got), len(tt.want))
			}
			for i := range got {
				if got[i].Price != tt.want[i].Price {
					t.Errorf("Level %d Price = %.2f, want %.2f", i, got[i].Price, tt.want[i].Price)
				}
				if got[i].Percentage != tt.want[i].Percentage {
					t.Errorf("Level %d Percentage = %.2f, want %.2f", i, got[i].Percentage, tt.want[i].Percentage)
				}
			}
		})
	}
}

func TestTransformWitResponse(t *testing.T) {
	tests := []struct {
		name  string
		resp  *WitAIResponse
		input string
		want  *intent.NormalizedCommand
	}{
		{
			name: "Open position command",
			resp: &WitAIResponse{
				Intents: []WitAIIntent{
					{Name: "open_position", Confidence: 0.95},
				},
				Entities: map[string][]WitAIEntity{
					"symbol":      {{Value: "btc"}},
					"side":        {{Value: "long"}},
					"entry_price": {{Value: "45000"}},
					"stop_loss":   {{Value: "44500"}},
					"risk":        {{Value: "2"}},
				},
			},
			input: "open long BTC at 45000 with SL 44500 risk 2%",
			want: &intent.NormalizedCommand{
				Intent:     intent.IntentOpenPosition,
				Confidence: 0.95,
				Symbol:     "BTC-USDT",
				RawInput:   "open long BTC at 45000 with SL 44500 risk 2%",
			},
		},
		{
			name: "Close position command",
			resp: &WitAIResponse{
				Intents: []WitAIIntent{
					{Name: "close_position", Confidence: 0.9},
				},
				Entities: map[string][]WitAIEntity{
					"symbol": {{Value: "ethereum"}},
				},
			},
			input: "close ETH position",
			want: &intent.NormalizedCommand{
				Intent:     intent.IntentClosePosition,
				Confidence: 0.9,
				Symbol:     "ETH-USDT",
				RawInput:   "close ETH position",
			},
		},
		{
			name: "Trailing stop command",
			resp: &WitAIResponse{
				Intents: []WitAIIntent{
					{Name: "trailing_stop", Confidence: 0.92},
				},
				Entities: map[string][]WitAIEntity{
					"symbol":        {{Value: "BTC"}},
					"trigger_price": {{Value: "46000"}},
					"callback_rate": {{Value: "1"}},
				},
			},
			input: "set trailing stop on BTC at 46000 with 1% callback",
			want: &intent.NormalizedCommand{
				Intent:     intent.IntentTrailingStop,
				Confidence: 0.92,
				Symbol:     "BTC-USDT",
				RawInput:   "set trailing stop on BTC at 46000 with 1% callback",
			},
		},
		{
			name: "Unknown intent",
			resp: &WitAIResponse{
				Intents: []WitAIIntent{
					{Name: "unknown", Confidence: 0.3},
				},
				Entities: map[string][]WitAIEntity{},
			},
			input: "what is the meaning of life",
			want: &intent.NormalizedCommand{
				Intent:     intent.IntentUnknown,
				Confidence: 0.3,
				RawInput:   "what is the meaning of life",
			},
		},
		{
			name: "Multiple TP levels",
			resp: &WitAIResponse{
				Intents: []WitAIIntent{
					{Name: "open_position", Confidence: 0.95},
				},
				Entities: map[string][]WitAIEntity{
					"symbol":      {{Value: "btc"}},
					"side":        {{Value: "buy"}},
					"entry_price": {{Value: "45000"}},
					"stop_loss":   {{Value: "44500"}},
					"risk":        {{Value: "2"}},
					"levels":      {{Value: "46000:50,47000:50"}},
				},
			},
			input: "open long BTC at 45000 SL 44500 TP 46000:50,47000:50 risk 2%",
			want: &intent.NormalizedCommand{
				Intent:     intent.IntentOpenPosition,
				Confidence: 0.95,
				Symbol:     "BTC-USDT",
				RawInput:   "open long BTC at 45000 SL 44500 TP 46000:50,47000:50 risk 2%",
				TPLevels: []types.TPLevel{
					{Price: 46000.0, Percentage: 50.0},
					{Price: 47000.0, Percentage: 50.0},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformWitResponse(tt.resp, tt.input)

			if got.Intent != tt.want.Intent {
				t.Errorf("Intent = %v, want %v", got.Intent, tt.want.Intent)
			}
			if got.Confidence != tt.want.Confidence {
				t.Errorf("Confidence = %.2f, want %.2f", got.Confidence, tt.want.Confidence)
			}
			if got.Symbol != tt.want.Symbol {
				t.Errorf("Symbol = %q, want %q", got.Symbol, tt.want.Symbol)
			}
			if got.RawInput != tt.want.RawInput {
				t.Errorf("RawInput = %q, want %q", got.RawInput, tt.want.RawInput)
			}

			// Check TPLevels if present
			if len(tt.want.TPLevels) > 0 {
				if !reflect.DeepEqual(got.TPLevels, tt.want.TPLevels) {
					t.Errorf("TPLevels = %v, want %v", got.TPLevels, tt.want.TPLevels)
				}
			}

			// Check timestamp was set
			if got.Timestamp.IsZero() {
				t.Error("Timestamp was not set")
			}
		})
	}
}

func TestTransformWitResponse_ExtractEntities(t *testing.T) {
	// Test detailed entity extraction
	resp := &WitAIResponse{
		Intents: []WitAIIntent{
			{Name: "open_position", Confidence: 0.95},
		},
		Entities: map[string][]WitAIEntity{
			"symbol":       {{Value: "ethereum"}},
			"side":         {{Value: "vender"}}, // Spanish
			"entry_price":  {{Value: "3000.50"}},
			"stop_loss":    {{Value: "3100.00"}},
			"take_profit":  {{Value: "2850.75"}},
			"risk":         {{Value: "1.5"}},
			"trigger_price": {{Value: "3050"}},
		},
	}

	got := transformWitResponse(resp, "test command")

	// Check all extracted values
	if got.Symbol != "ETH-USDT" {
		t.Errorf("Symbol = %q, want %q", got.Symbol, "ETH-USDT")
	}
	if got.Side == nil || *got.Side != types.SideShort {
		t.Errorf("Side = %v, want SHORT", got.Side)
	}
	if got.EntryPrice == nil || *got.EntryPrice != 3000.50 {
		t.Errorf("EntryPrice = %v, want 3000.50", got.EntryPrice)
	}
	if got.StopLoss == nil || *got.StopLoss != 3100.00 {
		t.Errorf("StopLoss = %v, want 3100.00", got.StopLoss)
	}
	if got.TakeProfit == nil || *got.TakeProfit != 2850.75 {
		t.Errorf("TakeProfit = %v, want 2850.75", got.TakeProfit)
	}
	if got.RiskPercent == nil || *got.RiskPercent != 1.5 {
		t.Errorf("RiskPercent = %v, want 1.5", got.RiskPercent)
	}
	if got.TriggerPrice == nil || *got.TriggerPrice != 3050 {
		t.Errorf("TriggerPrice = %v, want 3050", got.TriggerPrice)
	}
}
