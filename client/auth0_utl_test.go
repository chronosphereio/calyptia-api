package client_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type Auth0Client struct {
	BaseURL      string
	Client       *http.Client
	ClientID     string
	ClientSecret string
	Audience     string
	UserAgent    string
}

type Auth0UserClaims struct {
	Email         string
	EmailVerified bool
	Name          string
}

type Auth0AccessTokenReq struct {
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	Audience      string `json:"audience"`
	GranType      string `json:"grant_type"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
	Name          string `json:"name,omitempty"`
}

type Auth0Error struct {
	Msg         string `json:"error,omitempty"`
	Description string `json:"error_description,omitempty"`
	URI         string `json:"error_uri,omitempty"`
}

func (e Auth0Error) Error() string {
	if e.Description != "" {
		return fmt.Sprintf("%s: %s", e.Msg, e.Description)
	}
	return e.Msg
}

func (c *Auth0Client) AccessToken(ctx context.Context, claims Auth0UserClaims) (*oauth2.Token, error) {
	b, err := json.Marshal(Auth0AccessTokenReq{
		ClientID:      c.ClientID,
		ClientSecret:  c.ClientSecret,
		Audience:      c.Audience,
		GranType:      "client_credentials",
		Email:         claims.Email,
		EmailVerified: true,
		Name:          claims.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("could not json encode access token request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/oauth/token", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("could not make access token request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not do access token request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		e := &Auth0Error{}
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&e)
		if err != nil {
			return nil, fmt.Errorf("could not json decode access token error response body: status_code=%d: %w", resp.StatusCode, err)
		}

		return nil, e
	}

	out := &oauth2.Token{}
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return nil, fmt.Errorf("could not json decode access token response body: %w", err)
	}

	return out, nil
}
