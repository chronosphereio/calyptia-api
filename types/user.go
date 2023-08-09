package types

import (
	"strings"
	"time"
)

// User model.
type User struct {
	ID        string    `json:"id" yaml:"id"`
	Email     string    `json:"email" yaml:"email"`
	Name      string    `json:"name" yaml:"name"`
	AvatarURL *string   `json:"avatarURL" yaml:"avatarURL"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}

func (u *User) RedactEmail() string {
	atIndex := strings.Index(u.Email, "@")
	if atIndex == -1 {
		return u.Email
	}

	username := u.Email[:atIndex]
	if len(username) <= 1 {
		// Redact single-character prefix entirely
		return strings.Repeat("*", atIndex) + u.Email[atIndex:]
	}

	if len(username) == 2 {
		// Redact only the first character for two-character prefix
		prefix := string(username[0]) + "*" + u.Email[atIndex:]
		return prefix
	}

	prefix := string(username[0]) + strings.Repeat("*", len(username)-2) + string(username[len(username)-1])

	return prefix + u.Email[atIndex:]
}

type UpsertUser struct {
	Email     string
	Name      string
	AvatarURL *string
}
