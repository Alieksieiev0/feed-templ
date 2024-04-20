package web

import (
	"net/http"

	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/Alieksieiev0/feed-templ/internal/transport/web/handlers"
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

	ws.app.Get("/", handlers.HomeHandler(serv))
	ws.app.Get("/search", handlers.SearchPage)
	ws.app.Get("/signup", handlers.SignupPage)
	ws.app.Get("/signin", handlers.SigninPage)
	ws.app.Get("/posts", handlers.GetPostsHandler(serv))
	ws.app.Get("/users", handlers.GetUsersHandler(serv))
	ws.app.Get("/notifications", handlers.GetNotificationsHandler(serv))
	ws.app.Get("/notifications/listen", handlers.ListenHandler(serv))
	ws.app.Post("/signup", handlers.SignupHandler(serv))
	ws.app.Post("/signin", handlers.SigninHandler(serv))
	ws.app.Post("/posts", handlers.CreatePostHandler(serv))
	ws.app.Post("/subscribe/:id", handlers.SubscribeHandler(serv))
	ws.app.Use(NotFoundMiddleware)

	return ws.app.Listen(ws.addr)
}
