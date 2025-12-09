package main

import (
	"context"
	"fmt"
	"os"

	"github.com/agatticelli/intent-go"
	"github.com/agatticelli/intent-go/witai"
)

// This example demonstrates basic command parsing with Wit.ai
func main() {
	fmt.Println("=== Basic NLP Parsing Example ===\n")

	// Get Wit.ai token from environment
	token := os.Getenv("WIT_AI_TOKEN")
	if token == "" {
		fmt.Println("⚠️  WIT_AI_TOKEN not found in environment")
		fmt.Println("   Get your token from: https://wit.ai/")
		fmt.Println()
		fmt.Println("Example usage:")
		fmt.Println("  export WIT_AI_TOKEN=\"your-token\"")
		fmt.Println("  go run basic_parsing.go")
		fmt.Println()
		showExampleCommands()
		return
	}

	// Create Wit.ai processor
	fmt.Println("Creating Wit.ai processor...")
	processor, err := witai.New(token)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("✅ Processor: %s\n", processor.Name())
	fmt.Printf("   Languages: %v\n\n", processor.SupportedLanguages())

	ctx := context.Background()

	// Example commands to parse
	commands := []string{
		"open long BTC at 45000 with stop loss 44500 and risk 2%",
		"close 50% of ETH position",
		"show my positions",
		"abrir largo ETH en 3000 con stop 2900 y riesgo 1.5%",
	}

	for i, input := range commands {
		fmt.Printf("Command %d: %s\n", i+1, input)
		fmt.Println(string(make([]byte, len(input)+11)))

		cmd, err := processor.ParseCommand(ctx, input)
		if err != nil {
			fmt.Printf("❌ Error: %v\n\n", err)
			continue
		}

		printCommand(cmd)
		fmt.Println()
	}
}

func printCommand(cmd *intent.NormalizedCommand) {
	fmt.Printf("✅ Parsed Successfully\n")
	fmt.Printf("   Intent: %s\n", cmd.Intent)
	fmt.Printf("   Confidence: %.2f\n", cmd.Confidence)
	fmt.Printf("   Language: %s\n", cmd.Language)

	if cmd.Symbol != "" {
		fmt.Printf("   Symbol: %s\n", cmd.Symbol)
	}

	if cmd.Side != nil {
		fmt.Printf("   Side: %s\n", *cmd.Side)
	}

	if cmd.EntryPrice != nil {
		fmt.Printf("   Entry Price: $%.2f\n", *cmd.EntryPrice)
	}

	if cmd.StopLoss != nil {
		fmt.Printf("   Stop Loss: $%.2f\n", *cmd.StopLoss)
	}

	if cmd.TakeProfit != nil {
		fmt.Printf("   Take Profit: $%.2f\n", *cmd.TakeProfit)
	}

	if cmd.RiskPercent != nil {
		fmt.Printf("   Risk: %.1f%%\n", *cmd.RiskPercent)
	}

	if cmd.RRRatio != nil {
		fmt.Printf("   R:R Ratio: %.1f:1\n", *cmd.RRRatio)
	}

	fmt.Printf("   Valid: %v\n", cmd.Valid)

	if !cmd.Valid && len(cmd.Missing) > 0 {
		fmt.Printf("   Missing: %v\n", cmd.Missing)
	}

	if len(cmd.Errors) > 0 {
		fmt.Printf("   Errors: %v\n", cmd.Errors)
	}
}

func showExampleCommands() {
	fmt.Println("📝 Example Commands (English):")
	fmt.Println("   - open long BTC at 45000 with stop loss 44500 and risk 2%")
	fmt.Println("   - close 50% of ETH position")
	fmt.Println("   - set trailing stop on BTC at 1%")
	fmt.Println("   - show my positions")
	fmt.Println("   - cancel all orders")
	fmt.Println()
	fmt.Println("📝 Example Commands (Spanish):")
	fmt.Println("   - abrir largo BTC en 45000 con stop loss 44500 y riesgo 2%")
	fmt.Println("   - cerrar 50% de posición ETH")
	fmt.Println("   - poner trailing stop en BTC al 1%")
	fmt.Println("   - mostrar mis posiciones")
	fmt.Println("   - cancelar todas las órdenes")
}
