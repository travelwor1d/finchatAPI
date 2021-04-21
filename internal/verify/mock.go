package verify

import (
	"context"
)

type Mock struct {
}

func (Mock) Request(ctx context.Context, phoneNumber string) (string, error) {
	return "pending", nil
}

func (Mock) Verify(ctx context.Context, phoneNumber, code string) (string, error) {
	return "approved", nil
}
