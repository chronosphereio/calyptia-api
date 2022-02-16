package types

import "fmt"

const (
	ErrUnauthenticated  = UnauthenticatedError("unauthenticated")
	ErrInvalidArgument  = InvalidArgumentError("invalid argument")
	ErrNotFound         = NotFoundError("not found")
	ErrConflict         = ConflictError("conflict")
	ErrPermissionDenied = PermissionDeniedError("permission denied")
	ErrGone             = GoneError("gone")
)

type UnauthenticatedError string

func (e UnauthenticatedError) Error() string        { return string(e) }
func (e UnauthenticatedError) Is(target error) bool { return target == ErrUnauthenticated }

type InvalidArgumentError string

func (e InvalidArgumentError) Error() string        { return string(e) }
func (e InvalidArgumentError) Is(target error) bool { return target == ErrInvalidArgument }

type NotFoundError string

func (e NotFoundError) Error() string        { return string(e) }
func (e NotFoundError) Is(target error) bool { return target == ErrNotFound }

type ConflictError string

func (e ConflictError) Error() string        { return string(e) }
func (e ConflictError) Is(target error) bool { return target == ErrConflict }

type PermissionDeniedError string

func (e PermissionDeniedError) Error() string        { return string(e) }
func (e PermissionDeniedError) Is(target error) bool { return target == ErrPermissionDenied }

type GoneError string

func (e GoneError) Error() string        { return string(e) }
func (e GoneError) Is(target error) bool { return target == ErrGone }

// Error returned by the API.
type Error struct {
	Msg    string  `json:"error"`
	Detail *string `json:"detail"`
}

func (e *Error) Error() string {
	if e.Detail != nil {
		return fmt.Sprintf("%s: %s", e.Msg, *e.Detail)
	}
	return e.Msg
}
