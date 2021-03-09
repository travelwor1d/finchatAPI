package controller

import (
	"net/http"

	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/sub"
)

type addCreditCardPayload struct {
	CardToken string `json:"cardToken" validate:"required"`
}

func (ctr *Ctr) AddCreditCard(c *fiber.Ctx) error {
	var p addCreditCardPayload
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, fiber.StatusBadRequest, "failed to parse body", err).Send(c)
	}
	v := validate.Struct(p)
	if !v.Validate() {
		return httperr.New(codes.Omit, http.StatusBadRequest, v.Errors.One()).Send(c)
	}

	id, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	user, err := ctr.store.GetUser(c.Context(), id)
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}

	if user.StripeID == nil {
		// Create stripe customer.
		params := &stripe.CustomerParams{
			Email:  &user.Email,
			Name:   stripe.String(user.FirstName + " " + user.LastName),
			Phone:  &user.Phone,
			Source: &stripe.SourceParams{Token: &p.CardToken},
		}
		custmr, err := customer.New(params)
		if err != nil {
			return errInternal.SetDetail(err).Send(c)
		}
		err = ctr.store.SetStripeID(c.Context(), user.ID, custmr.ID)
		if err != nil {
			return errInternal.SetDetail(err).Send(c)
		}
	} else {
		params := &stripe.CustomerParams{
			Source: &stripe.SourceParams{Token: &p.CardToken},
		}
		_, err := customer.Update(
			*user.StripeID,
			params,
		)
		if err != nil {
			return errInternal.SetDetail(err).Send(c)
		}
	}
	return c.JSON(fiber.Map{"success": true})
}

type createSubscriptionPayload struct {
	StripePriceID string `json:"stripePriceId" validate:"required"`
}

func (ctr *Ctr) CreateSubscription(c *fiber.Ctx) error {
	var p createSubscriptionPayload
	if err := c.BodyParser(&p); err != nil {
		return httperr.New(codes.Omit, fiber.StatusBadRequest, "failed to parse body", err).Send(c)
	}
	v := validate.Struct(p)
	if !v.Validate() {
		return httperr.New(codes.Omit, http.StatusBadRequest, v.Errors.One()).Send(c)
	}

	id, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	user, err := ctr.store.GetUser(c.Context(), id)
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}

	if user.StripeID == nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "not a stripe customer").Send(c)
	}

	params := &stripe.SubscriptionParams{
		Customer: user.StripeID,
		Items:    []*stripe.SubscriptionItemsParams{{Price: &p.StripePriceID}},
	}
	_, err = sub.New(params)
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(fiber.Map{"success": true})
}
