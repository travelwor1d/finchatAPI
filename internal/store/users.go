package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/finchatapp/finchat-api/internal/model"
	"github.com/go-sql-driver/mysql"
)

var (
	ErrUserDeleted    = errors.New("cannot perform action: user was deleted")
	ErrUserNotDeleted = errors.New("cannot undelete user: user has not been deleted")
)

func (s *Store) GetUser(ctx context.Context, id int) (*model.User, error) {
	const query = `
	SELECT * FROM users WHERE id = ? AND deleted_at IS NULL
	`
	var user model.User
	err := s.db.GetContext(ctx, &user, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	const query = `
	SELECT * FROM users WHERE email = ? AND deleted_at IS NULL
	`
	var user model.User
	err := s.db.GetContext(ctx, &user, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) SearchUsers(ctx context.Context, searchInput, userTypes string, p *Pagination) ([]*model.User, error) {
	query := fmt.Sprintf(`
	SELECT * FROM users
		WHERE (
			lower(first_name) LIKE '%s' OR
			lower(last_name) LIKE '%s'
		) AND user_type IN (%s) AND deleted_at IS NULL
	LIMIT ? OFFSET ?
	`, "%"+searchInput+"%", "%"+searchInput+"%", userTypes)
	var users []*model.User
	err := s.db.SelectContext(ctx, &users, query, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Store) CreateUser(ctx context.Context, user *model.User, password string, inviteCode ...string) (*model.User, error) {
	const query = `
	INSERT INTO users(first_name, last_name, phone, email, user_type, profile_avatar)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	tx, err := s.Begin()
	if err != nil {
		return nil, err
	}

	result, err := tx.ExecContext(ctx, query, user.FirstName, user.LastName, user.Phone, user.Email, user.Type, user.ProfileAvatar)
	if err != nil {
		me, ok := err.(*mysql.MySQLError)
		if !ok {
			tx.Rollback()
			return nil, err
		}
		if me.Number == 1062 {
			tx.Rollback()
			return nil, ErrAlreadyExists
		}
		tx.Rollback()
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	user, err = s.WithTx(tx).GetUser(ctx, int(id))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if user.Type == "GOAT" {
		if len(inviteCode) != 1 {
			return nil, errors.New("invalid usage of store.CreateUser")
		}
		err := s.WithTx(tx).UseInviteCode(ctx, inviteCode[0], int(id))
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = s.WithTx(tx).SetPassword(ctx, user.ID, password)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) UpdateUser(ctx context.Context, userID int, firstName, lastName, profileAvatar *string) (*model.User, error) {
	isDeleted, err := s.isUserDeleted(ctx, userID)
	if err != nil {
		return nil, err
	}
	if isDeleted {
		return nil, ErrUserDeleted
	}
	const query = `
	UPDATE users SET
		first_name = coalesce(?, first_name),
		last_name = coalesce(?, last_name),
		profile_avatar = coalesce(?, profile_avatar)
	WHERE id = ?
	`
	_, err = s.db.ExecContext(ctx, query, firstName, lastName, profileAvatar, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return s.GetUser(ctx, userID)
}

func (s *Store) SoftDeleteUser(ctx context.Context, userID int) error {
	isDeleted, err := s.isUserDeleted(ctx, userID)
	if err != nil {
		return err
	}
	if isDeleted {
		// Returns error that user was already"soft deleted".
		return ErrUserDeleted
	}
	const query = `
	UPDATE users SET
		deleted_at = now()
	WHERE id = ?
	`
	_, err = s.db.ExecContext(ctx, query, userID)
	return err
}

func (s *Store) UndeleteUser(ctx context.Context, userID int) error {
	isDeleted, err := s.isUserDeleted(ctx, userID)
	if err != nil {
		return err
	}
	if !isDeleted {
		// Returns error that user has not been "soft deleted", yet.
		return ErrUserNotDeleted
	}
	const query = `
	UPDATE users SET
		deleted_at = NULL
	WHERE id = ?
	`
	_, err = s.db.ExecContext(ctx, query, userID)
	return err
}

func (s *Store) SetStripeID(ctx context.Context, userID int, stripeID string) error {
	isDeleted, err := s.isUserDeleted(ctx, userID)
	if err != nil {
		return err
	}
	if isDeleted {
		return ErrUserDeleted
	}
	const query = `
	UPDATE users SET
		stripe_id = ?
	WHERE id = ?
	`
	_, err = s.db.ExecContext(ctx, query, stripeID, userID)
	return err
}

func (s *Store) isUserDeleted(ctx context.Context, userID int) (bool, error) {
	const query = `
	SELECT deleted_at IS NOT NULL FROM users WHERE id = ?
	`
	var isDeleted bool
	err := s.db.GetContext(ctx, &isDeleted, query, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, ErrNotFound
	}
	if err != nil {
		return false, err
	}
	return isDeleted, nil
}
