package web

import (
	"github.com/gofiber/fiber/v2"
)

type WebServer struct {
	app *fiber.App
}

func NewWebServer(app *fiber.App) *WebServer {
	return &WebServer{
		app: app,
	}
}

func (ws *WebServer) Start(addr string) error {
	ws.app.Get("/", homeHandler())
	ws.app.Use(NotFoundMiddleware)

	return ws.app.Listen(addr)
}
