package types

import (
	"time"
)

// Invitation model.
type Invitation struct {
	ID        string    `json:"id" yaml:"id"`
	Email     string    `json:"email" yaml:"email"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
}

// CreateInvitation request payload for creating a project invitation.
type CreateInvitation struct {
	Email       string `json:"email"`
	RedirectURI string `json:"redirectURI"`
}

// AcceptInvitation request payload for accepting a project invitation.
type AcceptInvitation struct {
	Token string `json:"token"`
}
