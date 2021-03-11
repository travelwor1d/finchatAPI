package verify

import (
	"context"
)

type Mock struct {
}

func (Mock) Request(ctx context.Context, phone string) (string, error) {
	return "pending", nil
}

func (Mock) Verify(ctx context.Context, phone, code string) (string, error) {
	return "approved", nil
}
