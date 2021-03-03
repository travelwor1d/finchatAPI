package httperr

import "github.com/gofiber/fiber/v2"

type httperr struct {
	// Application specific code.
	Code int `json:"code,omitempty"`
	// HTTP status.
	Status int `json:"status"`
	// User readable message.
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func New(code, status int, message string, detail ...interface{}) httperr {
	e := httperr{Code: code, Status: status, Message: message}
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

func (e *httperr) SetDetail(detail interface{}) *httperr {
	switch d := detail.(type) {
	case string:
		e.Detail = d
	case error:
		e.Detail = d.Error()
	}
	return e
}

func (e httperr) Send(c *fiber.Ctx) error {
	return c.Status(e.Status).JSON(e)
}

func (e httperr) Error() string {
	return e.Message
}
