package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/finchatapp/finchat-api/internal/messaging"
	"github.com/finchatapp/finchat-api/internal/model"
)

func (s *Store) ListThreads(ctx context.Context, userID int, p *Pagination) ([]*model.Thread, error) {
	const query = `
	SELECT threads.* FROM threads
		JOIN thread_participants ON threads.id = thread_participants.thread_id
	WHERE thread_participants.user_id = ?
	LIMIT ? OFFSET ?
	`
	var threads []*model.Thread
	err := s.db.SelectContext(ctx, &threads, query, userID, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}
	return threads, nil
}

func (s *Store) GetThread(ctx context.Context, id int) (*model.Thread, error) {
	const query = `
	SELECT * FROM threads WHERE id = ?
	`
	var thread model.Thread
	err := s.db.GetContext(ctx, &thread, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func (s *Store) CreateThread(ctx context.Context, thread *model.Thread, participants []int) (*model.Thread, error) {
	if len(participants) < 1 {
		return nil, errors.New("invalid number of thread participants")
	}
	const query = `
	INSERT INTO threads(title, thread_type) VALUES (?, ?)
	`
	tx, err := s.Begin()
	if err != nil {
		return nil, err
	}

	result, err := tx.ExecContext(ctx, query, thread.Title, thread.Type)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	q := `INSERT INTO thread_participants(user_id, thread_id) VALUES (?, ?)`
	for _, participant := range participants {
		_, err = tx.ExecContext(ctx, q, participant, id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	thread, err = s.WithTx(tx).GetThread(ctx, int(id))
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return thread, nil
}

// store.Store should implement messaging.MessageSaver interface.
var _ messaging.MessageSaver = (*Store)(nil)

func (s *Store) SaveMessage(ctx context.Context, msg *model.Message) error {
	const query = `
		INSERT INTO thread_messages(thread_id, sender_id, message_type, message, timestamp)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := s.db.ExecContext(ctx, query, msg.ThreadID, msg.SenderID, msg.Type, msg.Message, msg.Timestamp)
	return err
}
