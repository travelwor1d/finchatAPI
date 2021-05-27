package controller

import (
	"github.com/finchatapp/finchat-api/internal/entities/contact"
	"github.com/finchatapp/finchat-api/internal/entities/user"
	"github.com/finchatapp/finchat-api/internal/logerr"
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
	lr       *logerr.Logerr
	contact  contact.Usecase
	user     user.Usecase
}

func New(s *store.Store,
	contactUsecase contact.Usecase,
	userUsecase user.Usecase,
	t token.Service,
	v verify.Verifier,
	u upload.Uploader,
	msg messaging.Messager,
	lr *logerr.Logerr,
) *Ctr {
	return &Ctr{
		store:    s,
		contact:  contactUsecase,
		user:     userUsecase,
		tokenSvc: t,
		verify:   v,
		upload:   u,
		msg:      msg,
		lr:       lr,
	}
}

func (ctr *Ctr) TokenSvc() token.Service {
	return ctr.tokenSvc
}
