package controller

import (
	"net/http"
	"time"

	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

type createUserWebhookPayload struct {
	FirebaseID string    `json:"firebaseId"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (ctr *Ctr) CreateUserWebhook(c *fiber.Ctx) error {
	if c.Get("X-Webhook-Token") != appconfig.Config.WebhookToken {
		return httperr.New(codes.Omit, http.StatusUnauthorized, "Webhook token was not found").Send(c)
	}

	var p createUserWebhookPayload
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "Failed to parse body", err).Send(c)
	}

	if err := ctr.store.CreateFirebaseUser(c.Context(), p.FirebaseID, p.Email, p.CreatedAt); err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.SendStatus(http.StatusCreated)
}
