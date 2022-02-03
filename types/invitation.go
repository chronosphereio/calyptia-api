package types

// CreateInvitation request payload for creating a project invitation.
type CreateInvitation struct {
	Email       string `json:"email"`
	RedirectURI string `json:"redirectURI"`
}

// AcceptInvitation request payload for accepting a project invitation.
type AcceptInvitation struct {
	Token string `json:"token"`
}
