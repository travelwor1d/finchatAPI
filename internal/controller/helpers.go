package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

var errInternal = httperr.New(codes.Omit, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

func userID(c *fiber.Ctx) (int, *httperr.HTTPErr) {
	uid, ok := c.Locals("uid").(string)
	if !ok || uid == "" {
		return 0, httperr.New(codes.Omit, http.StatusUnauthorized, "failed to retrieve uid from request context")
	}
	id, err := strconv.Atoi(uid)
	if err != nil {
		return 0, httperr.New(codes.Omit, http.StatusInternalServerError, "invalid uid")
	}
	return id, nil
}

func getUserTypes(t string) (string, error) {
	if t == "" {
		return "'GOAT','USER'", nil
	}
	types := strings.Split(t, ",")
	if len(types) < 0 {
		return "", errors.New("invaid `userTypes` format")
	}
	t = ""
	for i, typ := range types {
		typ = strings.ToUpper(strings.TrimSpace(typ))
		if typ != "USER" && typ != "GOAT" {
			return "", fmt.Errorf("invalid type: %s", typ)
		}
		t += "'" + typ + "'"
		if i != len(types)-1 {
			t += ","
		}
	}
	return t, nil
}

func sendSuccess(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"success": true})
}
