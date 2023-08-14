package errs

const (
	EmailNotVerified = UnauthenticatedError("email not verified")
	InvalidEmail     = InvalidArgumentError("invalid email")
	InvalidUserName  = InvalidArgumentError("invalid user name")
	InvalidAvatarURL = InvalidArgumentError("invalid avatar URL")
	EmailTaken       = ConflictError("email taken")
	UserNotFound     = NotFoundError("user not found")
)
