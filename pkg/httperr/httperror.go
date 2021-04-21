package httperr

import (
	"net/http"

	"github.com/finchatapp/finchat-api/pkg/codes"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

type HTTPErr struct {
	// Application specific code.
	Code int `json:"code,omitempty"`
	// HTTP status.
	Status int `json:"status"`
	// User readable message.
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

type HTTPValidationErr struct {
	*HTTPErr
}

func NewValidationErr(v *validate.Validation, message string) *HTTPValidationErr {
	var detail string
	if v != nil {
		detail = v.Errors.Error()
	}
	e := New(codes.ValidationError, http.StatusUnprocessableEntity, message).SetDetail(detail)
	return &HTTPValidationErr{e}
}

func New(code, status int, message string, detail ...interface{}) *HTTPErr {
	e := &HTTPErr{Code: code, Status: status, Message: message}
	if len(detail) != 1 {
		return e
	}
	switch d := detail[0].(type) {
	case string:
		e.Detail = d
	case error:
		e.Detail = d.Error()
	}
	return e
}

func (e *HTTPErr) SetDetail(detail interface{}) *HTTPErr {
	switch d := detail.(type) {
	case string:
		e.Detail = d
	case error:
		e.Detail = d.Error()
	}
	return e
}

func (e HTTPErr) Send(c *fiber.Ctx) error {
	return c.Status(e.Status).JSON(e)
}

func (e HTTPErr) Error() string {
	return e.Message
}

func (e HTTPValidationErr) Send(c *fiber.Ctx) error {
	return c.Status(e.Status).JSON(e)
}

func (e HTTPValidationErr) Error() string {
	return e.Message
}
