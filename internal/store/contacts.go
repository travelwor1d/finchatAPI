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

func (s *Store) ListContactRequests(ctx context.Context, contactOwnerID int, p *Pagination) ([]*model.Contact, error) {
	const query = `
	SELECT * FROM contact_requests user_id = ? LIMIT ? OFFSET ?
	`
	var posts []*model.Contact
	err := s.db.SelectContext(ctx, &posts, query, contactOwnerID, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *Store) CreateContactRequest(ctx context.Context, contactOwnerID, contactID int) error {
	const query = `
	INSERT INTO contacts(user_id, contact_id, request_status) VALUES (?, ?, 'REQUESTED')
	`
	_, err := s.db.ExecContext(ctx, query, contactOwnerID, contactID)
	return err
}

func (s *Store) ApproveContactRequest(ctx context.Context, contactOwnerID, contactID int) error {
	const query = `
	UPDATE contacts SET request_status = 'ACCEPTED' WHERE user_id = ? AND contact_id = ?
	`
	_, err := s.db.ExecContext(ctx, query, contactOwnerID, contactID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (s *Store) DenyContactRequest(ctx context.Context, contactOwnerID, contactID int) error {
	const query = `
	UPDATE contacts SET request_status = 'DENIED' WHERE user_id = ? AND contact_id = ?
	`
	_, err := s.db.ExecContext(ctx, query, contactOwnerID, contactID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
