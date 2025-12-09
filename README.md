# intent-go

Natural language processing for trading commands.

## Features

- Convert natural language to structured commands
- Support for multiple NLP providers (Wit.ai, OpenAI, etc.)
- Spanish and English support
- Intent classification and entity extraction
- Parameter validation

## Installation

```bash
go get github.com/gattimassimo/intent-go
```

## Usage

```go
import (
    "context"
    "os"
    "github.com/gattimassimo/intent-go"
    "github.com/gattimassimo/intent-go/witai"
)

// Create Wit.ai processor
processor, err := witai.New(os.Getenv("WIT_AI_TOKEN"))

// Parse natural language
cmd, err := processor.ParseCommand(context.Background(),
    "open long ETH at 3950 with stop loss 3900 and risk 2%")

// cmd contains:
// Intent: open_position
// Symbol: "ETH-USDT"
// Side: LONG
// EntryPrice: 3950
// StopLoss: 3900
// RiskPercent: 2.0
// Valid: true
```

## Supported Languages

- English: "open long BTC at 45000..."
- Spanish: "abrir largo BTC en 45000..."

## Supported NLP Providers

- [x] Wit.ai
- [ ] OpenAI (planned)
- [ ] Anthropic (planned)

## NormalizedCommand Structure

The central data structure that flows through the trading system:

```go
type NormalizedCommand struct {
    Intent       Intent      // open_position, close_position, etc.
    Confidence   float64     // 0.0 - 1.0
    Symbol       string      // "BTC-USDT", "ETH-USDT"
    Side         *Side       // LONG or SHORT
    EntryPrice   *float64
    StopLoss     *float64
    RiskPercent  *float64
    Valid        bool
    Missing      []string
    Errors       []string
}
```

## License

MIT
