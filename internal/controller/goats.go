package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (ctr *Ctr) InviteGoat(c *fiber.Ctx) error {
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	if user.Type == "USER" {
		code, err := ctr.store.CreateGoatInviteCode(c.Context(), user.ID)
		if err != nil {
			ctr.lr.LogError(err, c.Request())
			return errInternal.SetDetail(err).Send(c)
		}
		return c.JSON(fiber.Map{"inviteCode": code})
	} else if user.Type == "GOAT" {
		return errInternal.SetDetail("goat not supported").Send(c)
	}
	return errInternal.SetDetail("unknown user type").Send(c)
}

func (ctr *Ctr) VerifyInviteCode(c *fiber.Ctx) error {
	inviteCode := c.Params("inviteCode")
	status, found, err := ctr.store.GetInviteCodeStatus(c.Context(), inviteCode)
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	if !found {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"valid": false, "reason": "not found"})
	}
	if status != "ACTIVE" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"valid": false, "reason": "status: " + status})
	}
	return c.JSON(fiber.Map{"valid": true})
}
