package store

import "context"

func (s *Store) CreateUpvote(ctx context.Context, postID, userID int) error {
	const query = `
	INSERT INTO post_votes(post_id, user_id, value) VALUES (?, ?, 1)
		ON DUPLICATE KEY UPDATE value = VALUES(value)
	`
	_, err := s.db.ExecContext(ctx, query, postID, userID)
	return err
}

func (s *Store) DeleteUpvote(ctx context.Context, postID, userID int) error {
	const query = `
	DELETE FROM post_votes WHERE post_id = ? AND user_id = ? AND value = 1
	`
	_, err := s.db.ExecContext(ctx, query, postID, userID)
	return err
}

func (s *Store) CreateDownvote(ctx context.Context, postID, userID int) error {
	const query = `
	INSERT INTO post_votes(post_id, user_id, value) VALUES (?, ?, -1)
		ON DUPLICATE KEY UPDATE value = VALUES(value)
	`
	_, err := s.db.ExecContext(ctx, query, postID, userID)
	return err
}

func (s *Store) DeleteDownvote(ctx context.Context, postID, userID int) error {
	const query = `
	DELETE FROM post_votes WHERE post_id = ? AND user_id = ? AND value = -1
	`
	_, err := s.db.ExecContext(ctx, query, postID, userID)
	return err
}
