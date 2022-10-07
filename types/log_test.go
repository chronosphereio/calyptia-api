package types

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
)

func TestFluentBitLog_JSON(t *testing.T) {
	t.Run("unmarshal", func(t *testing.T) {
		now := time.Now().UTC()
		// trucating to second precision for now.
		// TODO: test correct parsing of unix fraction time.
		now = now.Truncate(time.Second)

		secs := now.Truncate(time.Second).Unix()

		var l FluentBitLog
		err := json.Unmarshal([]byte(fmt.Sprintf(`[%d,{"foo":"bar"}]`, secs)), &l)
		assert.NoError(t, err)

		assert.Equal(t, float64(l.Timestamp), float64(now.Unix()))
		assert.Equal(t, 1, len(l.Attrs))
		assert.Equal(t, "bar", l.Attrs["foo"])
	})

	t.Run("marshal", func(t *testing.T) {
		now := time.Now().UTC()

		// trucating to second precision for now.
		// TODO: test correct parsing of unix fraction time.
		now = now.Truncate(time.Second)
		l := FluentBitLog{
			Timestamp: FluentBitTime(now.Unix()),
			Attrs: FluentBitLogAttrs{
				"foo": "bar",
			},
		}

		b, err := json.Marshal(l)
		assert.NoError(t, err)

		nsecs := now.Unix()
		assert.Equal(t, fmt.Sprintf(`[%d,{"foo":"bar"}]`, nsecs), string(b))
	})

	t.Run("back_and_forth", func(t *testing.T) {
		now := time.Now().UTC()
		// trucating to second precision for now.
		// TODO: test correct parsing of unix fraction time.
		now = now.Truncate(time.Second)

		secs := now.Truncate(time.Second).Unix()
		in := fmt.Sprintf(`[%d,{"foo":"bar"}]`, secs)

		var l FluentBitLog
		err := json.Unmarshal([]byte(in), &l)
		assert.NoError(t, err)

		out, err := json.Marshal(l)
		assert.NoError(t, err)

		assert.Equal(t, in, string(out))
	})
}
