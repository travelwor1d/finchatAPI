package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/finchatapp/finchat-api/internal/model"
)

func (s *Store) ListContacts(ctx context.Context, contactOwnerID int, p *Pagination) ([]*model.Contact, error) {
	const query = `
	SELECT * FROM contacts WHERE user_id = ? LIMIT ? OFFSET ?
	`
	var posts []*model.Contact
	err := s.db.SelectContext(ctx, &posts, query, contactOwnerID, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *Store) GetContact(ctx context.Context, userID, id int) (*model.Contact, error) {
	const query = `
	SELECT * FROM contacts WHERE user_id = ? AND id = ?
	`
	var c model.Contact
	err := s.db.GetContext(ctx, &c, query, userID, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) DeleteContact(ctx context.Context, userID, id int) error {
	const query = `
	DELETE FROM users_contacts WHERE user_id = ? AND id = ?
	`
	result, err := s.db.ExecContext(ctx, query, userID, id)
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
