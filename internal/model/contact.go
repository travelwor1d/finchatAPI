package model

import "time"

type Contact struct {
	ID             int       `db:"id" json:"id"`
	ContactOwnerID int       `db:"user_id" json:"contactOwnerId"`
	ContactID      int       `db:"contact_id" json:"contactId"`
	FirstName      string    `db:"first_name" json:"firstName"`
	LastName       string    `db:"last_name" json:"lastname"`
	Phone          string    `db:"phone" json:"phone"`
	Email          string    `db:"email" json:"email"`
	Type           string    `db:"user_type" json:"userType"`
	ProfileAvatar  *string   `db:"profile_avatar" json:"profiveAvatar"`
	LastSeen       time.Time `db:"last_seen" json:"lastSeen"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`
}

type ContactRequest struct {
	ID             int       `db:"id" json:"id"`
	ContactOwnerID int       `db:"user_id" json:"contactOwnerId"`
	ContactID      int       `db:"contact_id" json:"contactId"`
	Status         string    `db:"request_status" json:"status"`
	FirstName      string    `db:"first_name" json:"firstName"`
	LastName       string    `db:"last_name" json:"lastName"`
	Phone          string    `db:"phone" json:"phone"`
	Email          string    `db:"email" json:"email"`
	Type           string    `db:"user_type" json:"userType"`
	ProfileAvatar  *string   `db:"profile_avatar" json:"profileAvatar"`
	LastSeen       time.Time `db:"last_seen" json:"lastSeen"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`
}
