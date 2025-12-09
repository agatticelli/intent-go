# intent-go

Natural Language Processing (NLP) for trading commands. Convert natural language inputs in Spanish or English into structured trading commands that can be executed by trading systems.

## Features

- **Natural Language Understanding**: Parse trading commands in plain language
- **Multi-Language Support**: Spanish and English
- **Intent Classification**: Identify the trading action (open, close, view, etc.)
- **Entity Extraction**: Extract symbols, prices, risk parameters, etc.
- **Parameter Validation**: Validate extracted parameters and report missing fields
- **Wit.ai Integration**: Production-ready NLP backend
- **Extensible**: Easy to add new NLP providers (OpenAI, Anthropic, etc.)

## Dependencies

- **[trading-common-types](https://github.com/agatticelli/trading-common-types)**: Shared type definitions (Side, Intent, NormalizedCommand, etc.)

All types are re-exported for convenience, so you can use `intent.Side` or `types.Side` interchangeably.

## Installation

```bash
go get github.com/agatticelli/intent-go
```

## Quick Start

```go
import (
    "context"
    "fmt"
    "os"

    "github.com/agatticelli/intent-go"
    "github.com/agatticelli/intent-go/witai"
)

// Create Wit.ai processor
processor, err := witai.New(os.Getenv("WIT_AI_TOKEN"))
if err != nil {
    panic(err)
}

// Parse natural language command
input := "open long BTC at 45000 with stop loss 44500 and risk 2%"
cmd, err := processor.ParseCommand(context.Background(), input)
if err != nil {
    panic(err)
}

// Use the structured command
fmt.Printf("Intent: %s\n", cmd.Intent)
fmt.Printf("Symbol: %s\n", cmd.Symbol)
fmt.Printf("Side: %s\n", *cmd.Side)
fmt.Printf("Entry: $%.2f\n", *cmd.EntryPrice)
fmt.Printf("Stop Loss: $%.2f\n", *cmd.StopLoss)
fmt.Printf("Risk: %.1f%%\n", *cmd.RiskPercent)
```

**Output:**
```
Intent: open_position
Symbol: BTC-USDT
Side: LONG
Entry: $45000.00
Stop Loss: $44500.00
Risk: 2.0%
```

For complete working examples, see the [examples/](examples/) directory.

## Architecture

intent-go is part of a 5-module trading system:

```
intent-go (v0.1.0)      → NLP processing (this module)
    ↓
trading-cli             → CLI orchestrator
```

**Key Design Decisions:**

1. **Processor Interface**: Abstract interface for any NLP provider
2. **NormalizedCommand**: Central data structure that flows through the system
3. **Own Type Definitions**: `Intent` and `Side` types independent of other modules
4. **Zero Dependencies**: Uses only Go standard library
5. **Language Agnostic**: Works with Spanish and English seamlessly

## Processor Interface

All NLP providers must implement this interface:

```go
type Processor interface {
    // Parse natural language into structured command
    ParseCommand(ctx context.Context, input string) (*NormalizedCommand, error)

    // Return processor name
    Name() string

    // List supported language codes
    SupportedLanguages() []string
}
```

## Core Types

### NormalizedCommand

The central data structure representing a parsed trading command:

```go
type NormalizedCommand struct {
    // Intent classification
    Intent     Intent    // open_position, close_position, etc.
    Confidence float64   // 0.0 - 1.0

    // Extracted parameters
    Symbol string       // "BTC-USDT", "ETH-USDT"
    Side   *Side        // LONG or SHORT

    // Price parameters
    EntryPrice   *float64
    StopLoss     *float64
    TakeProfit   *float64
    TriggerPrice *float64

    // Multi-level take profits
    TPLevels []TPLevel

    // Risk parameters
    RiskPercent *float64  // 0-100
    RRRatio     *float64  // e.g., 2.0 for 2:1

    // Trailing parameters
    CallbackRate *float64
    Distance     *float64

    // Validation
    Valid   bool
    Missing []string  // Missing required parameters
    Errors  []string  // Validation errors

    // Metadata
    RawInput  string
    Language  string
    Timestamp time.Time
}
```

### Intent

Supported trading actions:

```go
type Intent string

const (
    IntentOpenPosition  Intent = "open_position"
    IntentClosePosition Intent = "close_position"
    IntentViewPositions Intent = "view_positions"
    IntentViewOrders    Intent = "view_orders"
    IntentCancelOrders  Intent = "cancel_orders"
    IntentCheckBalance  Intent = "check_balance"
    IntentBreakEven     Intent = "break_even"
    IntentTrailingStop  Intent = "trailing_stop"
    IntentUnknown       Intent = "unknown"
)
```

### Side

Position direction:

```go
type Side string

const (
    SideLong  Side = "LONG"
    SideShort Side = "SHORT"
)
```

### TPLevel

Multi-level take profit:

```go
type TPLevel struct {
    Price      float64
    Percentage float64  // 0-100
}
```

## Wit.ai Integration

### Setup

1. Create a Wit.ai account: https://wit.ai/
2. Create a new app for trading commands
3. Train intents and entities (or import pre-trained model)
4. Get your Server Access Token

### Creating a Processor

```go
import "github.com/agatticelli/intent-go/witai"

processor, err := witai.New(token)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Processor: %s\n", processor.Name())
fmt.Printf("Languages: %v\n", processor.SupportedLanguages())
```

### Training Data Examples

**English Examples:**
- "open long BTC at 45000 with stop loss 44500 and risk 2%"
- "close 50% of ETH position"
- "set trailing stop on BTC at 1%"
- "show my positions"
- "cancel all orders"

**Spanish Examples:**
- "abrir largo BTC en 45000 con stop loss 44500 y riesgo 2%"
- "cerrar 50% de posición ETH"
- "poner trailing stop en BTC al 1%"
- "mostrar mis posiciones"
- "cancelar todas las órdenes"

## Supported Intents

### open_position

Open a new trading position.

**Required:**
- Symbol
- Side (LONG/SHORT)

**Optional:**
- EntryPrice
- StopLoss
- TakeProfit or RRRatio
- RiskPercent

**Examples:**
```
"open long BTC at 45000 with SL 44500 and TP 46000"
"abrir largo ETH en 3000 con stop 2900 y riesgo 2%"
"buy BTC at 45k, stop 44k, risk 1.5%"
```

### close_position

Close an existing position (full or partial).

**Required:**
- Symbol

**Optional:**
- Percentage (default: 100%)

**Examples:**
```
"close BTC position"
"close 50% of ETH"
"cerrar posición de BTC"
"cerrar 25% de ETH"
```

### trailing_stop

Set a trailing stop loss.

**Required:**
- Symbol
- CallbackRate or Distance

**Examples:**
```
"set trailing stop on BTC at 1%"
"trailing stop ETH 0.5%"
"poner trailing en BTC al 1%"
```

### break_even

Move stop loss to entry price (break-even).

**Required:**
- Symbol

**Examples:**
```
"move BTC to break even"
"break even ETH"
"mover ETH a break even"
```

### view_positions / view_orders / check_balance

View account information.

**Examples:**
```
"show positions"
"show my orders"
"check balance"
"mostrar posiciones"
"ver balance"
```

## Validation

NormalizedCommand includes validation status:

```go
cmd, err := processor.ParseCommand(ctx, input)
if err != nil {
    log.Fatal(err)
}

if !cmd.Valid {
    fmt.Println("Missing parameters:", cmd.Missing)
    fmt.Println("Errors:", cmd.Errors)
    return
}

// Command is valid, proceed with execution
```

## Symbol Normalization

Raw inputs are normalized to exchange format:

| Input | Normalized |
|-------|-----------|
| "BTC" | "BTC-USDT" |
| "ETH" | "ETH-USDT" |
| "bitcoin" | "BTC-USDT" |
| "ethereum" | "ETH-USDT" |

## Error Handling

```go
cmd, err := processor.ParseCommand(ctx, input)
if err != nil {
    // Network error, API error, etc.
    log.Printf("Parse error: %v", err)
    return
}

if cmd.Intent == intent.IntentUnknown {
    fmt.Println("Could not understand command")
    fmt.Printf("Confidence: %.2f\n", cmd.Confidence)
    return
}

if !cmd.Valid {
    fmt.Println("Command incomplete:")
    for _, missing := range cmd.Missing {
        fmt.Printf("  - Missing: %s\n", missing)
    }
    return
}

// Command is valid
```

## Implementing a Custom Processor

To add a new NLP provider:

```go
package mynlp

import (
    "context"
    "github.com/agatticelli/intent-go"
)

type MyProcessor struct {
    apiKey string
}

func New(apiKey string) *MyProcessor {
    return &MyProcessor{apiKey: apiKey}
}

func (p *MyProcessor) Name() string {
    return "mynlp"
}

func (p *MyProcessor) SupportedLanguages() []string {
    return []string{"en", "es", "fr"}
}

func (p *MyProcessor) ParseCommand(ctx context.Context, input string) (*intent.NormalizedCommand, error) {
    // 1. Call your NLP API
    // 2. Extract intent and entities
    // 3. Normalize to intent.NormalizedCommand
    // 4. Validate and return

    return &intent.NormalizedCommand{
        Intent:     intent.IntentOpenPosition,
        Confidence: 0.95,
        Symbol:     "BTC-USDT",
        // ... rest of fields
        Valid: true,
    }, nil
}
```

## Examples

See the [examples/](examples/) directory for complete working code:

- **basic_parsing.go**: Parse simple commands
- **validation.go**: Handle validation and errors
- **multi_language.go**: Spanish and English examples
- **all_intents.go**: Examples of all supported intents

## Dependencies

**None** - Uses only Go standard library:
- `context` - Request cancellation
- `net/http` - Wit.ai API calls
- `encoding/json` - JSON parsing
- `time` - Timestamps

## Testing

```bash
# Run tests
go test ./...

# Test with Wit.ai token
export WIT_AI_TOKEN="your-wit-ai-token"
go run examples/basic_parsing.go
```

## Supported NLP Providers

| Provider | Status | Languages |
|----------|--------|-----------|
| Wit.ai   | ✅ Complete | en, es |
| OpenAI   | 🚧 Planned | Any |
| Anthropic | 🚧 Planned | Any |

## License

MIT
