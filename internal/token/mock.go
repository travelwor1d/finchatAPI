package token

import (
	"context"
	"os"

	"firebase.google.com/go/v4/auth"
)

type Mock struct {
}

func (m Mock) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return &auth.Token{UID: os.Getenv("TEST_USER_ID")}, nil
}

func (m Mock) DeleteFirebaseUser(ctx context.Context, uid string) error {
	return nil
}
