package types

const (
	ErrInvalidInvitationToken = UnauthenticatedError("invalid invitation token")
	ErrInvitationExpired      = UnauthenticatedError("invitation expired")
	ErrInvalidRedirectURI     = InvalidArgumentError("invalid redirect URI")
	ErrUntrustedRedirectURI   = InvalidArgumentError("untrusted redirect URI")
	ErrEmailMismatch          = PermissionDeniedError("email mismatch")
	ErrCannotInviteYourSelf   = PermissionDeniedError("cannot invite yourself")
	ErrInvitationGone         = GoneError("invitation gone")
	ErrInvitationNotFound     = NotFoundError("invitation not found")
)

// CreateInvitation request payload for creating a project invitation.
type CreateInvitation struct {
	Email       string `json:"email"`
	RedirectURI string `json:"redirectURI"`
}

// AcceptInvitation request payload for accepting a project invitation.
type AcceptInvitation struct {
	Token string `json:"token"`
}
