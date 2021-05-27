package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/finchatapp/finchat-api/internal/utils"
	"github.com/sirupsen/logrus"
)

func (s *repository) DeleteContact(ctx context.Context, userID, contactID int) error {
	query := squirrel.Delete("users_contacts").
		Where("user_id=?", userID).Where("contact_id=?", contactID).
		PlaceholderFormat(squirrel.Question).RunWith(s.MasterNode())

	_, err := query.ExecContext(ctx)
	if err != nil {
		q, p, _ := query.ToSql()
		logrus.Error(utils.SqlErrLogMsg(err, q, p))
		return err
	}

	return nil
}
