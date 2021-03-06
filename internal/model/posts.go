package model

import "time"

type Post struct {
	ID          int         `db:"id" json:"id"`
	Title       string      `db:"title" json:"title"`
	ImageURLs   StringSlice `db:"image_urls" json:"imageUrls"`
	Tickers     StringSlice `db:"tickers" json:"tickers"`
	Content     string      `db:"content" json:"content"`
	Upvotes     int         `db:"upvotes" json:"upvotes"`
	Downvotes   int         `db:"downvotes" json:"downvotes"`
	PostedBy    int         `db:"posted_by" json:"postedBy"`
	PublishedAt *time.Time  `db:"published_at" json:"publishedAt"`
	CreatedAt   time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time   `db:"updated_at" json:"updatedAt"`
}
