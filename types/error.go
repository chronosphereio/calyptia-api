package types

import "fmt"

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
