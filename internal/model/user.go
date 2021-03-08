package model

import "time"

type User struct {
	ID            int       `db:"id"`
	FirstName     string    `db:"first_name"`
	LastName      string    `db:"last_name"`
	Phone         string    `db:"phone"`
	Email         string    `db:"email"`
	Verified      bool      `db:"verified"`
	Type          string    `db:"user_type"`
	ProfileAvatar *string   `db:"profile_avatar"`
	LastSeen      time.Time `db:"last_seen"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
