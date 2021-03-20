package messaging

import "context"

type Mock struct {
}

func (m Mock) User(id int) Messager {
	return m
}

func (m Mock) Register(ctx context.Context, firstName string, lastName string, email string) error {
	return nil
}
