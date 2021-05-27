package repositories

import (
	"context"

	srequest "github.com/finchatapp/finchat-api/internal/entities/_shared/models/request"
	sresponse "github.com/finchatapp/finchat-api/internal/entities/_shared/models/response"
	baserepo "github.com/finchatapp/finchat-api/internal/entities/_shared/repositories"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db []*sqlx.DB
}

func (r *repository) MasterNode() *sqlx.DB {
	return r.db[0]
}

func (s *repository) CalcQueryWithSortAndOffset(ctx context.Context, q string, sorts []srequest.SortItem, pageSize, pageNumber int, sortDefinition baserepo.SortDefinitionFunc) (string, error) {
	return baserepo.CalcQueryWithSortAndOffset(ctx, q, sorts, pageSize, pageNumber, sortDefinition)
}

func (s *repository) GetCountByQuery(ctx context.Context, query string, params []string) (int, error) {
	return baserepo.GetCountByQuery(s, ctx, query, params)
}

func (s *repository) CalcPages(ctx context.Context, query string, params []string, reqPageNumber, reqPageSize, currentCount int) (sresponse.Paging, error) {
	return baserepo.CalcPages(s, ctx, query, params, reqPageNumber, reqPageSize, currentCount)
}

func New(
	db []*sqlx.DB,
) *repository {
	return &repository{
		db: db,
	}
}
