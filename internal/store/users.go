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

func (s *Store) GetUserByFirebaseID(ctx context.Context, uid string) (*model.User, error) {
	const query = `
	SELECT * FROM users WHERE firebase_id = ?
	`
	var user model.User
	err := s.db.GetContext(ctx, &user, query, uid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) GetUser(ctx context.Context, id int) (*model.User, error) {
	const query = `
	SELECT * FROM users WHERE id = ?
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
	SELECT * FROM users WHERE email = ?
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
	SELECT * FROM verified_active_users
		WHERE (
			lower(first_name) LIKE '%s' OR
			lower(last_name) LIKE '%s'
		) AND user_type IN (%s)
	LIMIT ? OFFSET ?
	`, "%"+searchInput+"%", "%"+searchInput+"%", userTypes)
	var users []*model.User
	err := s.db.SelectContext(ctx, &users, query, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Store) CreateUser(ctx context.Context, user *model.User, inviteCode ...string) (*model.User, error) {
	const query = `
	INSERT INTO users(first_name, last_name, phone_number, country_code, email, username, user_type, profile_avatar)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	tx, err := s.Begin()
	if err != nil {
		return nil, err
	}

	result, err := tx.ExecContext(ctx, query,
		user.FirstName, user.LastName, user.Phonenumber, user.CountryCode, user.Email, user.Username, user.Type, user.ProfileAvatar,
	)
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

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) UpdateUser(ctx context.Context, userID int, firstName, lastName, username, profileAvatar *string) (*model.User, error) {
	isDeleted, err := s.isUserDeleted(ctx, userID)
	if err != nil {
		return nil, err
	}
	if isDeleted {
		return nil, ErrUserDeleted
	}
	const query = `
	UPDATE verified_active_users SET
		first_name = coalesce(?, first_name),
		last_name = coalesce(?, last_name),
		username = coalesce(?, username),
		profile_avatar = coalesce(?, profile_avatar)
	WHERE id = ?
	`
	result, err := s.db.ExecContext(ctx, query, firstName, lastName, username, profileAvatar, userID)
	if err != nil {
		return nil, err
	}
	r, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return nil, ErrNoRowsAffected
	}
	return s.GetUser(ctx, userID)
}

func (s *Store) SoftDeleteUser(ctx context.Context, userID int) error {
	isDeleted, err := s.isUserDeleted(ctx, userID)
	if err != nil {
		return err
	}
	if isDeleted {
		// Return error that user was already "soft deleted".
		return ErrUserDeleted
	}
	const query = `
	UPDATE verified_active_users SET
		deleted_at = now()
	WHERE id = ?
	`
	result, err := s.db.ExecContext(ctx, query, userID)
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

func (s *Store) DeleteUserByEmail(ctx context.Context, email string) error {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	const query = `
	DELETE FROM users WHERE id = ?
	`
	result, err := s.db.ExecContext(ctx, query, user.ID)
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

func (s *Store) UndeleteUser(ctx context.Context, userID int) error {
	isDeleted, err := s.isUserDeleted(ctx, userID)
	if err != nil {
		return err
	}
	if !isDeleted {
		// Return error that user has not been "soft deleted", yet.
		return ErrUserNotDeleted
	}
	const query = `
	UPDATE users SET
		deleted_at = NULL
	WHERE id = ?
	`
	result, err := s.db.ExecContext(ctx, query, userID)
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

func (s *Store) SetStripeID(ctx context.Context, userID int, stripeID string) error {
	isDeleted, err := s.isUserDeleted(ctx, userID)
	if err != nil {
		return err
	}
	if isDeleted {
		return ErrUserDeleted
	}
	const query = `
	UPDATE verified_active_users SET
		stripe_id = ?
	WHERE id = ?
	`
	result, err := s.db.ExecContext(ctx, query, stripeID, userID)
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

func (s *Store) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	const query = `
	SELECT EXISTS (SELECT 1 FROM users WHERE email = ?)
	`
	var exists bool
	err := s.db.GetContext(ctx, &exists, query, email)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Store) IsPhoneNumberTaken(ctx context.Context, phoneNumber string) (bool, error) {
	const query = `
	SELECT EXISTS (SELECT 1 FROM users WHERE phone_number = ?)
	`
	var exists bool
	err := s.db.GetContext(ctx, &exists, query, phoneNumber)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Store) IsUsernameTaken(ctx context.Context, username string) (bool, error) {
	const query = `
	SELECT EXISTS (SELECT 1 FROM users WHERE username = ?)
	`
	var exists bool
	err := s.db.GetContext(ctx, &exists, query, username)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Store) SetFirebaseIDByEmail(ctx context.Context, uid, email string) error {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	const query = `
	UPDATE users SET
		firebase_id = ?,
		is_active = true
	WHERE id = ? AND deleted_at IS NULL
	`
	result, err := s.db.ExecContext(ctx, query, uid, user.ID)
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

func (s *Store) SetVerifiedUser(ctx context.Context, id int) error {
	const query = `
	UPDATE users SET
		is_verified = true
	WHERE id = ? AND is_active AND deleted_at IS NULL
	`
	result, err := s.db.ExecContext(ctx, query, id)
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
