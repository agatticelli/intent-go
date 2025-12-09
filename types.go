package intent

import (
	"github.com/agatticelli/trading-common-types"
)

// Re-export common types for backward compatibility
type (
	Intent            = types.Intent
	Side              = types.Side
	NormalizedCommand = types.NormalizedCommand
	TPLevel           = types.TPLevel
)

// Re-export constants
const (
	IntentOpenPosition  = types.IntentOpenPosition
	IntentClosePosition = types.IntentClosePosition
	IntentViewPositions = types.IntentViewPositions
	IntentViewOrders    = types.IntentViewOrders
	IntentCancelOrders  = types.IntentCancelOrders
	IntentCheckBalance  = types.IntentCheckBalance
	IntentBreakEven     = types.IntentBreakEven
	IntentTrailingStop  = types.IntentTrailingStop
	IntentUnknown       = types.IntentUnknown

	SideLong  = types.SideLong
	SideShort = types.SideShort
)
