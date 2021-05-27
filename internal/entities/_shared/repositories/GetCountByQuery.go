package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/finchatapp/finchat-api/internal/utils"
	"github.com/sirupsen/logrus"
)

func GetCountByQuery(s BaseRepository, ctx context.Context, query string, params []string) (int, error) {
	res := 0
	newParams := make(map[string]interface{}, len(params))
	for i, v := range params {
		newParams[fmt.Sprintf("param_value_%d", i+1)] = v
	}
	q := fmt.Sprintf("select coalesce(count(*),0) from (%s) t", query)
	rows, err := s.MasterNode().NamedQuery(q, newParams)

	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		logrus.Error(utils.SqlErrLogMsg(err, q, nil))
		return res, err
	}

	for rows.Next() {
		item := 0
		err := rows.Scan(
			&item,
		)
		if err != nil {
			return 0, err
		}
		res = item
	}
	return res, nil
}
