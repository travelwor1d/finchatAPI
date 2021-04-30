package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrNoRowsAffected = errors.New("no rows affected")
	ErrAlreadyExists  = errors.New("already exists")
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
