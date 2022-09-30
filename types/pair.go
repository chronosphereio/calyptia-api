package types

import (
	"strings"
)

// Pair struct. It can be used to define any sorted key value pairs.
// For example fluent-bit config section properties.
type Pair struct {
	Key   string `json:"key" yaml:"key"`
	Value any    `json:"value" yaml:"value"`
}

// Pairs list.
type Pairs []Pair

// Get a value matching its key case-insensitively.
func (pp Pairs) Get(key string) (any, bool) {
	for _, p := range pp {
		if strings.EqualFold(p.Key, key) {
			return p.Value, true
		}
	}
	return nil, false
}

// Has checks whether the key exists (case-insensitively).
func (pp Pairs) Has(key string) bool {
	for _, p := range pp {
		if strings.EqualFold(p.Key, key) {
			return true
		}
	}
	return false
}

// Set a value, overriding any axisting key (case-insensitively).
func (pp *Pairs) Set(key string, value any) {
	for i, p := range *pp {
		if strings.EqualFold(p.Key, key) {
			(*pp)[i].Value = value
			return
		}
	}
	*pp = append(*pp, Pair{Key: key, Value: value})
}

// Remove a value using the given key (case-insensitively).
func (pp *Pairs) Remove(key string) {
	deleteIndex := -1
	for i, p := range *pp {
		if strings.EqualFold(p.Key, key) {
			deleteIndex = i
			break
		}
	}

	*pp = append((*pp)[:deleteIndex], (*pp)[deleteIndex+1:]...)
}

func (pp Pairs) AsMap() map[string]any {
	if pp == nil {
		return nil
	}

	m := make(map[string]any, len(pp))
	for _, p := range pp {
		m[p.Key] = p.Value
	}
	return m
}

// FixFloats handles JSON side-effect of converting numbers into float64
// when using `any` type.
// This tries to convert them back into int when possible.
// Call this method after json unmarshalling.
func (pp Pairs) FixFloats() {
	for i, p := range pp {
		if f, ok := p.Value.(float64); ok {
			if f == float64(int(f)) {
				pp[i].Value = int(f)
			} else if f == float64(int64(f)) {
				pp[i].Value = int64(f)
			}
		}
	}
}
