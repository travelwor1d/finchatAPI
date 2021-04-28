package token

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

type Service interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
	DeleteFirebaseUser(ctx context.Context, uid string) error
}

type svc struct {
	authClient *auth.Client
}

func NewService(c *auth.Client) Service {
	return &svc{c}
}

func (s *svc) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return s.authClient.VerifyIDToken(ctx, idToken)
}

func (s *svc) DeleteFirebaseUser(ctx context.Context, uid string) error {
	return s.authClient.DeleteUser(ctx, uid)
}
