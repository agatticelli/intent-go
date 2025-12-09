package validators

import (
	"testing"

	"github.com/agatticelli/intent-go"
	"github.com/agatticelli/trading-common-types"
)

func float64Ptr(v float64) *float64 {
	return &v
}

func sidePtr(s types.Side) *types.Side {
	return &s
}

func TestValidateCommand_OpenPosition(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *intent.NormalizedCommand
		wantValid   bool
		wantMissing []string
		wantErrors  []string
	}{
		{
			name: "Valid open position - complete",
			cmd: &intent.NormalizedCommand{
				Intent:      intent.IntentOpenPosition,
				Symbol:      "BTC-USDT",
				Side:        sidePtr(types.SideLong),
				EntryPrice:  float64Ptr(45000.0),
				StopLoss:    float64Ptr(44500.0),
				RiskPercent: float64Ptr(2.0),
			},
			wantValid:   true,
			wantMissing: []string{},
			wantErrors:  []string{},
		},
		{
			name: "Missing symbol",
			cmd: &intent.NormalizedCommand{
				Intent:      intent.IntentOpenPosition,
				Side:        sidePtr(types.SideLong),
				EntryPrice:  float64Ptr(45000.0),
				StopLoss:    float64Ptr(44500.0),
				RiskPercent: float64Ptr(2.0),
			},
			wantValid:   false,
			wantMissing: []string{"symbol"},
		},
		{
			name: "Missing side",
			cmd: &intent.NormalizedCommand{
				Intent:      intent.IntentOpenPosition,
				Symbol:      "BTC-USDT",
				EntryPrice:  float64Ptr(45000.0),
				StopLoss:    float64Ptr(44500.0),
				RiskPercent: float64Ptr(2.0),
			},
			wantValid:   false,
			wantMissing: []string{"side"},
		},
		{
			name: "Missing entry price",
			cmd: &intent.NormalizedCommand{
				Intent:      intent.IntentOpenPosition,
				Symbol:      "BTC-USDT",
				Side:        sidePtr(types.SideLong),
				StopLoss:    float64Ptr(44500.0),
				RiskPercent: float64Ptr(2.0),
			},
			wantValid:   false,
			wantMissing: []string{"entry_price"},
		},
		{
			name: "Missing stop loss",
			cmd: &intent.NormalizedCommand{
				Intent:      intent.IntentOpenPosition,
				Symbol:      "BTC-USDT",
				Side:        sidePtr(types.SideLong),
				EntryPrice:  float64Ptr(45000.0),
				RiskPercent: float64Ptr(2.0),
			},
			wantValid:   false,
			wantMissing: []string{"stop_loss"},
		},
		{
			name: "Missing risk percent",
			cmd: &intent.NormalizedCommand{
				Intent:     intent.IntentOpenPosition,
				Symbol:     "BTC-USDT",
				Side:       sidePtr(types.SideLong),
				EntryPrice: float64Ptr(45000.0),
				StopLoss:   float64Ptr(44500.0),
			},
			wantValid:   false,
			wantMissing: []string{"risk_percent"},
		},
		{
			name: "Invalid risk percent - too high",
			cmd: &intent.NormalizedCommand{
				Intent:      intent.IntentOpenPosition,
				Symbol:      "BTC-USDT",
				Side:        sidePtr(types.SideLong),
				EntryPrice:  float64Ptr(45000.0),
				StopLoss:    float64Ptr(44500.0),
				RiskPercent: float64Ptr(150.0),
			},
			wantValid:  false,
			wantErrors: []string{"risk_percent must be between 0 and 100"},
		},
		{
			name: "Invalid risk percent - zero",
			cmd: &intent.NormalizedCommand{
				Intent:      intent.IntentOpenPosition,
				Symbol:      "BTC-USDT",
				Side:        sidePtr(types.SideLong),
				EntryPrice:  float64Ptr(45000.0),
				StopLoss:    float64Ptr(44500.0),
				RiskPercent: float64Ptr(0.0),
			},
			wantValid:  false,
			wantErrors: []string{"risk_percent must be between 0 and 100"},
		},
		{
			name: "Invalid LONG - SL above entry",
			cmd: &intent.NormalizedCommand{
				Intent:      intent.IntentOpenPosition,
				Symbol:      "BTC-USDT",
				Side:        sidePtr(types.SideLong),
				EntryPrice:  float64Ptr(45000.0),
				StopLoss:    float64Ptr(46000.0),
				RiskPercent: float64Ptr(2.0),
			},
			wantValid:  false,
			wantErrors: []string{"stop_loss must be below entry_price for LONG"},
		},
		{
			name: "Invalid SHORT - SL below entry",
			cmd: &intent.NormalizedCommand{
				Intent:      intent.IntentOpenPosition,
				Symbol:      "ETH-USDT",
				Side:        sidePtr(types.SideShort),
				EntryPrice:  float64Ptr(3000.0),
				StopLoss:    float64Ptr(2900.0),
				RiskPercent: float64Ptr(2.0),
			},
			wantValid:  false,
			wantErrors: []string{"stop_loss must be above entry_price for SHORT"},
		},
		{
			name: "Invalid TP levels - exceed 100%",
			cmd: &intent.NormalizedCommand{
				Intent:      intent.IntentOpenPosition,
				Symbol:      "BTC-USDT",
				Side:        sidePtr(types.SideLong),
				EntryPrice:  float64Ptr(45000.0),
				StopLoss:    float64Ptr(44500.0),
				RiskPercent: float64Ptr(2.0),
				TPLevels: []types.TPLevel{
					{Price: 46000.0, Percentage: 60.0},
					{Price: 47000.0, Percentage: 50.0},
				},
			},
			wantValid:  false,
			wantErrors: []string{"TP percentages sum to 110.0%, cannot exceed 100%"},
		},
		{
			name: "Valid SHORT position",
			cmd: &intent.NormalizedCommand{
				Intent:      intent.IntentOpenPosition,
				Symbol:      "ETH-USDT",
				Side:        sidePtr(types.SideShort),
				EntryPrice:  float64Ptr(3000.0),
				StopLoss:    float64Ptr(3100.0),
				RiskPercent: float64Ptr(2.0),
			},
			wantValid:   true,
			wantMissing: []string{},
			wantErrors:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ValidateCommand(tt.cmd)

			if tt.cmd.Valid != tt.wantValid {
				t.Errorf("Valid = %v, want %v", tt.cmd.Valid, tt.wantValid)
			}

			if len(tt.wantMissing) > 0 {
				if len(tt.cmd.Missing) != len(tt.wantMissing) {
					t.Errorf("Missing = %v, want %v", tt.cmd.Missing, tt.wantMissing)
				} else {
					for i, missing := range tt.wantMissing {
						if tt.cmd.Missing[i] != missing {
							t.Errorf("Missing[%d] = %q, want %q", i, tt.cmd.Missing[i], missing)
						}
					}
				}
			}

			if len(tt.wantErrors) > 0 {
				if len(tt.cmd.Errors) != len(tt.wantErrors) {
					t.Errorf("Errors = %v, want %v", tt.cmd.Errors, tt.wantErrors)
				} else {
					for i, err := range tt.wantErrors {
						if tt.cmd.Errors[i] != err {
							t.Errorf("Errors[%d] = %q, want %q", i, tt.cmd.Errors[i], err)
						}
					}
				}
			}
		})
	}
}

func TestValidateCommand_ClosePosition(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *intent.NormalizedCommand
		wantValid   bool
		wantMissing []string
	}{
		{
			name: "Valid close position",
			cmd: &intent.NormalizedCommand{
				Intent: intent.IntentClosePosition,
				Symbol: "BTC-USDT",
			},
			wantValid:   true,
			wantMissing: []string{},
		},
		{
			name: "Missing symbol",
			cmd: &intent.NormalizedCommand{
				Intent: intent.IntentClosePosition,
			},
			wantValid:   false,
			wantMissing: []string{"symbol"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ValidateCommand(tt.cmd)

			if tt.cmd.Valid != tt.wantValid {
				t.Errorf("Valid = %v, want %v", tt.cmd.Valid, tt.wantValid)
			}

			if len(tt.wantMissing) > 0 && len(tt.cmd.Missing) != len(tt.wantMissing) {
				t.Errorf("Missing = %v, want %v", tt.cmd.Missing, tt.wantMissing)
			}
		})
	}
}

func TestValidateCommand_TrailingStop(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *intent.NormalizedCommand
		wantValid   bool
		wantMissing []string
	}{
		{
			name: "Valid trailing stop with callback rate",
			cmd: &intent.NormalizedCommand{
				Intent:       intent.IntentTrailingStop,
				Symbol:       "BTC-USDT",
				TriggerPrice: float64Ptr(46000.0),
				CallbackRate: float64Ptr(1.0),
			},
			wantValid:   true,
			wantMissing: []string{},
		},
		{
			name: "Valid trailing stop with distance",
			cmd: &intent.NormalizedCommand{
				Intent:       intent.IntentTrailingStop,
				Symbol:       "BTC-USDT",
				TriggerPrice: float64Ptr(46000.0),
				Distance:     float64Ptr(500.0),
			},
			wantValid:   true,
			wantMissing: []string{},
		},
		{
			name: "Missing symbol",
			cmd: &intent.NormalizedCommand{
				Intent:       intent.IntentTrailingStop,
				TriggerPrice: float64Ptr(46000.0),
				CallbackRate: float64Ptr(1.0),
			},
			wantValid:   false,
			wantMissing: []string{"symbol"},
		},
		{
			name: "Missing trigger price",
			cmd: &intent.NormalizedCommand{
				Intent:       intent.IntentTrailingStop,
				Symbol:       "BTC-USDT",
				CallbackRate: float64Ptr(1.0),
			},
			wantValid:   false,
			wantMissing: []string{"trigger_price"},
		},
		{
			name: "Missing callback rate and distance",
			cmd: &intent.NormalizedCommand{
				Intent:       intent.IntentTrailingStop,
				Symbol:       "BTC-USDT",
				TriggerPrice: float64Ptr(46000.0),
			},
			wantValid:   false,
			wantMissing: []string{"callback_rate or distance"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ValidateCommand(tt.cmd)

			if tt.cmd.Valid != tt.wantValid {
				t.Errorf("Valid = %v, want %v", tt.cmd.Valid, tt.wantValid)
			}

			if len(tt.wantMissing) > 0 && len(tt.cmd.Missing) != len(tt.wantMissing) {
				t.Errorf("Missing = %v, want %v", tt.cmd.Missing, tt.wantMissing)
			}
		})
	}
}

func TestValidateCommand_BreakEven(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *intent.NormalizedCommand
		wantValid   bool
		wantMissing []string
	}{
		{
			name: "Valid break even",
			cmd: &intent.NormalizedCommand{
				Intent: intent.IntentBreakEven,
				Symbol: "BTC-USDT",
			},
			wantValid:   true,
			wantMissing: []string{},
		},
		{
			name: "Missing symbol",
			cmd: &intent.NormalizedCommand{
				Intent: intent.IntentBreakEven,
			},
			wantValid:   false,
			wantMissing: []string{"symbol"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ValidateCommand(tt.cmd)

			if tt.cmd.Valid != tt.wantValid {
				t.Errorf("Valid = %v, want %v", tt.cmd.Valid, tt.wantValid)
			}

			if len(tt.wantMissing) > 0 && len(tt.cmd.Missing) != len(tt.wantMissing) {
				t.Errorf("Missing = %v, want %v", tt.cmd.Missing, tt.wantMissing)
			}
		})
	}
}

func TestValidateCommand_ViewIntents(t *testing.T) {
	// View intents don't require validation
	intents := []intent.Intent{
		intent.IntentViewPositions,
		intent.IntentViewOrders,
		intent.IntentCheckBalance,
		intent.IntentCancelOrders,
	}

	for _, intentType := range intents {
		t.Run(string(intentType), func(t *testing.T) {
			cmd := &intent.NormalizedCommand{
				Intent: intentType,
			}

			ValidateCommand(cmd)

			if !cmd.Valid {
				t.Errorf("Valid = false, want true for %s", intentType)
			}
		})
	}
}

func TestValidateCommand_UnknownIntent(t *testing.T) {
	cmd := &intent.NormalizedCommand{
		Intent: intent.IntentUnknown,
	}

	ValidateCommand(cmd)

	if cmd.Valid {
		t.Error("Valid = true, want false for unknown intent")
	}
	if len(cmd.Errors) == 0 {
		t.Error("Expected error for unknown intent")
	}
}
