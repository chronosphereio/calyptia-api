package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/calyptia/api/types"
)

type Client struct {
	BaseURL string
	// Tip: Use oauth2.NewClient(context.Context, *oauth2.TokenSource)
	Client          *http.Client
	userAgent       string
	projectToken    string
	agentToken      string
	aggregatorToken string
}

func New(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		Client:  http.DefaultClient,
	}
}

func (c *Client) SetUserAgent(s string) {
	c.userAgent = s
}

func (c *Client) SetProjectToken(s string) {
	c.projectToken = s
}

func (c *Client) SetAgentToken(s string) {
	c.agentToken = s
}

func (c *Client) SetAggregatorToken(s string) {
	c.aggregatorToken = s
}

func (c *Client) do(ctx context.Context, method, path string, v, dest interface{}) error {
	var body io.Reader

	if b, ok := v.(io.Reader); ok {
		body = b
	} else if b, ok := v.([]byte); ok {
		body = bytes.NewReader(b)
	} else if v != nil {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}

		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, body)
	if err != nil {
		return err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	if c.projectToken != "" {
		req.Header.Set("X-Project-Token", c.projectToken)
	}

	if c.agentToken != "" {
		req.Header.Set("X-Agent-Token", c.agentToken)
	}

	if c.aggregatorToken != "" {
		req.Header.Set("X-Aggregator-Token", c.aggregatorToken)
	}

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		e := &types.Error{}
		err = json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return fmt.Errorf("could not json decode error response: %w", err)
		}

		return e
	}

	if dest != nil {
		err = json.NewDecoder(resp.Body).Decode(dest)
		if err != nil {
			return fmt.Errorf("could not json decode response: %w", err)
		}
	}

	return nil
}
