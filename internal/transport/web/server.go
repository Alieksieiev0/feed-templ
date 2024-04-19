package web

import (
	"net/http"

	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type WebServer struct {
	app  *fiber.App
	addr string
}

func NewWebServer(app *fiber.App, addr string) *WebServer {
	return &WebServer{
		app:  app,
		addr: addr,
	}
}

func (ws *WebServer) Start(serv services.Service) error {
	ws.app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${method} | ${path} | ${error}\nResponse Body: ${resBody}\n",
	}))

	ws.app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.Dir("./"),
		PathPrefix: "static",
		Browse:     true,
	}))

	ws.app.Get("/", homeHandler(serv, serv))
	ws.app.Get("/search", searchPage)
	ws.app.Get("/signup", signupPage)
	ws.app.Get("/signin", signinPage)
	ws.app.Get("/posts", getPostsHandler(serv))
	ws.app.Get("/users", getUsersHandler(serv))
	ws.app.Get("/notifications", getNotificationsHandler(serv))
	ws.app.Get("/notifications/listen", listenHandler(serv))
	ws.app.Post("/signup", signupHandler(serv))
	ws.app.Post("/signin", signinHandler(serv))
	ws.app.Post("/posts", createPostHandler(serv))
	ws.app.Post("/subscribe/:id", subscribeHandler(serv))
	ws.app.Use(NotFoundMiddleware)

	return ws.app.Listen(ws.addr)
}
