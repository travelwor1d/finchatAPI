package controller

import (
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/token"
	"github.com/kevinburke/twilio-go"
)

type Ctr struct {
	store      *store.Store
	jwtManager *token.JWTManager
	verify     *twilio.VerifyPhoneNumberService
}

func New(s *store.Store, jw *token.JWTManager, v *twilio.VerifyPhoneNumberService) *Ctr {
	return &Ctr{s, jw, v}
}

func (ctr *Ctr) JWTManager() *token.JWTManager {
	return ctr.jwtManager
}
