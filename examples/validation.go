package main

import (
	"fmt"

	"github.com/agatticelli/intent-go"
)

// This example demonstrates validation handling without requiring Wit.ai
func main() {
	fmt.Println("=== Command Validation Example ===\n")

	// Simulate parsed commands (in real usage, these come from processor.ParseCommand)

	// Example 1: Valid command
	fmt.Println("Example 1: Valid Command")
	validCmd := &intent.NormalizedCommand{
		Intent:     intent.IntentOpenPosition,
		Confidence: 0.95,
		Symbol:     "BTC-USDT",
		Side:       ptrSide(intent.SideLong),
		EntryPrice: ptrFloat(45000.0),
		StopLoss:   ptrFloat(44500.0),
		RiskPercent: ptrFloat(2.0),
		Valid:      true,
		Missing:    []string{},
		Errors:     []string{},
	}
	handleCommand(validCmd)
	fmt.Println()

	// Example 2: Missing required parameters
	fmt.Println("Example 2: Missing Parameters")
	missingCmd := &intent.NormalizedCommand{
		Intent:     intent.IntentOpenPosition,
		Confidence: 0.92,
		Symbol:     "ETH-USDT",
		Side:       ptrSide(intent.SideLong),
		// Missing: EntryPrice, StopLoss, RiskPercent
		Valid:   false,
		Missing: []string{"entry_price", "stop_loss", "risk_percent"},
		Errors:  []string{},
	}
	handleCommand(missingCmd)
	fmt.Println()

	// Example 3: Validation errors
	fmt.Println("Example 3: Validation Errors")
	errorCmd := &intent.NormalizedCommand{
		Intent:     intent.IntentOpenPosition,
		Confidence: 0.88,
		Symbol:     "BTC-USDT",
		Side:       ptrSide(intent.SideLong),
		EntryPrice: ptrFloat(45000.0),
		StopLoss:   ptrFloat(46000.0), // Invalid: SL above entry for LONG
		Valid:      false,
		Missing:    []string{},
		Errors:     []string{"stop_loss must be below entry_price for LONG positions"},
	}
	handleCommand(errorCmd)
	fmt.Println()

	// Example 4: Low confidence
	fmt.Println("Example 4: Low Confidence")
	lowConfidenceCmd := &intent.NormalizedCommand{
		Intent:     intent.IntentUnknown,
		Confidence: 0.45,
		Symbol:     "",
		Valid:      false,
		Missing:    []string{},
		Errors:     []string{"could not understand command"},
		RawInput:   "do something with crypto",
	}
	handleCommand(lowConfidenceCmd)
	fmt.Println()

	fmt.Println("💡 Best Practices:")
	fmt.Println("   - Always check cmd.Valid before executing")
	fmt.Println("   - Check cmd.Intent != IntentUnknown")
	fmt.Println("   - Verify confidence threshold (e.g., > 0.7)")
	fmt.Println("   - Provide helpful error messages for missing parameters")
	fmt.Println("   - Log validation errors for debugging")
}

func handleCommand(cmd *intent.NormalizedCommand) {
	fmt.Printf("Intent: %s (confidence: %.2f)\n", cmd.Intent, cmd.Confidence)

	// Check if intent was recognized
	if cmd.Intent == intent.IntentUnknown {
		fmt.Println("❌ Command not understood")
		if len(cmd.Errors) > 0 {
			fmt.Printf("   Errors: %v\n", cmd.Errors)
		}
		return
	}

	// Check validation status
	if !cmd.Valid {
		fmt.Println("❌ Command is invalid")

		if len(cmd.Missing) > 0 {
			fmt.Println("   Missing parameters:")
			for _, param := range cmd.Missing {
				fmt.Printf("     - %s\n", param)
			}
		}

		if len(cmd.Errors) > 0 {
			fmt.Println("   Validation errors:")
			for _, err := range cmd.Errors {
				fmt.Printf("     - %s\n", err)
			}
		}

		// In a real application, you would prompt the user for missing info
		fmt.Println("   → Prompt user for missing information")
		return
	}

	// Command is valid - proceed with execution
	fmt.Println("✅ Command is valid")
	fmt.Printf("   Symbol: %s\n", cmd.Symbol)
	if cmd.Side != nil {
		fmt.Printf("   Side: %s\n", *cmd.Side)
	}
	fmt.Println("   → Ready to execute")
}

func ptrFloat(f float64) *float64 {
	return &f
}

func ptrSide(s intent.Side) *intent.Side {
	return &s
}
