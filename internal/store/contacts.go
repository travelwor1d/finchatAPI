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

func (s *Store) ListContactRequests(ctx context.Context, contactOwnerID int, p *Pagination) ([]*model.ContactRequest, error) {
	const query = `
	SELECT * FROM contact_requests WHERE user_id = ? LIMIT ? OFFSET ?
	`
	var posts []*model.ContactRequest
	err := s.db.SelectContext(ctx, &posts, query, contactOwnerID, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *Store) GetContactRequest(ctx context.Context, id int) (*model.ContactRequest, error) {
	const query = `
	SELECT * FROM contact_requests WHERE id = ?
	`
	var r model.ContactRequest
	err := s.db.GetContext(ctx, &r, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *Store) CreateContactRequest(ctx context.Context, contactOwnerID, contactID int) (*model.ContactRequest, error) {
	const query = `
	INSERT INTO contacts(user_id, contact_id, request_status) VALUES (?, ?, 'REQUESTED')
	`
	result, err := s.db.ExecContext(ctx, query, contactOwnerID, contactID)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetContactRequest(ctx, int(id))
}

func (s *Store) ApproveContactRequest(ctx context.Context, id int) error {
	const query = `
	UPDATE contacts SET request_status = 'ACCEPTED' WHERE id = ?
	`
	_, err := s.db.ExecContext(ctx, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (s *Store) DenyContactRequest(ctx context.Context, id int) error {
	const query = `
	UPDATE contacts SET request_status = 'DENIED' WHERE id = ?
	`
	_, err := s.db.ExecContext(ctx, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
