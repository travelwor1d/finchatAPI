package model

import "time"

type Comment struct {
	ID          int        `db:"id" json:"id"`
	PostID      int        `db:"post_id" json:"postId"`
	Content     string     `db:"content" json:"content"`
	PostedBy    int        `db:"posted_by" json:"postedBy"`
	PublishedAt *time.Time `db:"published_at" json:"publishedAt"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
}
