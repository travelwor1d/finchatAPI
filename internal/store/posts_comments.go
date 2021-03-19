package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/finchatapp/finchat-api/internal/model"
)

func (s *Store) ListPosts(ctx context.Context, p *Pagination) ([]*model.Post, error) {
	const query = `
	SELECT * FROM posts LIMIT ? OFFSET ?
	`
	var posts []*model.Post
	err := s.db.SelectContext(ctx, &posts, query, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *Store) GetPost(ctx context.Context, id int) (*model.Post, error) {
	const query = `
	SELECT * FROM posts WHERE id = ?
	`
	var post model.Post
	err := s.db.GetContext(ctx, &post, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *Store) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	const query = `
	INSERT INTO posts(title, content, image_urls, posted_by, published_at) VALUES (?, ?, ?, ?, ?)
	`
	result, err := s.db.ExecContext(ctx, query, post.Title, post.Content, post.ImageURLs, post.PostedBy, post.PublishedAt)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetPost(ctx, int(id))
}

func (s *Store) ListComments(ctx context.Context, p *Pagination) ([]*model.Comment, error) {
	const query = `
	SELECT * FROM comments LIMIT ? OFFSET ?
	`
	var comments []*model.Comment
	err := s.db.SelectContext(ctx, &comments, query, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *Store) GetComment(ctx context.Context, id int) (*model.Comment, error) {
	const query = `
	SELECT * FROM comments WHERE id = ?
	`
	var comment model.Comment
	err := s.db.GetContext(ctx, &comment, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (s *Store) CreateComment(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	const query = `
	INSERT INTO comments(post_id, content, posted_by, published_at) VALUES (?, ?, ?, ?)
	`
	result, err := s.db.ExecContext(ctx, query, comment.PostID, comment.Content, comment.PostedBy, comment.PublishedAt)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetComment(ctx, int(id))
}
