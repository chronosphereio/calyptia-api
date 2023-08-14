package errs

const (
	InvalidMemberID     = InvalidArgumentError("invalid member ID")
	MemberAlreadyExists = ConflictError("member already exists")
	MemberNotFound      = NotFoundError("member not found")
)
