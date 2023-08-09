// Package errs provides common error types that can be used across the application.
package errs

const (
	Unauthenticated  = UnauthenticatedError("unauthenticated")
	InvalidArgument  = InvalidArgumentError("invalid argument")
	NotFound         = NotFoundError("not found")
	Conflict         = ConflictError("conflict")
	PermissionDenied = PermissionDeniedError("permission denied")
	Gone             = GoneError("gone")
)

type UnauthenticatedError string

func (e UnauthenticatedError) Error() string        { return string(e) }
func (e UnauthenticatedError) Is(target error) bool { return target == Unauthenticated }

type InvalidArgumentError string

func (e InvalidArgumentError) Error() string        { return string(e) }
func (e InvalidArgumentError) Is(target error) bool { return target == InvalidArgument }

type NotFoundError string

func (e NotFoundError) Error() string        { return string(e) }
func (e NotFoundError) Is(target error) bool { return target == NotFound }

type ConflictError string

func (e ConflictError) Error() string        { return string(e) }
func (e ConflictError) Is(target error) bool { return target == Conflict }

type PermissionDeniedError string

func (e PermissionDeniedError) Error() string        { return string(e) }
func (e PermissionDeniedError) Is(target error) bool { return target == PermissionDenied }

type GoneError string

func (e GoneError) Error() string        { return string(e) }
func (e GoneError) Is(target error) bool { return target == Gone }

type DetailedError struct {
	Msg    string  `json:"error"`
	Detail *string `json:"detail,omitempty"`
	Parent error   `json:"-"`
}

func NewDetailedError(parent error, detail string) *DetailedError {
	return &DetailedError{
		Msg:    parent.Error(),
		Detail: &detail,
		Parent: parent,
	}
}

func (e *DetailedError) Error() string { return e.Msg }
func (e *DetailedError) Unwrap() error { return e.Parent }

type IndexedError struct {
	Msg    string `json:"error"`
	Index  uint   `json:"index"`
	Parent error  `json:"-"`
}

func (e *IndexedError) Error() string { return e.Msg }
func (e *IndexedError) Unwrap() error { return e.Parent }

func NewIndexedError(e error, index uint) *IndexedError {
	return &IndexedError{
		Msg:    e.Error(),
		Index:  index,
		Parent: e,
	}
}

type QuotaExceededError struct {
	Msg       string `json:"error"`
	Limit     uint   `json:"quotaLimit"`
	Remaining uint   `json:"quotaRemaining"`
	Parent    error  `json:"-"`
}

func (e *QuotaExceededError) Error() string { return e.Msg }
func (e *QuotaExceededError) Unwrap() error { return e.Parent }

func NewQuotaExceededError(base error, limit, remaining uint) *QuotaExceededError {
	return &QuotaExceededError{
		Msg:       base.Error(),
		Limit:     limit,
		Remaining: remaining,
		Parent:    base,
	}
}
