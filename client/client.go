// Package client provides a client over the REST HTTP API of Calyptia Cloud.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/peterhellberg/link"
	"gopkg.in/yaml.v2"

	"github.com/calyptia/api/types"
)

const (
	// DefaultBaseURL of the API.
	DefaultBaseURL    = "https://cloud-api.calyptia.com"
	uintBase          = 10
	CalyptiaEnvAPIURL = "CALYPTIA_API_URL"
	//nolint: gosec // no credential leaks just a variable.
	CalyptiaEnvAPIToken = "CALYPTIA_API_TOKEN"
)

// Client is the client over the REST HTTP API of Calyptia Cloud.
type Client struct {
	BaseURL string
	// Tip: Use oauth2.NewClient(context.Context, *oauth2.TokenSource)
	Client            *http.Client
	userAgent         string
	projectToken      string
	agentToken        string
	coreInstanceToken string
}

// New default client.
func New() *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		Client:  http.DefaultClient,
	}
}

// NewFromEnv creates a new client using default environment variables.
func NewFromEnv() *Client {
	c := New()

	// Update the BaseURL if environment variable is set
	if apiURL := os.Getenv(CalyptiaEnvAPIURL); apiURL != "" {
		c.BaseURL = apiURL
	}

	// Set the ProjectToken if environment variable is set
	if token := os.Getenv(CalyptiaEnvAPIToken); token != "" {
		c.SetProjectToken(token)
	}
	return c
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

// SetCoreInstanceToken sets the "X-Aggregator-Token" header of the client.
func (c *Client) SetCoreInstanceToken(s string) {
	c.coreInstanceToken = s
}

func (c *Client) setRequestHeaders(req *http.Request) {
	if c.projectToken != "" {
		req.Header.Set("X-Project-Token", c.projectToken)
	}

	if c.agentToken != "" {
		req.Header.Set("X-Agent-Token", c.agentToken)
	}

	if c.coreInstanceToken != "" {
		req.Header.Set("X-Aggregator-Token", c.coreInstanceToken)
	}

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
}

func decodeResponse(resp *http.Response, dest interface{}) error {
	var err error

	switch resp.Header.Get("Content-Type") {
	default:
		fallthrough
	case "application/json":
		err = json.NewDecoder(resp.Body).Decode(dest)
		if err != nil {
			return fmt.Errorf("could not json decode response: %w", err)
		}
	case "application/yaml":
		err = yaml.NewDecoder(resp.Body).Decode(dest)
		if err != nil {
			return fmt.Errorf("could not yaml decode response: %w", err)
		}
	case "text/plain":
		strp, ok := dest.(*string)
		if !ok {
			return fmt.Errorf("must decode plain text response into string")
		}
		r, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body text: %w", err)
		}
		*strp = string(r)
	}
	return nil
}

func (c *Client) do(ctx context.Context, method, path string, v, dest interface{}, oo ...opt) error {
	var options opts
	for _, o := range oo {
		o(&options)
	}

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

	endpoint := strings.TrimSuffix(c.BaseURL, "/") + path

	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
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
			return fmt.Errorf("could not json decode error response: status_code=%d: %w", resp.StatusCode, err)
		}

		return e
	}

	if s := resp.Header.Get("Link"); s != "" && options.cursor != nil {
		before, err := nextLinkBefore(s)
		if err != nil {
			return err
		}

		if before != "" {
			*options.cursor = &before
		}
	}

	if dest != nil {
		return decodeResponse(resp, dest)
	}

	return nil
}

type opts struct {
	cursor **string
}

type opt func(*opts)

func withCursor(s **string) opt {
	return func(o *opts) {
		o.cursor = s
	}
}

// nextLinkBefore parses the given `Link` header
// and extracts the `?before=` query string parameter
// from the URI.
// It can return an empty string.
func nextLinkBefore(s string) (string, error) {
	for _, l := range link.Parse(s) {
		if l.Rel != "next" {
			continue
		}

		u, err := url.Parse(l.URI)
		if err != nil {
			return "", fmt.Errorf("could not parse link header uri: %w", err)
		}

		if !u.Query().Has("before") {
			continue
		}

		return u.Query().Get("before"), nil
	}

	return "", nil
}

func disableRedirect(c *http.Client) func() {
	old := c.CheckRedirect
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return func() {
		c.CheckRedirect = old
	}
}
