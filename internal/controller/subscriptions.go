package controller

import (
	"net/http"

	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/sub"
)

type addCreditCardPayload struct {
	CardToken string `json:"cardToken" validate:"required"`
}

func (ctr *Ctr) AddCreditCard(c *fiber.Ctx) error {
	var p addCreditCardPayload
	if err := c.BodyParser(&p); err != nil {
		return errParseBody.SetDetail(err).Send(c)
	}
	if v := validate.Struct(p); !v.Validate() {
		return httperr.New(codes.Omit, http.StatusBadRequest, v.Errors.One()).Send(c)
	}

	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}

	if user.StripeID == nil {
		// Create stripe customer.
		params := &stripe.CustomerParams{
			Email:  &user.Email,
			Name:   stripe.String(user.FirstName + " " + user.LastName),
			Phone:  &user.Phonenumber,
			Source: &stripe.SourceParams{Token: &p.CardToken},
		}
		custmr, err := customer.New(params)
		if err != nil {
			ctr.lr.LogError(err, c.Request())
			return errInternal.SetDetail(err).Send(c)
		}
		err = ctr.store.SetStripeID(c.Context(), user.ID, custmr.ID)
		if err != nil {
			ctr.lr.LogError(err, c.Request())
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
			ctr.lr.LogError(err, c.Request())
			return errInternal.SetDetail(err).Send(c)
		}
	}
	return sendSuccess(c)
}

type createSubscriptionPayload struct {
	StripePriceID string `json:"stripePriceId" validate:"required"`
}

type createSubscriptionPlanPayload struct {
	PriceUSD int64 `json:"priceUsd" validate:"required|uint"`
}

// CreateSubscriptionPlan creates stripe product and price for it.
func (ctr *Ctr) CreateSubscriptionPlan(c *fiber.Ctx) error {
	var p createSubscriptionPlanPayload
	if err := c.BodyParser(&p); err != nil {
		return errParseBody.SetDetail(err).Send(c)
	}
	if v := validate.Struct(p); !v.Validate() {
		return httperr.New(codes.Omit, http.StatusBadRequest, v.Errors.One()).Send(c)
	}

	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	prodParams := &stripe.ProductParams{
		Name: stripe.String(user.FirstName + " " + user.LastName + " subscription"),
	}
	prod, err := product.New(prodParams)
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}

	priceParams := &stripe.PriceParams{
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Product:  &prod.ID,
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String("month"),
		},
		UnitAmount: stripe.Int64(p.PriceUSD),
	}
	_, err = price.New(priceParams)
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}

	return sendSuccess(c)
}

func (ctr *Ctr) CreateSubscription(c *fiber.Ctx) error {
	var p createSubscriptionPayload
	if err := c.BodyParser(&p); err != nil {
		return errParseBody.SetDetail(err).Send(c)
	}
	if v := validate.Struct(p); !v.Validate() {
		return httperr.New(codes.Omit, http.StatusBadRequest, v.Errors.One()).Send(c)
	}

	user, httpErr := ctr.userFromCtx(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}

	if user.StripeID == nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "not a stripe customer").Send(c)
	}

	params := &stripe.SubscriptionParams{
		Customer: user.StripeID,
		Items:    []*stripe.SubscriptionItemsParams{{Price: &p.StripePriceID}},
	}
	_, err := sub.New(params)
	if err != nil {
		ctr.lr.LogError(err, c.Request())
		return errInternal.SetDetail(err).Send(c)
	}
	return sendSuccess(c)
}
