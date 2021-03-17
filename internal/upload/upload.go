package upload

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/finchatapp/finchat-api/internal/model"
	"github.com/finchatapp/finchat-api/pkg/unique"
)

type Client struct {
	bucket *storage.BucketHandle
}

func New(b *storage.BucketHandle) *Client {
	return &Client{b}
}

func (c *Client) ProfiveAvatarFileName(user *model.User, ext string) string {
	return fmt.Sprintf("profile-avatar-%s-%s-%s%s", user.FirstName, user.LastName, unique.New(12), ext)
}

func (c *Client) UploadProfileAvatar(ctx context.Context, filename string, r io.Reader) error {
	wc := c.bucket.Object(filename).NewWriter(ctx)
	if _, err := io.Copy(wc, r); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}
