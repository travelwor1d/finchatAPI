package upload

import (
	"context"
	"fmt"
	"io"

	"github.com/finchatapp/finchat-api/internal/model"
	"github.com/finchatapp/finchat-api/pkg/unique"
)

type Mock struct {
}

func (m Mock) ProfiveAvatarFileName(user *model.User, ext string) string {
	return fmt.Sprintf("profile-avatar-%s-%s-%s%s", user.FirstName, user.LastName, unique.New(12), ext)
}

func (m Mock) Upload(ctx context.Context, filename string, r io.Reader) error {
	return nil
}
