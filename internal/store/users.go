package store

import (
	"context"
	"errors"

	"github.com/finchatapp/finchat-api/internal/model"
	"github.com/go-sql-driver/mysql"
)

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
