package service

import (
	"context"

	srequest "github.com/finchatapp/finchat-api/internal/entities/_shared/models/request"
	sresponse "github.com/finchatapp/finchat-api/internal/entities/_shared/models/response"
	"github.com/finchatapp/finchat-api/internal/entities/user/models/response"
)

func (svc *usecase) SearchUsers(ctx context.Context, userID int, req srequest.NewGridList) (*sresponse.Paging, []response.User, error) {
	query, sortDef, params, err := svc.repo.SearchUsersCalcQuery(ctx, userID, req.CustomFilters)
	if err != nil {
		return nil, nil, err
	}
	queryWithSortAndOffset, err := svc.repo.CalcQueryWithSortAndOffset(ctx, query, req.Sorts, req.PageSize, req.PageNumber, sortDef)
	if err != nil {
		return nil, nil, err
	}

	resItems, err := svc.repo.SearchUsersExecQuery(ctx, queryWithSortAndOffset, params)
	if err != nil {
		return nil, nil, err
	}
	paging, err := svc.repo.CalcPages(ctx, query, params, req.PageNumber, req.PageSize, len(resItems))
	if err != nil {
		return nil, nil, err
	}
	return &paging, resItems, nil
}
