package controller

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

type webhookPayload struct {
	FirebaseID string    `json:"firebaseId"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (ctr *Ctr) CreateUserWebhook(c *fiber.Ctx) error {
	if c.Get("X-Webhook-Token") != appconfig.Config.WebhookToken {
		return httperr.New(codes.Omit, http.StatusUnauthorized, "Webhook token was not found").Send(c)
	}

	var p webhookPayload
	if err := c.BodyParser(&p); err != nil {
		log.Printf("Failed to parse body: %v", err)
		return httperr.New(codes.Omit, http.StatusBadRequest, "Failed to parse body", err).Send(c)
	}
	err := ctr.store.SetFirebaseIDByEmail(c.Context(), p.FirebaseID, p.Email)
	log.Printf("Some error here: %v", err)
	if errors.Is(err, store.ErrNotFound) {
		if err := ctr.tokenSvc.DeleteFirebaseUser(c.Context(), p.FirebaseID); err != nil {
			return errInternal.SetDetail(err).Send(c)
		}
		return httperr.New(codes.Omit, http.StatusBadRequest, "No user with such email").Send(c)
	}
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}

func (ctr *Ctr) DeleteUserWebhook(c *fiber.Ctx) error {
	if c.Get("X-Webhook-Token") != appconfig.Config.WebhookToken {
		return httperr.New(codes.Omit, http.StatusUnauthorized, "Webhook token was not found").Send(c)
	}

	var p webhookPayload
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Failed to parse body", err).Send(c)
	}
	err := ctr.store.DeleteUserByEmail(c.Context(), p.Email)
	if errors.Is(err, store.ErrNotFound) {
		return httperr.New(codes.Omit, http.StatusNotFound, "No user with such email").Send(c)
	}
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}
