package controller

import (
	"net/http"
	"path/filepath"

	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/finchatapp/finchat-api/pkg/httperr"
	"github.com/gofiber/fiber/v2"
)

func (ctr *Ctr) UploadProvileAvatar(c *fiber.Ctx) error {
	id, httpErr := userID(c)
	if httpErr != nil {
		return httpErr.Send(c)
	}
	user, err := ctr.store.GetUser(c.Context(), id)
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	file, err := c.FormFile("profileAvatar")
	if err != nil {
		return httperr.New(codes.Omit, http.StatusBadRequest, "cannot get file `profileAvatar`")
	}
	f, err := file.Open()
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	filename := ctr.upload.ProfiveAvatarFileName(user, filepath.Ext(file.Filename))
	err = ctr.upload.UploadProfileAvatar(c.Context(), filename, f)
	if err != nil {
		return errInternal.SetDetail(err).Send(c)
	}
	return c.JSON(fiber.Map{"filename": filename})
}
