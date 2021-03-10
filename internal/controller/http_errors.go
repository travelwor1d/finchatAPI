package controller

import (
	"net/http"

	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

var errInternal = httperr.New(codes.Omit, fiber.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
