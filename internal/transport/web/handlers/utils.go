package handlers

import (
	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/Alieksieiev0/feed-templ/internal/view/layout"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func Render(
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

func baseWithAuth(c *fiber.Ctx, contents templ.Component) templ.Component {
	return layout.Base(isLoggedIn(c), contents)
}

func isLoggedIn(c *fiber.Ctx) bool {
	return c.Cookies("jwt") != ""
}

func setLimitOffsetCookies(c *fiber.Ctx, limit, offset string) {
	c.Cookie(sessionCookie("limit", limit))
	c.Cookie(sessionCookie("offset", offset))
}

func sessionCookie(name, value string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:        name,
		Value:       value,
		SessionOnly: true,
	}
}

func redirectToAuth(c *fiber.Ctx, r *services.Response) {
	c.ClearCookie("jwt", "id", "username", "email")
	r.Redirect(c, "/signin")
}
