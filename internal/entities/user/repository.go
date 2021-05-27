package user

import (
	"context"

	srequest "github.com/finchatapp/finchat-api/internal/entities/_shared/models/request"
	baserepo "github.com/finchatapp/finchat-api/internal/entities/_shared/repositories"
	"github.com/finchatapp/finchat-api/internal/entities/user/models/response"
)

type Repository interface {
	baserepo.BaseRepository
	SearchUsersCalcQuery(ctx context.Context, userID int, filters []srequest.CustomFilterItem) (string, baserepo.SortDefinitionFunc, []string, error)
	SearchUsersExecQuery(ctx context.Context, query string, params []string) ([]response.User, error)
}
