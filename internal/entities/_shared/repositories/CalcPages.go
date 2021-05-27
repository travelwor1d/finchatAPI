package repositories

import (
	"context"

	sresponse "github.com/finchatapp/finchat-api/internal/entities/_shared/models/response"
)

func CalcPages(s BaseRepository, ctx context.Context, query string, params []string, reqPageNumber, reqPageSize, currentCount int) (sresponse.Paging, error) {
	res := sresponse.Paging{}
	totalCount, err := s.GetCountByQuery(ctx, query, params)
	if err != nil {
		return res, err
	}
	res.TotalCount = totalCount
	res.CurrentCount = currentCount
	res.Page = reqPageNumber
	res.Size = reqPageSize
	return res, nil
}
