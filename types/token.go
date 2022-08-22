package types

import "time"

// Token model.
type Token struct {
	ID        string    `json:"id" yaml:"id"`
	Token     string    `json:"token" yaml:"token"`
	Name      string    `json:"name" yaml:"name"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
}

// Tokens paginated list.
type Tokens struct {
	Items     []Token
	EndCursor *string
}

// CreateToken request payload for creating a new token.
type CreateToken struct {
	Name string `json:"name"`
}

// TokensParams request payload for querying tokens.
type TokensParams struct {
	Last   *uint
	Before *string
	Name   *string
}

// UpdateToken request payload for updating a token.
type UpdateToken struct {
	Name *string `json:"name"`
}
