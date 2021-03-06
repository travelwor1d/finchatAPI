package controller

import (
	"net/http"

	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func (ctr *Ctr) RequestVerification(c *fiber.Ctx) error {
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	if user.IsVerified {
		return httperr.New(codes.Omit, http.StatusBadRequest, "already verified").Send(c)
	}
	status, err := ctr.verify.Request(c.Context(), user.Phonenumber)
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(fiber.Map{"verificationStatus": status})
}

type VerifyPayload struct {
	Code string `json:"code"`
}

func (ctr *Ctr) Verify(c *fiber.Ctx) error {
	var p VerifyPayload
	if err := c.BodyParser(&p); err != nil {
		return errParseBody.SetDetail(err).Send(c)
	}
	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	status, err := ctr.verify.Verify(c.Context(), user.Phonenumber, p.Code)
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	if status == "pending" {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid verification code").Send(c)
	}
	if status == "approved" {
		if err = ctr.store.SetVerifiedUser(c.Context(), user.ID); err != nil {
			ctr.lr.LogError(err, c.Request())
			return errInternal.SetDetail(err).Send(c)
		}
		return c.JSON(fiber.Map{"verificationStatus": status})
	}
	return c.JSON(fiber.Map{"verificationStatus": status})
}
