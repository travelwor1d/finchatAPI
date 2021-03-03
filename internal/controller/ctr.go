package controller

import (
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/pkg/token"
)

type Ctr struct {
	store      *store.Store
	jwtManager *token.JWTManager
}

func New(s *store.Store, jw *token.JWTManager) *Ctr {
	return &Ctr{s, jw}
}
