package controller

import (
	"github.com/finchatapp/finchat-api/internal/messaging"
	"github.com/finchatapp/finchat-api/internal/store"
	"github.com/finchatapp/finchat-api/internal/token"
	"github.com/finchatapp/finchat-api/internal/upload"
	"github.com/finchatapp/finchat-api/internal/verify"
)

type Ctr struct {
	store    *store.Store
	tokenSvc token.Service
	verify   verify.Verifier
	upload   upload.Uploader
	msg      messaging.Messager
}

func New(s *store.Store, t token.Service, v verify.Verifier, u upload.Uploader, msg messaging.Messager) *Ctr {
	return &Ctr{s, t, v, u, msg}
}

func (ctr *Ctr) TokenSvc() token.Service {
	return ctr.tokenSvc
}
