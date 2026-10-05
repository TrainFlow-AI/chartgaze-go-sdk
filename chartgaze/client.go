package chartgaze

// Client is the ChartGaze Go SDK entrypoint for market intelligence.
type Client struct {
	APIKey  string
	BaseURL string
}

type Option func(*Client)

func WithAPIKey(key string) Option {
	return func(c *Client) { c.APIKey = key }
}

func NewClient(opts ...Option) (*Client, error) {
	c := &Client{BaseURL: "https://api.chartgaze.live"}
	for _, o := range opts {
		o(c)
	}
	return c, nil
}
