package types

import "time"

type Invitation struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"-"`
	SigningKey  []byte    `json:"-"`
	Email       string    `json:"email"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"createdAt"`
}

// CreateInvitation request payload for creating a project invitation.
type CreateInvitation struct {
	Email       string   `json:"email"`
	RedirectURI string   `json:"redirectURI"`
	Permissions []string `json:"permissions"`

	id         string
	signingKey []byte
}

func (in *CreateInvitation) SetID(id string) {
	in.id = id
}

func (in *CreateInvitation) SetSigningKey(signingKey []byte) {
	in.signingKey = signingKey
}

func (in CreateInvitation) ID() string {
	return in.id
}

func (in CreateInvitation) SigningKey() []byte {
	return in.signingKey
}

// AcceptInvitation request payload for accepting a project invitation.
type AcceptInvitation struct {
	Token string `json:"token"`
}
