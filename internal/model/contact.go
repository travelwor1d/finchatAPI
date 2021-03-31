package model

import "time"

type Contact struct {
	ID             int       `db:"id"`
	ContactOwnerID int       `db:"user_id"`
	ContactID      int       `db:"contact_id"`
	FirstName      string    `db:"first_name"`
	LastName       string    `db:"last_name"`
	Phone          string    `db:"phone"`
	Email          string    `db:"email"`
	Type           string    `db:"user_type"`
	ProfileAvatar  *string   `db:"profile_avatar"`
	LastSeen       time.Time `db:"last_seen"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type ContactRequest struct {
	ID             int       `db:"id"`
	ContactOwnerID int       `db:"user_id"`
	ContactID      int       `db:"contact_id"`
	Status         string    `db:"request_status"`
	FirstName      string    `db:"first_name"`
	LastName       string    `db:"last_name"`
	Phone          string    `db:"phone"`
	Email          string    `db:"email"`
	Type           string    `db:"user_type"`
	ProfileAvatar  *string   `db:"profile_avatar"`
	LastSeen       time.Time `db:"last_seen"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
