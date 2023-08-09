package types

import (
	"fmt"

	"github.com/go-logfmt/logfmt"
)

// Error returned by the API.
type Error struct {
	Msg            string  `json:"error"`
	Detail         *string `json:"detail,omitempty"`
	Index          *uint   `json:"index,omitempty"`
	QuotaLimit     *uint   `json:"quotaLimit,omitempty"`
	QuotaRemaining *uint   `json:"quotaRemaining,omitempty"`
}

func (e *Error) Error() string {
	keyvals := []any{}
	if e.Detail != nil {
		keyvals = append(keyvals, "detail", *e.Detail)
	}
	if e.Index != nil {
		keyvals = append(keyvals, "index", *e.Index)
	}
	if e.QuotaLimit != nil {
		keyvals = append(keyvals, "quota_limit", *e.QuotaLimit)
	}
	if e.QuotaRemaining != nil {
		keyvals = append(keyvals, "quota_remaining", *e.QuotaRemaining)
	}

	if len(keyvals) != 0 {
		extra, err := logfmt.MarshalKeyvals(keyvals...)
		if err == nil && len(extra) != 0 {
			return fmt.Sprintf("%s: %s", e.Msg, extra)
		}
	}

	return e.Msg
}
