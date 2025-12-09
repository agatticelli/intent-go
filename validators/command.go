package validators

import (
	"fmt"

	"github.com/agatticelli/intent-go"
)

// ValidateCommand validates a NormalizedCommand and populates errors
func ValidateCommand(cmd *intent.NormalizedCommand) {
	cmd.Valid = true
	cmd.Missing = []string{}
	cmd.Errors = []string{}

	switch cmd.Intent {
	case intent.IntentOpenPosition:
		validateOpenPosition(cmd)
	case intent.IntentClosePosition:
		validateClosePosition(cmd)
	case intent.IntentTrailingStop:
		validateTrailingStop(cmd)
	case intent.IntentBreakEven:
		validateBreakEven(cmd)
	case intent.IntentCancelOrders, intent.IntentViewPositions, intent.IntentViewOrders, intent.IntentCheckBalance:
		// These intents don't require validation (optional symbol filter)
	default:
		cmd.Valid = false
		cmd.Errors = append(cmd.Errors, fmt.Sprintf("unknown intent: %s", cmd.Intent))
	}
}

func validateOpenPosition(cmd *intent.NormalizedCommand) {
	// Required: symbol, side, entry price, stop loss, risk
	if cmd.Symbol == "" {
		cmd.Missing = append(cmd.Missing, "symbol")
		cmd.Valid = false
	}
	if cmd.Side == nil {
		cmd.Missing = append(cmd.Missing, "side")
		cmd.Valid = false
	}
	if cmd.EntryPrice == nil {
		cmd.Missing = append(cmd.Missing, "entry_price")
		cmd.Valid = false
	}
	if cmd.StopLoss == nil {
		cmd.Missing = append(cmd.Missing, "stop_loss")
		cmd.Valid = false
	}
	if cmd.RiskPercent == nil {
		cmd.Missing = append(cmd.Missing, "risk_percent")
		cmd.Valid = false
	}

	// Validate ranges
	if cmd.RiskPercent != nil && (*cmd.RiskPercent <= 0 || *cmd.RiskPercent > 100) {
		cmd.Errors = append(cmd.Errors, "risk_percent must be between 0 and 100")
		cmd.Valid = false
	}

	// Validate price logic
	if cmd.Side != nil && cmd.EntryPrice != nil && cmd.StopLoss != nil {
		if *cmd.Side == intent.SideLong && *cmd.StopLoss >= *cmd.EntryPrice {
			cmd.Errors = append(cmd.Errors, "stop_loss must be below entry_price for LONG")
			cmd.Valid = false
		}
		if *cmd.Side == intent.SideShort && *cmd.StopLoss <= *cmd.EntryPrice {
			cmd.Errors = append(cmd.Errors, "stop_loss must be above entry_price for SHORT")
			cmd.Valid = false
		}
	}

	// Validate TP levels
	if len(cmd.TPLevels) > 0 {
		totalPct := 0.0
		for _, tp := range cmd.TPLevels {
			totalPct += tp.Percentage
		}
		if totalPct > 100 {
			cmd.Errors = append(cmd.Errors, fmt.Sprintf("TP percentages sum to %.1f%%, cannot exceed 100%%", totalPct))
			cmd.Valid = false
		}
	}
}

func validateClosePosition(cmd *intent.NormalizedCommand) {
	// Symbol is required
	if cmd.Symbol == "" {
		cmd.Missing = append(cmd.Missing, "symbol")
		cmd.Valid = false
	}
}

func validateTrailingStop(cmd *intent.NormalizedCommand) {
	// Required: symbol, trigger price, callback rate or distance
	if cmd.Symbol == "" {
		cmd.Missing = append(cmd.Missing, "symbol")
		cmd.Valid = false
	}
	if cmd.TriggerPrice == nil {
		cmd.Missing = append(cmd.Missing, "trigger_price")
		cmd.Valid = false
	}
	if cmd.CallbackRate == nil && cmd.Distance == nil {
		cmd.Missing = append(cmd.Missing, "callback_rate or distance")
		cmd.Valid = false
	}
}

func validateBreakEven(cmd *intent.NormalizedCommand) {
	// Symbol is required
	if cmd.Symbol == "" {
		cmd.Missing = append(cmd.Missing, "symbol")
		cmd.Valid = false
	}
}
