package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/finchatapp/finchat-api/internal/model"
	"github.com/go-sql-driver/mysql"
)

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
	SELECT * FROM users
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
	fmt.Println(users, err)
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

func (s *Store) SetStripeID(ctx context.Context, userID int, stripeID string) error {
	const query = `
	UPDATE users SET
		stripe_id = ?
	WHERE id = ?
	`
	_, err := s.db.ExecContext(ctx, query, stripeID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	return nil
}
