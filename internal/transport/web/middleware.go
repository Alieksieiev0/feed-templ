package web

import (
	"net/http"

	"github.com/Alieksieiev0/feed-templ/internal/transport/web/handlers"
	"github.com/Alieksieiev0/feed-templ/internal/view/pages"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func NotFoundMiddleware(c *fiber.Ctx) error {
	return handlers.Render(c, pages.NotFound(), templ.WithStatus(http.StatusNotFound))
}
