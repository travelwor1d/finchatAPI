package service

import "context"

func (svc *usecase) DeleteContact(ctx context.Context, userID, contactID int) error {
	err := svc.repo.DeleteContact(ctx, userID, contactID)
	if err != nil {
		return err
	}
	return nil
}
