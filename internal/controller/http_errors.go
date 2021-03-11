package controller

import (
	"net/http"

	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
)

var errInternal = httperr.New(codes.Omit, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
