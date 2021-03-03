package model

import "time"

type Creds struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Hash      string    `db:"hash"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
