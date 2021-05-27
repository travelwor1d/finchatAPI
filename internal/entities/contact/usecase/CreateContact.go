package service

import (
	"context"

	"github.com/finchatapp/finchat-api/internal/entities/contact/models"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/internal/utils"
)

func (svc *usecase) CreateContact(ctx context.Context, userID, contactID int) (*models.Contact, error) {
	contactUserExists, err := svc.repo.ValidateUserIdExists(ctx, contactID)
	if err != nil {
		return nil, err
	}
	if !contactUserExists {
		return nil, store.ErrNotFound
	}

	newUUID := utils.NewUUID()
	err = svc.repo.CreateContact(ctx, userID, contactID, newUUID)
	if err != nil {
		return nil, err
	}
	res, err := svc.repo.GetUserByContactUUID(ctx, newUUID)
	if err != nil {
		return nil, err
	}
	return res, nil
}
