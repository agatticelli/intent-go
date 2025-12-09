package witai

import (
	"strconv"
	"strings"
	"time"

	"github.com/agatticelli/intent-go"
)

// transformWitResponse converts Wit.ai response to NormalizedCommand
func transformWitResponse(resp *WitAIResponse, rawInput string) *intent.NormalizedCommand {
	cmd := &intent.NormalizedCommand{
		RawInput:  rawInput,
		Timestamp: time.Now(),
	}

	// Extract intent
	if len(resp.Intents) > 0 {
		cmd.Intent = mapWitIntent(resp.Intents[0].Name)
		cmd.Confidence = resp.Intents[0].Confidence
	}

	// Extract entities
	for entityName, entityValues := range resp.Entities {
		if len(entityValues) == 0 {
			continue
		}

		entity := entityValues[0] // Take highest confidence

		switch entityName {
		case "symbol":
			cmd.Symbol = normalizeSymbol(entity.Value)

		case "side":
			side := normalizeSide(entity.Value)
			cmd.Side = &side

		case "entry_price", "price:entry":
			if price, err := strconv.ParseFloat(entity.Value, 64); err == nil {
				cmd.EntryPrice = &price
			}

		case "stop_loss", "price:stop_loss":
			if sl, err := strconv.ParseFloat(entity.Value, 64); err == nil {
				cmd.StopLoss = &sl
			}

		case "take_profit", "price:take_profit":
			if tp, err := strconv.ParseFloat(entity.Value, 64); err == nil {
				cmd.TakeProfit = &tp
			}

		case "risk":
			if risk, err := strconv.ParseFloat(entity.Value, 64); err == nil {
				cmd.RiskPercent = &risk
			}

		case "trigger_price":
			if trigger, err := strconv.ParseFloat(entity.Value, 64); err == nil {
				cmd.TriggerPrice = &trigger
			}

		case "callback_rate":
			if cb, err := strconv.ParseFloat(entity.Value, 64); err == nil {
				cmd.CallbackRate = &cb
			}

		case "levels":
			// Parse multiple TP levels: "3000:30,3100:70"
			cmd.TPLevels = parseTPLevels(entity.Value)
		}
	}

	return cmd
}

// normalizeSymbol converts various formats to standard "BTC-USDT"
func normalizeSymbol(symbol string) string {
	symbolMap := map[string]string{
		"bitcoin":  "BTC-USDT",
		"btc":      "BTC-USDT",
		"ethereum": "ETH-USDT",
		"eth":      "ETH-USDT",
		"solana":   "SOL-USDT",
		"sol":      "SOL-USDT",
		"bnb":      "BNB-USDT",
		"xrp":      "XRP-USDT",
		"ada":      "ADA-USDT",
		"cardano":  "ADA-USDT",
		"doge":     "DOGE-USDT",
		"dogecoin": "DOGE-USDT",
	}

	normalized := strings.ToLower(strings.TrimSpace(symbol))
	if mapped, ok := symbolMap[normalized]; ok {
		return mapped
	}

	// Assume it's already a symbol, format it
	symbol = strings.ToUpper(symbol)
	if !strings.HasSuffix(symbol, "-USDT") {
		return symbol + "-USDT"
	}
	return symbol
}

// normalizeSide converts various formats to LONG/SHORT
// Supports Spanish and English
func normalizeSide(side string) intent.Side {
	side = strings.ToLower(strings.TrimSpace(side))

	// Long synonyms (English + Spanish)
	longSynonyms := []string{
		"buy", "long", "bullish", "comprar", "largo", "alcista",
	}

	// Short synonyms (English + Spanish)
	shortSynonyms := []string{
		"sell", "short", "bearish", "vender", "corto", "bajista",
	}

	for _, synonym := range longSynonyms {
		if side == synonym {
			return intent.SideLong
		}
	}

	for _, synonym := range shortSynonyms {
		if side == synonym {
			return intent.SideShort
		}
	}

	// Default to LONG if unknown
	return intent.SideLong
}

// mapWitIntent maps Wit.ai intent names to our Intent enum
func mapWitIntent(witIntent string) intent.Intent {
	intentMap := map[string]intent.Intent{
		"open_position":  intent.IntentOpenPosition,
		"close_position": intent.IntentClosePosition,
		"view_positions": intent.IntentViewPositions,
		"view_orders":    intent.IntentViewOrders,
		"cancel_orders":  intent.IntentCancelOrders,
		"check_balance":  intent.IntentCheckBalance,
		"break_even":     intent.IntentBreakEven,
		"trailing_stop":  intent.IntentTrailingStop,
	}

	if mapped, ok := intentMap[witIntent]; ok {
		return mapped
	}

	return intent.IntentUnknown
}

// parseTPLevels parses "3000:30,3100:70" format
func parseTPLevels(input string) []intent.TPLevel {
	var levels []intent.TPLevel

	parts := strings.Split(input, ",")
	for _, part := range parts {
		pricePct := strings.Split(strings.TrimSpace(part), ":")
		if len(pricePct) != 2 {
			continue
		}

		price, err1 := strconv.ParseFloat(pricePct[0], 64)
		pct, err2 := strconv.ParseFloat(pricePct[1], 64)

		if err1 == nil && err2 == nil {
			levels = append(levels, intent.TPLevel{
				Price:      price,
				Percentage: pct,
			})
		}
	}

	return levels
}
