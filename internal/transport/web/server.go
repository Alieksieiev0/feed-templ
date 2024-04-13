package web

import (
	"net/http"

	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type WebServer struct {
	app *fiber.App
}

func NewWebServer(app *fiber.App) *WebServer {
	return &WebServer{
		app: app,
	}
}

func (ws *WebServer) Start(addr string, authServ services.AuthService) error {
	ws.app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${method} | ${path} | ${error}\nResponse Body: ${resBody}\n",
	}))

	ws.app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.Dir("./"),
		PathPrefix: "static",
		Browse:     true,
	}))

	ws.app.Get("/", homePage)
	ws.app.Get("/signup", signupPage)
	ws.app.Get("/signin", signinPage)
	ws.app.Post("/signup", signupHandler(authServ))
	ws.app.Post("/signin", signinHandler(authServ))
	ws.app.Use(NotFoundMiddleware)

	return ws.app.Listen(addr)
}
