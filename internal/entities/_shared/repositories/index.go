package repositories

import (
	"context"

	srequest "github.com/finchatapp/finchat-api/internal/entities/_shared/models/request"
	sresponse "github.com/finchatapp/finchat-api/internal/entities/_shared/models/response"
	"github.com/jmoiron/sqlx"
)

type SortDefinitionFunc func(sort srequest.SortItem) string

type BaseRepository interface {
	MasterNode() *sqlx.DB
	CalcQueryWithSortAndOffset(ctx context.Context, q string, sorts []srequest.SortItem, pageSize, pageNumber int, sortDefinition SortDefinitionFunc) (string, error)
	CalcPages(ctx context.Context, query string, params []string, reqPageNumber, reqPageSize, currentCount int) (sresponse.Paging, error)
	GetCountByQuery(ctx context.Context, query string, params []string) (int, error)
}
