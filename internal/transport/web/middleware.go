package web

import (
	"net/http"

	"github.com/Alieksieiev0/feed-templ/internal/view/core"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func NotFoundMiddleware(c *fiber.Ctx) error {
	return Render(c, core.NotFound(), templ.WithStatus(http.StatusNotFound))
}
