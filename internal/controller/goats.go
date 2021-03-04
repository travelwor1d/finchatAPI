package controller

import (
	"net/http"

	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func (ctr *Ctr) InviteGoat(c *fiber.Ctx) error {
	id, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	user, err := ctr.store.GetUser(c.Context(), id)
	if err != nil {
		return httperr.New(codes.Omit, http.StatusInternalServerError, "").Send(c)
	}
	if user.Type == "USER" {

	}
	return nil
}
