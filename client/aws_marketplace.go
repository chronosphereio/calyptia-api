package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/calyptia/api/types"
)

func (c *Client) AWSCustomerRedirect(ctx context.Context, amazonMarketplaceToken string) (*url.URL, error) {
	form := url.Values{}
	form.Add("x-amzn-marketplace-token", amazonMarketplaceToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/v1/aws-customer-redirect", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	c.setRequestHeaders(req)

	defer disableRedirect(c.Client)()

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read body: %w", err)
		}

		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(b))
	}

	loc, err := url.Parse(resp.Header.Get("Location"))
	if err != nil {
		return nil, fmt.Errorf("parse location header: %w", err)
	}

	if !loc.IsAbs() {
		return nil, errors.New("invalid location header")
	}

	q := loc.Query()
	if q.Has("error") {
		return nil, errors.New(q.Get("error"))
	}

	return loc, nil
}

func (c *Client) CreateAWSContractFromToken(ctx context.Context, in types.CreateAWSContractFromToken) error {
	return c.do(ctx, http.MethodPost, "/v1/projects/"+url.PathEscape(in.ProjectID)+"/aws-contracts", in, nil)
}
