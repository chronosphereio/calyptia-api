package types

import "time"

// User model.
type User struct {
	ID        string    `json:"id" yaml:"id"`
	Email     string    `json:"email" yaml:"email"`
	Name      string    `json:"name" yaml:"name"`
	AvatarURL *string   `json:"avatarURL" yaml:"avatarURL"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}
