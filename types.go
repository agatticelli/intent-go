package intent

import "time"

// NormalizedCommand represents a parsed and normalized trading command
// This is the central data structure that flows through the entire system
type NormalizedCommand struct {
	// Intent classification
	Intent     Intent
	Confidence float64 // 0.0 - 1.0

	// Extracted parameters
	Symbol string // Normalized: "BTC-USDT", "ETH-USDT"
	Side   *Side  // Optional: LONG or SHORT

	// Price parameters
	EntryPrice   *float64 // Entry price
	StopLoss     *float64 // Stop loss price
	TakeProfit   *float64 // Single TP price
	TriggerPrice *float64 // For trailing stops

	// Multi-level TP/SL
	TPLevels []TPLevel // Multiple take profit levels

	// Risk parameters
	RiskPercent *float64 // Risk percentage (0-100)
	RRRatio     *float64 // Risk-reward ratio (e.g., 2.0)

	// Trailing parameters
	CallbackRate *float64 // Trailing callback rate (0.005 = 0.5%)
	Distance     *float64 // Fixed distance for trailing

	// Validation status
	Valid   bool
	Missing []string // List of missing required parameters
	Errors  []string // Validation errors

	// Metadata
	RawInput  string
	Language  string // Detected language
	Timestamp time.Time
}

// Intent represents the trading action to perform
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

// Side represents position direction
type Side string

const (
	SideLong  Side = "LONG"
	SideShort Side = "SHORT"
)

// TPLevel represents a take profit level for partial closing
type TPLevel struct {
	Price      float64
	Percentage float64 // 0-100
}
