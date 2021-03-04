package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/finchatapp/finchat-api/internal/model"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	GetContext(context.Context, interface{}, string, ...interface{}) error
	SelectContext(context.Context, interface{}, string, ...interface{}) error
}

type Store struct {
	db   DBTX
	conn *sqlx.DB
}

func New(db *sqlx.DB) *Store {
	return &Store{db, db}
}

func (s *Store) Begin() (*sqlx.Tx, error) {
	return s.conn.Beginx()
}

func (s *Store) WithTx(tx *sqlx.Tx) *Store {
	return &Store{tx, s.conn}
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

func (s *Store) GetUserCredsByEmail(ctx context.Context, email string) (*model.Creds, error) {
	const query = `
	SELECT c.* FROM credentials c JOIN users u ON c.user_id = u.id WHERE email = ?
	`
	var creds model.Creds
	err := s.db.GetContext(ctx, &creds, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &creds, nil
}

func (s *Store) CreateUser(ctx context.Context, user *model.User, password string) (*model.User, error) {
	const query = `
	INSERT INTO users(first_name, last_name, phone, email, user_type, profile_avatar)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	tx, err := s.Begin()
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx, query, user.FirstName, user.LastName, user.Phone, user.Email, user.Type, user.ProfileAvatar)
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

	user, err = s.WithTx(tx).GetUserByEmail(ctx, user.Email)
	if err != nil {
		tx.Rollback()
		return nil, err
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

func (s *Store) SetPassword(ctx context.Context, id int, password string) error {
	const query = `
	INSERT INTO credentials(hash, user_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE hash = ?
	`
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, query, string(hash), id, string(hash))
	return err
}
