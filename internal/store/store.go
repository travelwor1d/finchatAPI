package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/finchatapp/finchat-api/internal/model"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Pagination struct {
	Limit, Offset int
}

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

func (s *Store) SetVerifiedUser(ctx context.Context, id int) error {
	const query = `
	UPDATE users SET
		verified = true
	WHERE id = ?
	`
	_, err := s.db.ExecContext(ctx, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
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
