package repositories

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/finchatapp/finchat-api/internal/entities/contact/models"
	"github.com/finchatapp/finchat-api/internal/utils"
	"github.com/sirupsen/logrus"
)

func (s *repository) GetUserByContactUUID(ctx context.Context, uuid string) (*models.Contact, error) {
	res := models.Contact{}

	query := squirrel.Select(
		"u.id",
		"u.first_name",
		"u.last_name",
		"u.phone_number",
		"u.country_code",
		"u.email",
		"u.user_type",
		"u.profile_avatar",
		"u.last_seen",
		"u.created_at",
		"u.updated_at",
	).From("core.users_contacts uc").
		Join("core.users u on u.id = uc.contact_id").
		Where("uc.uuid = ?", uuid).
		PlaceholderFormat(squirrel.Question).RunWith(s.MasterNode())
	row := query.QueryRowContext(ctx)
	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Phonenumber,
		&res.CountryCode,
		&res.Email,
		&res.Type,
		&res.ProfileAvatar,
		&res.LastSeen,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		q, p, _ := query.ToSql()
		logrus.Error(utils.SqlErrLogMsg(err, q, p))
		return nil, err
	}

	return &res, err
}
