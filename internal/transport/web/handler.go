package web

import (
	"github.com/Alieksieiev0/feed-templ/internal/view/core"
	"github.com/Alieksieiev0/feed-templ/internal/view/layout"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

type TemplOption func(*templ.ComponentHandler)

func homeHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return render(c, layout.Base(core.Home("")))
	}
}

func render(
	c *fiber.Ctx,
	component templ.Component,
	options ...TemplOption,
) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}
