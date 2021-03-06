package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/finchatapp/finchat-api/pkg/unique"
)

func (s *Store) CreateGoatInviteCode(ctx context.Context, userID int) (string, error) {
	const query = `
	INSERT INTO goat_invite_codes(invite_code, created_by) VALUES (?, ?)
	`
	code := unique.New(6)
	_, err := s.db.ExecContext(ctx, query, code, userID)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (s *Store) GetInviteCodeStatus(ctx context.Context, code string) (string, bool, error) {
	const query = `
	SELECT status FROM goat_invite_codes WHERE invite_code = ?
	`
	var status string
	err := s.db.GetContext(ctx, &status, query, code)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return status, true, nil
}

func (s *Store) UseInviteCode(ctx context.Context, code string, userID int) error {
	const query = `
	UPDATE goat_invite_codes SET
		status = 'USED',
		used_by = ?
	WHERE invite_code = ? 
	`
	result, err := s.db.ExecContext(ctx, query, userID, code)
	if err != nil {
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
