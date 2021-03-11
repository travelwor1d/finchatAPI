package controller

import (
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/internal/verify"
	"github.com/finchatapp/finchat-api/pkg/token"
)

type Ctr struct {
	store      *store.Store
	jwtManager *token.JWTManager
	verify     verify.Verifier
}

func New(s *store.Store, jw *token.JWTManager, v verify.Verifier) *Ctr {
	return &Ctr{s, jw, v}
}

func (ctr *Ctr) JWTManager() *token.JWTManager {
	return ctr.jwtManager
}
