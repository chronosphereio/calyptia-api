package types

import "time"

// Token model.
type Token struct {
	ID          string    `json:"id" yaml:"id"`
	Token       string    `json:"token" yaml:"token"`
	Name        string    `json:"name" yaml:"name"`
	Permissions []string  `json:"permissions" yaml:"permissions"`
	CreatedAt   time.Time `json:"createdAt" yaml:"createdAt"`
}

// Tokens paginated list.
type Tokens struct {
	Items     []Token
	EndCursor *string
}

// CreateToken request payload for creating a new token.
type CreateToken struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`

	id         string
	signingKey []byte
	token      string
}

func (in *CreateToken) SetID(id string) {
	in.id = id
}

func (in *CreateToken) SetSigningKey(signingKey []byte) {
	in.signingKey = signingKey
}

func (in *CreateToken) SetToken(token string) {
	in.token = token
}

func (in CreateToken) ID() string {
	return in.id
}

func (in CreateToken) SigningKey() []byte {
	return in.signingKey
}

func (in CreateToken) Token() string {
	return in.token
}

// TokensParams request payload for querying tokens.
type TokensParams struct {
	Last   *uint
	Before *string
	Name   *string
}

// UpdateToken request payload for updating a token.
type UpdateToken struct {
	TokenID     string    `json:"-"`
	Name        *string   `json:"name"`
	Permissions *[]string `json:"permissions"`
}
