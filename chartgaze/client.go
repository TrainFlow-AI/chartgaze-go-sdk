package chartgaze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultBaseURL = "https://api.chartgaze.live"

// Client is the ChartGaze Go SDK entrypoint for market intelligence (MCP + REST).
type Client struct {
	APIKey  string
	BaseURL string
	HTTP    *http.Client
}

type Option func(*Client)

func WithAPIKey(key string) Option {
	return func(c *Client) { c.APIKey = key }
}

func WithBaseURL(u string) Option {
	return func(c *Client) { c.BaseURL = strings.TrimRight(u, "/") }
}

func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) { c.HTTP = h }
}

func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		BaseURL: DefaultBaseURL,
		HTTP:    &http.Client{Timeout: 90 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	if strings.TrimSpace(c.APIKey) == "" {
		return nil, fmt.Errorf("chartgaze: API key required")
	}
	return c, nil
}

func (c *Client) request(method, path string, body any) (json.RawMessage, error) {
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, c.BaseURL+path, rdr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("chartgaze: HTTP %d: %s", res.StatusCode, string(raw))
	}
	if len(bytes.TrimSpace(raw)) == 0 {
		return json.RawMessage("null"), nil
	}
	return json.RawMessage(raw), nil
}

func (c *Client) ListTools() (json.RawMessage, error) {
	return c.request(http.MethodGet, "/mcp/tools", nil)
}

func (c *Client) CallTool(name string, input map[string]any) (json.RawMessage, error) {
	if input == nil {
		input = map[string]any{}
	}
	q := url.QueryEscape(name)
	return c.request(http.MethodPost, "/mcp/call?tool_name="+q, input)
}

func (c *Client) GetMarketSnapshot(symbol string) (json.RawMessage, error) {
	return c.CallTool("get_market_snapshot", map[string]any{"symbol": symbol})
}

func (c *Client) GetMarketContext(symbol string, includeEvents bool) (json.RawMessage, error) {
	return c.CallTool("get_market_context", map[string]any{
		"symbol":         symbol,
		"include_events": includeEvents,
	})
}

func (c *Client) CompareMarkets(symbols []string) (json.RawMessage, error) {
	return c.CallTool("compare_markets", map[string]any{"symbols": symbols})
}

func (c *Client) GetHistoricalContext(symbol, timestamp string) (json.RawMessage, error) {
	return c.CallTool("get_historical_context", map[string]any{
		"symbol":    symbol,
		"timestamp": timestamp,
	})
}

func (c *Client) GetEconomicCalendarWeek() (json.RawMessage, error) {
	return c.CallTool("get_economic_calendar_week", map[string]any{})
}

func (c *Client) GetTradingAccounts() (json.RawMessage, error) {
	return c.CallTool("get_trading_accounts", map[string]any{})
}

// ProposeTrade queues a trade for Activity approval (Review/Live). accountID required.
func (c *Client) ProposeTrade(accountID, symbol, side string, volume float64, stopLoss, takeProfit *float64) (json.RawMessage, error) {
	in := map[string]any{
		"account_id": accountID,
		"symbol":     symbol,
		"side":       side,
		"volume":     volume,
	}
	if stopLoss != nil {
		in["stop_loss"] = *stopLoss
	}
	if takeProfit != nil {
		in["take_profit"] = *takeProfit
	}
	return c.CallTool("propose_trade", in)
}

// ExecuteTrade places immediately — Live mode only. accountID required.
func (c *Client) ExecuteTrade(accountID, symbol, side string, volume float64, stopLoss, takeProfit *float64) (json.RawMessage, error) {
	in := map[string]any{
		"account_id": accountID,
		"symbol":     symbol,
		"side":       side,
		"volume":     volume,
	}
	if stopLoss != nil {
		in["stop_loss"] = *stopLoss
	}
	if takeProfit != nil {
		in["take_profit"] = *takeProfit
	}
	return c.CallTool("execute_trade", in)
}

func (c *Client) GetUsageSummary() (json.RawMessage, error) {
	return c.request(http.MethodGet, "/usage/summary", nil)
}

func (c *Client) ListAccounts() (json.RawMessage, error) {
	return c.request(http.MethodGet, "/accounts/list", nil)
}
