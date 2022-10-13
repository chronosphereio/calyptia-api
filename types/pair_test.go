package types

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestPairs_Set(t *testing.T) {
	t.Run("override", func(t *testing.T) {
		pp := Pairs{
			{Key: "name", Value: "dummy"},
		}
		assert.NotPanics(t, func() {
			pp.Set("name", "cpu")
		})
		assert.Equal(t, Pairs{{Key: "name", Value: "cpu"}}, pp)
	})

	t.Run("add", func(t *testing.T) {
		pp := Pairs{}
		assert.NotPanics(t, func() {
			pp.Set("name", "dummy")
		})
		assert.Equal(t, Pairs{{Key: "name", Value: "dummy"}}, pp)
	})
}

func TestPairs_Remove(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		pp := Pairs{
			{Key: "name", Value: "dummy"},
			{Key: "match", Value: "*"},
		}
		assert.NotPanics(t, func() {
			pp.Remove("match")
		})
		assert.Equal(t, Pairs{{Key: "name", Value: "dummy"}}, pp)
	})

	t.Run("not_existing", func(t *testing.T) {
		pp := Pairs{}
		assert.NotPanics(t, func() {
			pp.Remove("match")
		})
		assert.Equal(t, Pairs{}, pp)
	})
}
