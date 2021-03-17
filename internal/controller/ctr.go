package controller

import (
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/internal/upload"
	"github.com/finchatapp/finchat-api/internal/verify"
	"github.com/finchatapp/finchat-api/pkg/token"
)

type Ctr struct {
	store      *store.Store
	jwtManager *token.JWTManager
	verify     verify.Verifier
	upload     *upload.Client
}

func New(s *store.Store, jw *token.JWTManager, v verify.Verifier, u *upload.Client) *Ctr {
	return &Ctr{s, jw, v, u}
}

func (ctr *Ctr) JWTManager() *token.JWTManager {
	return ctr.jwtManager
}
