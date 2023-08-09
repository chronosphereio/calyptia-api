package errs

const (
	InvalidToken     = UnauthenticatedError("invalid token")
	InvalidTokenID   = InvalidArgumentError("invalid token ID")
	InvalidTokenName = InvalidArgumentError("invalid token name")
	TokenNameTaken   = ConflictError("token name taken")
	TokenGone        = GoneError("token gone")
	TokenNotFound    = NotFoundError("token not found")
)
