package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/finchatapp/finchat-api/internal/entities/user/models/response"
	"github.com/finchatapp/finchat-api/internal/utils"

	"github.com/sirupsen/logrus"
)

type NamedParam struct {
}

func (s *repository) SearchUsersExecQuery(ctx context.Context, query string, params []string) ([]response.User, error) {
	res := make([]response.User, 0)

	newParams := make(map[string]interface{}, len(params))
	for i, v := range params {
		newParams[fmt.Sprintf("param_value_%d", i+1)] = v
	}
	rows, err := s.MasterNode().NamedQuery(query, newParams)

	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		logrus.Error(utils.SqlErrLogMsg(err, query, nil))
		return res, err
	}

	for rows.Next() {
		item := response.User{}
		err := rows.Scan(
			&item.ID,
			&item.IsActive,
			&item.FirstName,
			&item.LastName,
			&item.Phonenumber,
			&item.CountryCode,
			&item.Email,
			&item.IsVerified,
			&item.Type,
			&item.ProfileAvatar,
			&item.LastSeen,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.Username,
			&item.IsContact,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}
