// Package client provides a client over the REST HTTP API of Calyptia Cloud.
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

// Client is the client over the REST HTTP API of Calyptia Cloud.
type Client struct {
	BaseURL string
	// Tip: Use oauth2.NewClient(context.Context, *oauth2.TokenSource)
	Client          *http.Client
	userAgent       string
	projectToken    string
	agentToken      string
	aggregatorToken string
}

const (
	// DefaultBaseURL of the API.
	DefaultBaseURL        = "https://cloud-api.calyptia.com"
	defaultUintFormatBase = 10
)

// New default client.
func New() *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		Client:  http.DefaultClient,
	}
}

// SetUserAgent sets the "User-Agent" header of the client.
func (c *Client) SetUserAgent(s string) {
	c.userAgent = s
}

// SetProjectToken sets the "X-Project-Token" header of the client.
func (c *Client) SetProjectToken(s string) {
	c.projectToken = s
}

// SetAgentToken sets the "X-Agent-Token" header of the client.
func (c *Client) SetAgentToken(s string) {
	c.agentToken = s
}

// SetAggregatorToken sets the "X-Aggregator-Token" header of the client.
func (c *Client) SetAggregatorToken(s string) {
	c.aggregatorToken = s
}

func (c *Client) setRequestHeaders(req *http.Request) {
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

	c.setRequestHeaders(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= http.StatusBadRequest {
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
