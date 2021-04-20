package verify

import (
	"context"
)

type Mock struct {
}

func (Mock) Request(ctx context.Context, phonenumber string) (string, error) {
	return "pending", nil
}

func (Mock) Verify(ctx context.Context, phonenumber, code string) (string, error) {
	return "approved", nil
}
