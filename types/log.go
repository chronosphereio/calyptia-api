package types

import (
	"encoding/json"
	"errors"
)

// FluentBitLog is the Go representation of a fluent-bit log record.
// Normally, the format is an array of mixed types.
// [timestamp, attrs]
// Where timestampt is a float with seconds and nanoseconds as fraction.
// And attrs being an object with string key and values.
type FluentBitLog struct {
	Timestamp FluentBitTime
	Attrs     FluentBitLogAttrs
}

// FluentBitTime wrapper.
type FluentBitTime float64

// FluentBitLogAttrs wrapper.
type FluentBitLogAttrs map[string]string

func (l FluentBitLog) AsSlice() []any {
	return []any{float64(l.Timestamp), map[string]string(l.Attrs)}
}

// UnmarshalJSON parses fluent-big json representation of
// [timestamp, attrs].
func (l *FluentBitLog) UnmarshalJSON(b []byte) error {
	var ss [2]json.RawMessage

	err := json.Unmarshal(b, &ss)
	if err != nil {
		return err
	}

	if len(ss) != 2 {
		return errors.New("unexpected log record parts length")
	}

	err = json.Unmarshal(ss[0], &l.Timestamp)
	if err != nil {
		return err
	}

	err = json.Unmarshal(ss[1], &l.Attrs)
	if err != nil {
		return err
	}

	return nil
}

// MarshalJSON into fluent-bit json representation of
// [timestamp, attrs].
func (l FluentBitLog) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.AsSlice())
}
