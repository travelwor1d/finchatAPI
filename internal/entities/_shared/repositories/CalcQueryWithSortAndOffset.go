package repositories

import (
	"context"
	"fmt"
	"strings"

	srequest "github.com/finchatapp/finchat-api/internal/entities/_shared/models/request"
	"github.com/finchatapp/finchat-api/internal/utils"
)

func CalcQueryWithSortAndOffset(ctx context.Context, q string, sorts []srequest.SortItem, pageSize, pageNumber int, sortDefinition SortDefinitionFunc) (string, error) {
	offset := (pageNumber - 1) * pageSize

	orderStatements := make([]string, 0)

	for _, sort := range sorts {
		if sort.Dir == "descending" {
			sort.Dir = "desc"
		} else if sort.Dir == "ascending" {
			sort.Dir = "asc"
		} else if !utils.SliceStringsContains([]string{"asc", "desc"}, sort.Dir) {
			sort.Dir = ""
		}
		if sort.Dir != "" {
			sortString := sortDefinition(sort)
			if sortString != "" {
				orderStatements = append(orderStatements, sortString)
			}
		}
	}
	if len(orderStatements) > 0 {
		q = fmt.Sprintf("%s order by %s", q, strings.Join(orderStatements, " , "))
	}
	q = fmt.Sprintf("%s limit %d offset %d", q, pageSize, offset)

	return q, nil
}
