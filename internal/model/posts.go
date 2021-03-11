package model

import "time"

type Post struct {
	ID          int        `db:"id"`
	Title       string     `db:"title"`
	Content     string     `db:"content"`
	PostedBy    int        `db:"posted_by"`
	PublishedAt *time.Time `db:"published_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}
