package model

import "time"

type User struct {
	ID            int       `db:"id" json:"id"`
	StripeID      *string   `db:"stripe_id" json:"-"`
	FirstName     string    `db:"first_name" json:"firstName"`
	LastName      string    `db:"last_name" json:"lastName"`
	Phone         string    `db:"phone" json:"phone"`
	Email         string    `db:"email" json:"email"`
	Verified      bool      `db:"verified" json:"verified"`
	Type          string    `db:"user_type" json:"userType"`
	ProfileAvatar *string   `db:"profile_avatar" json:"profileAvatar"`
	LastSeen      time.Time `db:"last_seen" json:"lastSeen"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
}
