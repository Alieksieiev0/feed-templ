package services

import "github.com/gofiber/fiber/v2"

type Response struct {
	StatusCode int
	Err        error
}

func (r *Response) SendError(c *fiber.Ctx) error {
	return c.Status(r.StatusCode).Send([]byte("Error: " + r.Err.Error()))
}

func (r *Response) Redirect(c *fiber.Ctx, url string) {
	c.Set("HX-Redirect", url)
	c.Status(r.StatusCode)
}

func NewResponse(statusCode int, err error) *Response {
	return &Response{
		StatusCode: statusCode,
		Err:        err,
	}
}
