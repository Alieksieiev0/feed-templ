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

func (ws *WebServer) Start(
	authServ services.AuthService,
	feedServ services.FeedService,
	userSev services.UserService,
) error {
	ws.app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${method} | ${path} | ${error}\nResponse Body: ${resBody}\n",
	}))

	ws.app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.Dir("./"),
		PathPrefix: "static",
		Browse:     true,
	}))

	ws.app.Get("/", homeHandler(feedServ))
	ws.app.Get("/search", searchPage)
	ws.app.Get("/signup", signupPage)
	ws.app.Get("/signin", signinPage)
	ws.app.Get("/posts", getPostsHandler(feedServ))
	ws.app.Get("/users", getUsersHandler(userSev))
	ws.app.Post("/signup", signupHandler(authServ))
	ws.app.Post("/signin", signinHandler(authServ))
	ws.app.Post("/posts", createPostHandler(feedServ))
	ws.app.Post("/subscribe/:id", subscribeHandler(feedServ))
	ws.app.Use(NotFoundMiddleware)

	return ws.app.Listen(ws.addr)
}
