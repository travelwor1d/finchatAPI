package repositories

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/finchatapp/finchat-api/internal/utils"
	"github.com/sirupsen/logrus"
)

func (s *repository) ValidateUserIdExists(ctx context.Context, userID int) (bool, error) {
	res := 0
	q := squirrel.Select(
		"id",
	).From("core.users").Where("id = ?", userID).
		PlaceholderFormat(squirrel.Question).RunWith(s.MasterNode())
	row := q.QueryRowContext(ctx)
	err := row.Scan(
		&res,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return res > 0, nil
		}
		q1, p, _ := q.ToSql()
		logrus.Error(utils.SqlErrLogMsg(err, q1, p))
		return res > 0, err
	}
	return res > 0, nil
}
