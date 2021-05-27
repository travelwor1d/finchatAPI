package user

import (
	"context"

	srequest "github.com/finchatapp/finchat-api/internal/entities/_shared/models/request"
	sresponse "github.com/finchatapp/finchat-api/internal/entities/_shared/models/response"
	"github.com/finchatapp/finchat-api/internal/entities/user/models/response"
)

type Usecase interface {
	SearchUsers(ctx context.Context, userID int, req srequest.NewGridList) (*sresponse.Paging, []response.User, error)
}
