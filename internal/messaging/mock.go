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

func (m Mock) SendMessage(ctx context.Context, threadID, senderID int, message string) (int64, error) {
	return 0, nil
}
