# ChartGaze by TrainFlow AI — Go SDK

**ChartGaze** is the market-intelligence MCP / SDK from **[TrainFlow AI](https://www.trainflow.dev)** — give any AI live multi-TF structure, news, and your trading accounts.

> **Want the AI to trade for you?** ChartGaze is context + accounts.  
> For full autonomous scan → alert → execute on MT4 / MT5 / cTrader / crypto, use **[TrainFlow](https://www.trainflow.dev)** (Nexus / Hunter).  
> Same family: [chartgaze.live](https://chartgaze.live) · [trainflow.dev](https://www.trainflow.dev)

| | |
|---|---|
| Product | ChartGaze by TrainFlow AI |
| Site | https://chartgaze.live |
| Org | https://github.com/TrainFlow-AI |
| Autonomous trading | https://www.trainflow.dev |
| License | MIT |

---

## Overview

The **ChartGaze by TrainFlow AI** Go SDK provides high-performance market intelligence for agents and services. ChartGaze is context + accounts; for full autonomous trading use **[TrainFlow](https://www.trainflow.dev)**. Ideal for:

- **Scalable bots**: Handle thousands of concurrent market queries
- **Trading infrastructure**: Build brokers or fintech platforms on top
- **Microservices**: Integrate ChartGaze into your service mesh
- **WebSocket applications**: Real-time market subscriptions

---

## Installation

```bash
go get github.com/TrainFlow-AI/chartgaze-go-sdk
```

Or add to `go.mod`:

```
require github.com/TrainFlow-AI/chartgaze-go-sdk v0.1.0
```

**Requirements**: Go 1.20+

---

## Quick Start

### 1. Basic Market Query

```go
package main

import (
	"fmt"
	"log"

	"github.com/TrainFlow-AI/chartgaze-go-sdk/chartgaze"
)

func main() {
	// Initialize client
	client, err := chartgaze.NewClient(
		chartgaze.WithAPIKey("cg_pk_live_xxx"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Get market snapshot
	snapshot, err := client.GetMarketSnapshot("EURUSD")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("EURUSD: %.5f (bid: %.5f, ask: %.5f)\n",
		snapshot.Price, snapshot.Bid, snapshot.Ask)

	// Get multi-timeframe context
	context, err := client.GetMarketContext("XAUUSD")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Daily trend: %s\n", context.Daily.Trend)
	fmt.Printf("Hourly structure: %s\n", context.Hourly.Structure)
}
```

### 2. Account Integration

```go
// List connected accounts
accounts, err := client.GetTradingAccounts(ctx)
if err != nil {
	log.Fatal(err)
}

for _, account := range accounts {
	fmt.Printf("%s - %s: %s\n",
		account.Broker, account.Platform, account.Currency)
}

// Get positions
positions, err := client.GetPositions(ctx, accounts[0].ID)
if err != nil {
	log.Fatal(err)
}

for _, pos := range positions {
	fmt.Printf("%s: %s %.2f @ %.5f\n",
		pos.Symbol, pos.Side, pos.Volume, pos.EntryPrice)
}

// Get account info
info, err := client.GetAccountInfo(ctx, accounts[0].ID)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Balance: $%.2f | Equity: $%.2f\n", info.Balance, info.Equity)
```

### 3. Market Analysis

```go
// Historical context
historical, err := client.GetHistoricalContext(ctx, &chartgaze.HistoricalRequest{
	Symbol:            "BTC/USD",
	Timestamp:         time.Parse(time.RFC3339, "2026-09-15T14:30:00Z"),
	LookbackMinutes:   60,
})
if err != nil {
	log.Fatal(err)
}

fmt.Printf("BTC at 14:30: $%.0f\n", historical.Price)
fmt.Printf("Next 30 minutes: %v\n", historical.ForwardContext)

// Compare markets
comparison, err := client.CompareMarkets(ctx, &chartgaze.ComparisonRequest{
	Symbols:      []string{"EURUSD", "GBPUSD", "DXY"},
	AnalysisType: "correlation",
})
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Relationships: %+v\n", comparison.Relationships)

// Get events
events, err := client.GetMarketEvents(ctx, "XAUUSD", 24)
if err != nil {
	log.Fatal(err)
}

for _, event := range events {
	fmt.Printf("%s: %s (impact: %s)\n", event.Time, event.Name, event.Impact)
}
```

### 4. Trading (Optional)

```go
// Propose a trade
proposal, err := client.ProposeTrade(ctx, &chartgaze.TradeRequest{
	AccountID:  "acct_123",
	Symbol:     "EURUSD",
	Side:       "BUY",
	Volume:     1.0,
	StopLoss:   1.0800,
	TakeProfit: 1.0920,
})
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Proposal ID: %s - awaiting approval\n", proposal.ID)

// Execute a trade
trade, err := client.ExecuteTrade(ctx, &chartgaze.TradeRequest{
	AccountID:  "acct_123",
	Symbol:     "EURUSD",
	Side:       "BUY",
	Volume:     1.0,
	OrderType:  "market",
	StopLoss:   1.0800,
	TakeProfit: 1.0920,
})
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Trade executed: %s @ %.5f\n", trade.ID, trade.ExecutionPrice)

// Modify position
err = client.ModifyPosition(ctx, "acct_123", "pos_456", &chartgaze.ModifyRequest{
	StopLoss:   1.0810,
	TakeProfit: 1.0930,
})

// Close position
err = client.ClosePosition(ctx, "acct_123", "pos_456")
```

---

## Authentication

### API Key (Recommended for Bots)

```go
client, err := chartgaze.NewClient(
	chartgaze.WithAPIKey("cg_pk_live_xxx"),
)
```

### OAuth Token (For User-Facing Apps)

```go
import "github.com/TrainFlow-AI/chartgaze-go-sdk/auth"

oauth := auth.NewOAuthClient(
	auth.WithClientID("cg_client_123"),
	auth.WithClientSecret("cg_secret_456"),
)

// Get authorization URL
authURL := oauth.GetAuthorizationURL([]string{
	"market.read",
	"portfolio.read",
	"trading.execute",
})
fmt.Println("Visit:", authURL)

// After user authorizes, exchange code for token
token, err := oauth.ExchangeCode("auth_code_from_callback")
if err != nil {
	log.Fatal(err)
}

client, err := chartgaze.NewClient(
	chartgaze.WithOAuthToken(token),
)
```

---

## Concurrent Operations

Go SDK is designed for concurrency:

```go
package main

import (
	"context"
	"sync"

	"github.com/TrainFlow-AI/chartgaze-go-sdk/chartgaze"
)

func main() {
	client, _ := chartgaze.NewClient(chartgaze.WithAPIKey("cg_pk_live_xxx"))
	ctx := context.Background()

	// Fetch multiple snapshots concurrently
	symbols := []string{"EURUSD", "GBPUSD", "XAUUSD", "BTC/USD"}
	snapshots := make([]*chartgaze.MarketSnapshot, len(symbols))

	var wg sync.WaitGroup
	for i, symbol := range symbols {
		wg.Add(1)
		go func(idx int, sym string) {
			defer wg.Done()
			snapshot, _ := client.GetMarketSnapshot(sym)
			snapshots[idx] = snapshot
		}(i, symbol)
	}

	wg.Wait()

	// Process results
	for _, snapshot := range snapshots {
		fmt.Printf("%s: %.5f\n", snapshot.Symbol, snapshot.Price)
	}
}
```

---

## Error Handling

```go
import "github.com/TrainFlow-AI/chartgaze-go-sdk/errors"

snapshot, err := client.GetMarketSnapshot("INVALID")

if err != nil {
	switch err := err.(type) {
	case *errors.NotFoundError:
		fmt.Println("Symbol not found")
	case *errors.PermissionError:
		fmt.Println("Permission denied")
	case *errors.RateLimitError:
		fmt.Printf("Rate limit. Retry after: %d seconds\n", err.RetryAfter)
	case *errors.ChartGazeError:
		fmt.Printf("Error: %s - %s\n", err.Code, err.Message)
	default:
		fmt.Printf("Unknown error: %v\n", err)
	}
}
```

---

## Advanced: Autonomous Trading Bot

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/TrainFlow-AI/chartgaze-go-sdk/chartgaze"
)

type MarketBot struct {
	client    *chartgaze.Client
	accountID string
}

func NewMarketBot(apiKey, accountID string) (*MarketBot, error) {
	client, err := chartgaze.NewClient(chartgaze.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &MarketBot{client: client, accountID: accountID}, nil
}

func (b *MarketBot) AnalyzeAndTrade(ctx context.Context, symbol string) error {
	// 1. Get market context
	context, err := b.client.GetMarketContext(symbol)
	if err != nil {
		return err
	}

	// 2. Check account
	account, err := b.client.GetAccountInfo(ctx, b.accountID)
	if err != nil {
		return err
	}

	if account.FreeMargin < 1000 {
		fmt.Println("Insufficient margin")
		return nil
	}

	// 3. Decide
	if context.Daily.Trend == "bullish" && context.Hourly.Structure == "higher_high" {
		// 4. Execute
		trade, err := b.client.ExecuteTrade(ctx, &chartgaze.TradeRequest{
			AccountID:  b.accountID,
			Symbol:     symbol,
			Side:       "BUY",
			Volume:     b.calculatePositionSize(account),
			StopLoss:   context.Daily.Support,
			TakeProfit: context.Daily.Resistance,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Trade opened: %s\n", trade.ID)
	}

	// 5. Set alert
	err = b.client.CreateMarketAlert(ctx, &chartgaze.AlertRequest{
		Symbol: symbol,
		Condition: map[string]interface{}{
			"type":  "price_level",
			"level": context.Daily.Resistance,
		},
	})

	return err
}

func (b *MarketBot) MonitorPositions(ctx context.Context) error {
	positions, err := b.client.GetPositions(ctx, b.accountID)
	if err != nil {
		return err
	}

	for _, pos := range positions {
		pnlPercent := ((pos.CurrentPrice - pos.EntryPrice) / pos.EntryPrice) * 100
		fmt.Printf("%s: %.2f%% PnL\n", pos.Symbol, pnlPercent)
	}

	return nil
}

func (b *MarketBot) calculatePositionSize(account *chartgaze.AccountInfo) float64 {
	// Risk management: 2% of account per trade
	return (account.Equity * 0.02) / 100
}

func main() {
	bot, _ := NewMarketBot("cg_pk_live_xxx", "acct_123")
	ctx := context.Background()

	// Analyze multiple symbols
	symbols := []string{"EURUSD", "XAUUSD", "GBPUSD"}
	for _, symbol := range symbols {
		if err := bot.AnalyzeAndTrade(ctx, symbol); err != nil {
			log.Printf("Error analyzing %s: %v\n", symbol, err)
		}
	}

	// Monitor positions
	bot.MonitorPositions(ctx)
}
```

---

## Data Types

### MarketSnapshot
```go
type MarketSnapshot struct {
	Symbol           string    `json:"symbol"`
	Price            float64   `json:"price"`
	Bid              float64   `json:"bid"`
	Ask              float64   `json:"ask"`
	Spread           float64   `json:"spread"`
	DailyOpen        float64   `json:"daily_open"`
	DailyHigh        float64   `json:"daily_high"`
	DailyLow         float64   `json:"daily_low"`
	PrevClose        float64   `json:"prev_close"`
	ChangePercent    float64   `json:"change_percent"`
	Volume           int64     `json:"volume"`
	RelativeVolume   float64   `json:"relative_volume"`
	Volatility       float64   `json:"volatility"`
	Session          string    `json:"session"`
	Timestamp        time.Time `json:"timestamp"`
	Venue            string    `json:"venue"`
}
```

### Position
```go
type Position struct {
	ID            string    `json:"id"`
	Symbol        string    `json:"symbol"`
	Side          string    `json:"side"`
	Volume        float64   `json:"volume"`
	EntryPrice    float64   `json:"entry_price"`
	CurrentPrice  float64   `json:"current_price"`
	PnL           float64   `json:"pnl"`
	PnLPercent    float64   `json:"pnl_percent"`
	OpenedAt      time.Time `json:"opened_at"`
	StopLoss      float64   `json:"stop_loss"`
	TakeProfit    float64   `json:"take_profit"`
}
```

---

## Configuration

### Environment Variables

```bash
export CHARTGAZE_API_KEY="cg_pk_live_xxx"
export CHARTGAZE_ENV="production"
export CHARTGAZE_TIMEOUT="30"
```

### Programmatic Configuration

```go
client, err := chartgaze.NewClient(
	chartgaze.WithAPIKey("cg_pk_live_xxx"),
	chartgaze.WithEnvironment("production"),
	chartgaze.WithTimeout(30 * time.Second),
	chartgaze.WithRetries(3),
	chartgaze.WithDebug(true),
)
```

---

## Testing

```go
package main

import (
	"context"
	"testing"

	"github.com/TrainFlow-AI/chartgaze-go-sdk/mock"
)

func TestMarketSnapshot(t *testing.T) {
	client := mock.NewMockClient()

	snapshot, err := client.GetMarketSnapshot("EURUSD")
	if err != nil {
		t.Fatal(err)
	}

	if snapshot.Symbol != "EURUSD" {
		t.Errorf("Expected EURUSD, got %s", snapshot.Symbol)
	}

	if snapshot.Price <= 0 {
		t.Error("Price should be positive")
	}
}
```

---

## Performance Tips

### 1. Reuse Client
```go
// Good: Reuse client
client, _ := chartgaze.NewClient(chartgaze.WithAPIKey("..."))
snapshot1, _ := client.GetMarketSnapshot("EURUSD")
snapshot2, _ := client.GetMarketSnapshot("GBPUSD")

// Bad: Create new client for each call
client1, _ := chartgaze.NewClient(chartgaze.WithAPIKey("..."))
snapshot1, _ := client1.GetMarketSnapshot("EURUSD")
client2, _ := chartgaze.NewClient(chartgaze.WithAPIKey("..."))
snapshot2, _ := client2.GetMarketSnapshot("GBPUSD")
```

### 2. Batch Operations
```go
// Use concurrent requests
var wg sync.WaitGroup
for _, symbol := range symbols {
	wg.Add(1)
	go func(sym string) {
		defer wg.Done()
		client.GetMarketSnapshot(sym)
	}(symbol)
}
wg.Wait()
```

### 3. Context Timeouts
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

snapshot, err := client.GetMarketSnapshot(ctx, "EURUSD")
```

---

## Contributing

ChartGaze Go SDK is open-source:

```bash
git clone https://github.com/TrainFlow-AI/chartgaze-go-sdk.git
cd go-sdk

# Install dev deps
go get -t ./...

# Run tests
go test ./...

# Submit PR
```

---

## Troubleshooting

### "API key not found"
```go
// Ensure CHARTGAZE_API_KEY is set or passed to NewClient
client, err := chartgaze.NewClient(chartgaze.WithAPIKey("cg_pk_live_xxx"))
```

### "Connection timeout"
```go
// Increase timeout
client, _ := chartgaze.NewClient(
	chartgaze.WithAPIKey("..."),
	chartgaze.WithTimeout(60 * time.Second),
)
```

---

## Support

- **Docs**: https://docs.chartgaze.dev/go
- **GitHub**: https://github.com/TrainFlow-AI/chartgaze-go-sdk
- **Issues**: https://github.com/TrainFlow-AI/chartgaze-go-sdk/issues
- **Email**: dev@chartgaze.dev

---

**SDK Version**: 0.1  
**Last Updated**: 2026-09-29  
**License**: MIT

### TrainFlow AI
- **ChartGaze**: https://chartgaze.live
- **TrainFlow (autonomous trading)**: https://www.trainflow.dev
- **Org**: https://github.com/TrainFlow-AI
