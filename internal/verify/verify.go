package verify

import (
	"context"
	"net/url"

	"github.com/kevinburke/twilio-go"
)

type Verifier interface {
	Request(ctx context.Context, phonenumber string) (string, error)
	Verify(ctx context.Context, phonenumber, code string) (string, error)
}

type verify struct {
	svc *twilio.VerifyPhoneNumberService
	id  string
}

func New(s *twilio.VerifyPhoneNumberService, id string) Verifier {
	return &verify{s, id}
}

func (v *verify) Request(ctx context.Context, phonenumber string) (string, error) {
	resp, err := v.svc.Create(ctx, v.id, url.Values{"To": []string{phonenumber}, "Channel": []string{"sms"}})
	return resp.Status, err
}

func (v *verify) Verify(ctx context.Context, phonenumber, code string) (string, error) {
	resp, err := v.svc.Check(ctx, v.id, url.Values{"To": []string{phonenumber}, "Code": []string{code}})
	return resp.Status, err
}
