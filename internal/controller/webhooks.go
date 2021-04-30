package controller

import (
	"errors"
	"fmt"
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
		return httperr.New(codes.Omit, http.StatusBadRequest, "Failed to parse body", err).Send(c)
	}
	err := ctr.store.SetActiveUserByEmail(c.Context(), p.Email)
	if errors.Is(err, store.ErrNotFound) {
		if err := ctr.tokenSvc.DeleteFirebaseUser(c.Context(), p.FirebaseID); err != nil {
			// TODO: report an error.
			fmt.Printf("Failed to delete firebase user: %v\n", err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		return c.SendStatus(http.StatusBadRequest)
	}
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendStatus(http.StatusCreated)
}

func (ctr *Ctr) DeleteUserWebhook(c *fiber.Ctx) error {
	if c.Get("X-Webhook-Token") != appconfig.Config.WebhookToken {
		return httperr.New(codes.Omit, http.StatusUnauthorized, "Webhook token was not found").Send(c)
	}

	var p webhookPayload
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Failed to parse body", err).Send(c)
	}
	err := ctr.store.SetActiveUserByEmail(c.Context(), p.Email)
	if errors.Is(err, store.ErrNotFound) {
		return c.SendStatus(http.StatusNotFound)
	}
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendStatus(http.StatusOK)
}
