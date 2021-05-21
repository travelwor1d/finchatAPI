package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"github.com/finchatapp/finchat-api/internal/model"
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

var space = regexp.MustCompile(`\s+`)
var spaceOrPlus = regexp.MustCompile(`[ +]`)

func (s *Store) SearchUsers(ctx context.Context, userID int, searchInput, userTypes string, ignoreContacts bool, p *Pagination) ([]*model.User, error) {
	// Remove all duplicate whitespace.
	phoneNumber := spaceOrPlus.ReplaceAllString(searchInput, "")
	// Remove all duplicate whitespace.
	searchInput = space.ReplaceAllString(searchInput, " ")
	query := fmt.Sprintf(`
	SELECT u.*, c.id IS NOT NULL AS is_contact FROM verified_active_users u
	LEFT JOIN users_contacts c ON c.user_id = %d AND u.id = c.contact_id
		WHERE (
			username LIKE '%s' OR
			lower(concat(first_name, ' ', last_name)) LIKE '%s' OR
			lower(concat(last_name, ' ', first_name)) LIKE '%s' OR
			email = '%s' OR
			phone_number LIKE '%s'
		) AND user_type IN (%s)
	ORDER BY first_name ASC, last_name ASC
	LIMIT ? OFFSET ?
	`, userID, "%"+searchInput+"%", "%"+searchInput+"%", "%"+searchInput+"%", searchInput, "%"+phoneNumber+"%", userTypes)
	ignoreContactsQuery := fmt.Sprintf(`
	SELECT u.* FROM verified_active_users u
	JOIN users_contacts c ON u.id <> c.contact_id OR c.user_id <> 5
		WHERE (
			username LIKE '%s' OR
			lower(concat(first_name, ' ', last_name)) LIKE '%s' OR
			lower(concat(last_name, ' ', first_name)) LIKE '%s' OR
			email = '%s' OR
			phone_number = '%s'
		) AND user_type IN (%s)
	ORDER BY first_name ASC, last_name ASC
	LIMIT ? OFFSET ?
	`, "%"+searchInput+"%", "%"+searchInput+"%", "%"+searchInput+"%", searchInput, "%"+phoneNumber+"%", userTypes)
	var users []*model.User
	var err error
	if ignoreContacts {
		err = s.db.SelectContext(ctx, &users, ignoreContactsQuery, p.Limit, p.Offset)
	} else {
		err = s.db.SelectContext(ctx, &users, query, p.Limit, p.Offset)
	}
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Store) UpsertUser(ctx context.Context, user *model.User, inviteCode ...string) (*model.User, error) {
	const query = `
	INSERT INTO users(first_name, last_name, phone_number, country_code, email, username, user_type, profile_avatar)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
		first_name = VALUES(first_name),
		last_name = VALUES(last_name),
		phone_number = VALUES(phone_number),
		country_code = VALUES(country_code),
		email = VALUES(email),
		username = VALUES(username),
		user_type = VALUES(user_type),
		profile_avatar = VALUES(profile_avatar)
	`
	tx, err := s.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, query,
		user.FirstName, user.LastName, user.Phonenumber, user.CountryCode, user.Email, user.Username, user.Type, user.ProfileAvatar,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	user, err = s.WithTx(tx).GetUserByEmail(ctx, user.Email)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if user.Type == "GOAT" {
		if len(inviteCode) != 1 {
			return nil, errors.New("invalid usage of store.CreateUser")
		}
		err := s.WithTx(tx).UseInviteCode(ctx, inviteCode[0], user.ID)
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
	SELECT EXISTS (SELECT 1 FROM active_users WHERE email = ?)
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
	SELECT EXISTS (SELECT 1 FROM active_users WHERE phone_number = ?)
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
