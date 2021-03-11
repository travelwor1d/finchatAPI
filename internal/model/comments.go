package model

import "time"

type Comment struct {
	ID          int        `db:"id"`
	PostID      int        `db:"post_id"`
	Content     string     `db:"content"`
	PostedBy    int        `db:"posted_by"`
	PublishedAt *time.Time `db:"published_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}
