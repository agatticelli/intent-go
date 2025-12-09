# Intent-Go Examples

This directory contains examples demonstrating how to use intent-go for natural language processing of trading commands.

## Prerequisites

### Get a Wit.ai Token

1. Create a Wit.ai account: https://wit.ai/
2. Create a new app for trading commands
3. Train your app with trading intents and entities
4. Get your Server Access Token from Settings

**Export the token:**
```bash
export WIT_AI_TOKEN="your-wit-ai-token"
```

## Running the Examples

### 1. Basic Parsing
Parse natural language commands and see structured output.

```bash
go run basic_parsing.go
```

**What it demonstrates:**
- Creating a Wit.ai processor
- Parsing English and Spanish commands
- Extracting intents and entities
- Viewing parsed command structure

**Sample Output:**
```
Command: open long BTC at 45000 with stop loss 44500 and risk 2%
==================================================================
✅ Parsed Successfully
   Intent: open_position
   Confidence: 0.95
   Language: en
   Symbol: BTC-USDT
   Side: LONG
   Entry Price: $45000.00
   Stop Loss: $44500.00
   Risk: 2.0%
   Valid: true
```

### 2. Validation
Handle validation and error cases (works without Wit.ai token).

```bash
go run validation.go
```

**What it demonstrates:**
- Checking command validity
- Handling missing parameters
- Validation error messages
- Low confidence handling
- Best practices for validation

## Example Commands

### English Commands

**Open Position:**
```
"open long BTC at 45000 with stop loss 44500 and risk 2%"
"buy ETH at 3000, stop 2900, take profit 3200"
"long BTC entry 45k sl 44k rr 2"
```

**Close Position:**
```
"close BTC position"
"close 50% of ETH"
"sell all BTC"
```

**Trailing Stop:**
```
"set trailing stop on BTC at 1%"
"trailing stop ETH 0.5%"
"trail BTC 2%"
```

**View Information:**
```
"show my positions"
"show orders"
"check balance"
```

### Spanish Commands

**Abrir Posición:**
```
"abrir largo BTC en 45000 con stop loss 44500 y riesgo 2%"
"comprar ETH en 3000, stop 2900, take profit 3200"
```

**Cerrar Posición:**
```
"cerrar posición de BTC"
"cerrar 50% de ETH"
"vender todo BTC"
```

**Trailing Stop:**
```
"poner trailing stop en BTC al 1%"
"trailing stop ETH 0.5%"
```

**Ver Información:**
```
"mostrar mis posiciones"
"ver órdenes"
"ver balance"
```

## Integration Patterns

### Basic Usage

```go
import (
    "context"
    "github.com/agatticelli/intent-go/witai"
)

// Create processor
processor, err := witai.New(token)
if err != nil {
    log.Fatal(err)
}

// Parse command
cmd, err := processor.ParseCommand(ctx, "open long BTC at 45000")
if err != nil {
    log.Fatal(err)
}

// Validate
if !cmd.Valid {
    fmt.Println("Missing:", cmd.Missing)
    return
}

// Execute command
executeCommand(cmd)
```

### With CLI Integration

```go
// Interactive trading CLI
for {
    fmt.Print("> ")
    input, _ := reader.ReadString('\n')

    cmd, err := processor.ParseCommand(ctx, input)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        continue
    }

    if cmd.Intent == intent.IntentUnknown {
        fmt.Println("Sorry, I didn't understand that")
        continue
    }

    if !cmd.Valid {
        fmt.Println("Missing information:")
        for _, param := range cmd.Missing {
            fmt.Printf("  - %s\n", param)
        }
        continue
    }

    // Execute the valid command
    result, err := executor.Execute(cmd)
    if err != nil {
        fmt.Printf("Execution error: %v\n", err)
    } else {
        fmt.Println("✅ Command executed successfully")
    }
}
```

### With Voice Integration

```go
// Voice-to-text → NLP → Trading
func handleVoiceCommand(audioFile string) error {
    // 1. Convert speech to text
    text, err := speechToText(audioFile)
    if err != nil {
        return err
    }

    // 2. Parse with intent-go
    cmd, err := processor.ParseCommand(ctx, text)
    if err != nil {
        return err
    }

    // 3. Confirm with user (safety!)
    fmt.Printf("Understood: %s %s %s\n", cmd.Intent, cmd.Symbol, *cmd.Side)
    fmt.Print("Confirm? (y/n): ")

    var confirm string
    fmt.Scan(&confirm)

    if confirm == "y" {
        return executor.Execute(cmd)
    }

    return nil
}
```

## Supported Intents

| Intent | Description | Required Fields |
|--------|-------------|-----------------|
| `open_position` | Open new position | Symbol, Side |
| `close_position` | Close position | Symbol |
| `trailing_stop` | Set trailing stop | Symbol, CallbackRate |
| `break_even` | Move SL to entry | Symbol |
| `view_positions` | Show positions | None |
| `view_orders` | Show orders | None |
| `check_balance` | Show balance | None |
| `cancel_orders` | Cancel orders | Symbol (optional) |

## Common Entities

| Entity | Description | Examples |
|--------|-------------|----------|
| `symbol` | Trading pair | BTC, ETH, bitcoin, ethereum |
| `side` | Position direction | long, short, buy, sell |
| `entry_price` | Entry price | 45000, 45k, 3000 |
| `stop_loss` | Stop loss price | 44500, 44.5k, sl 2900 |
| `take_profit` | Take profit price | 46000, tp 3200 |
| `risk_percent` | Risk percentage | 2%, riesgo 1.5%, risk 3 |
| `rr_ratio` | Risk-reward ratio | 2:1, rr 3, ratio 1.5 |
| `percentage` | Partial close % | 50%, 25 percent |
| `callback_rate` | Trailing % | 1%, 0.5 percent |

## Wit.ai Training Tips

### Create These Intents
- `open_position`
- `close_position`
- `trailing_stop`
- `break_even`
- `view_positions`
- `view_orders`
- `check_balance`
- `cancel_orders`

### Create These Entities
- `symbol` (keywords: BTC, ETH, bitcoin, ethereum, etc.)
- `side` (keywords: long, short, buy, sell, largo, corto, comprar, vender)
- `entry_price` (wit/number)
- `stop_loss` (wit/number)
- `take_profit` (wit/number)
- `risk_percent` (wit/number)
- `rr_ratio` (wit/number)
- `percentage` (wit/number)
- `callback_rate` (wit/number)

### Training Phrases

Add variations for each language:

**English:**
- "open long BTC at 45000 stop 44500 risk 2%"
- "buy BTC 45k sl 44k rr 2"
- "long bitcoin entry 45000 stop loss 44500 tp 46000"

**Spanish:**
- "abrir largo BTC en 45000 stop 44500 riesgo 2%"
- "comprar BTC 45k sl 44k rr 2"
- "largo bitcoin entrada 45000 stop loss 44500 tp 46000"

Add at least 10-20 variations per intent for good accuracy.

## Troubleshooting

### "WIT_AI_TOKEN not found"
- Export the token: `export WIT_AI_TOKEN="your-token"`
- Verify: `echo $WIT_AI_TOKEN`

### Low Confidence (<0.7)
- Add more training phrases to Wit.ai
- Include variations of your commands
- Train both English and Spanish

### Missing Entities
- Check entity names match in Wit.ai
- Verify entity extraction in Wit.ai test panel
- Add more examples for that entity

### Wrong Intent
- Add more training data for that intent
- Check for ambiguous phrases
- Review similar intents in Wit.ai

## Further Reading

- [Main README](../README.md) - Complete API documentation
- [Wit.ai Docs](https://wit.ai/docs) - Official Wit.ai documentation
- [MIGRATION_STATUS.md](../../trading-cli/MIGRATION_STATUS.md) - Architecture overview
