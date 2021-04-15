package model

import "time"

type Thread struct {
	ID        int       `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Type      string    `db:"thread_type" json:"type"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type Message struct {
	ID        int    `db:"id" json:"id"`
	ThreadID  int    `db:"thread_id" json:"threadId"`
	SenderID  int    `db:"sender_id" json:"senderId"`
	Type      string `db:"message_type" json:"type"`
	Message   string `db:"message" json:"message"`
	Timestamp int64  `db:"timestamp" json:"timestamp"`
}
