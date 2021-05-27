package contact

import (
	"context"

	baserepo "github.com/finchatapp/finchat-api/internal/entities/_shared/repositories"
	"github.com/finchatapp/finchat-api/internal/entities/contact/models"
)

type Repository interface {
	baserepo.BaseRepository

	CreateContact(ctx context.Context, userID, contactID int, uuid string) error
	DeleteContact(ctx context.Context, userID, contactID int) error
	ValidateUserIdExists(ctx context.Context, userID int) (bool, error)

	GetUserByContactUUID(ctx context.Context, uuid string) (*models.Contact, error)
}
