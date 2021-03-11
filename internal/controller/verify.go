package controller

import (
	"net/http"
	"net/url"

	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func (ctr *Ctr) RequestVerification(c *fiber.Ctx) error {
	id, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	user, err := ctr.store.GetUser(c.Context(), id)
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	if user.Verified {
		return httperr.New(codes.Omit, http.StatusBadRequest, "already verified").Send(c)
	}
	v, err := ctr.verify.Create(c.Context(), viper.GetString("twilio.verify"), url.Values{"To": []string{user.Phone}, "Channel": []string{"sms"}})
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(fiber.Map{"verificationStatus": v.Status})
}

type VerifyPayload struct {
	Code string `json:"code"`
}

func (ctr *Ctr) Verify(c *fiber.Ctx) error {
	var p VerifyPayload
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "failed to parse body", err).Send(c)
	}
	id, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	user, err := ctr.store.GetUser(c.Context(), id)
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	v, err := ctr.verify.Check(c.Context(), viper.GetString("twilio.verify"), url.Values{"To": []string{user.Phone}, "Code": []string{p.Code}})
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	if v.Status == "pending" {
		return httperr.New(codes.Omit, http.StatusBadRequest, "invalid verification code").Send(c)
	}
	if v.Status == "approved" {
		if err = ctr.store.SetVerifiedUser(c.Context(), user.ID); err != nil {
			return errInternal.SetDetail(err).Send(c)
		}
		return c.JSON(fiber.Map{"verificationStatus": v.Status})
	}
	return c.JSON(fiber.Map{"verificationStatus": v.Status})
}
