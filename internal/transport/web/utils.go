package web

import (
	"github.com/Alieksieiev0/feed-templ/internal/view/layout"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func baseWithAuth(c *fiber.Ctx, contents templ.Component) templ.Component {
	return layout.Base(isLoggedIn(c), contents)
}

func isLoggedIn(c *fiber.Ctx) bool {
	return c.Cookies("jwt") != ""
}

func redirect(c *fiber.Ctx, url string, statusCode int) {
	c.Set("HX-Redirect", url)
	c.Status(statusCode)
}

func render(
	c *fiber.Ctx,
	component templ.Component,
	options ...func(*templ.ComponentHandler),
) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}
