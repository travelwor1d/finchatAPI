package models

import "time"

type Contact struct {
	ID            int       `db:"id" json:"id"`
	FirstName     string    `db:"first_name" json:"firstName"`
	LastName      string    `db:"last_name" json:"lastName"`
	Phonenumber   string    `db:"phone_number" json:"phoneNumber"`
	CountryCode   string    `db:"country_code" json:"countryCode"`
	Email         string    `db:"email" json:"email"`
	Type          string    `db:"user_type" json:"userType"`
	ProfileAvatar *string   `db:"profile_avatar" json:"profiveAvatar"`
	LastSeen      time.Time `db:"last_seen" json:"lastSeen"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
}
