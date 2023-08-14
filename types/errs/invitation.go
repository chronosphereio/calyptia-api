package errs

const (
	InvalidInvitationToken = UnauthenticatedError("invalid invitation token")
	InvitationExpired      = UnauthenticatedError("invitation expired")
	InvalidRedirectURI     = InvalidArgumentError("invalid redirect URI")
	UntrustedRedirectURI   = InvalidArgumentError("untrusted redirect URI")
	EmailMismatch          = PermissionDeniedError("email mismatch")
	CannotInviteYourSelf   = PermissionDeniedError("cannot invite yourself")
	InvitationGone         = GoneError("invitation gone")
	InvitationNotFound     = NotFoundError("invitation not found")
)
