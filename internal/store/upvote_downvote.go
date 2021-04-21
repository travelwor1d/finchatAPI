package store

import "context"

func (s *Store) CreateUpvote(ctx context.Context, postID, userID int) error {
	const query = `
	INSERT INTO post_votes(post_id, user_id, value) VALUES (?, ?, 1)
		ON DUPLICATE KEY UPDATE value = VALUES(value)
	`
	result, err := s.db.ExecContext(ctx, query, postID, userID)
	if err != nil {
		// TODO: move err code to a constant.
		if checkErrCode(err, 1452) {
			return ErrNotFound
		}
		return err
	}
	r, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func (s *Store) DeleteUpvote(ctx context.Context, postID, userID int) error {
	const query = `
	DELETE FROM post_votes WHERE post_id = ? AND user_id = ? AND value = 1
	`
	result, err := s.db.ExecContext(ctx, query, postID, userID)
	if err != nil {
		return err
	}
	r, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) CreateDownvote(ctx context.Context, postID, userID int) error {
	const query = `
	INSERT INTO post_votes(post_id, user_id, value) VALUES (?, ?, -1)
		ON DUPLICATE KEY UPDATE value = VALUES(value)
	`
	result, err := s.db.ExecContext(ctx, query, postID, userID)
	if err != nil {
		// TODO: move err code to a constant.
		if checkErrCode(err, 1452) {
			return ErrNotFound
		}
		return err
	}
	r, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func (s *Store) DeleteDownvote(ctx context.Context, postID, userID int) error {
	const query = `
	DELETE FROM post_votes WHERE post_id = ? AND user_id = ? AND value = -1
	`
	result, err := s.db.ExecContext(ctx, query, postID, userID)
	if err != nil {
		return err
	}
	r, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		return ErrNotFound
	}
	return nil
}
