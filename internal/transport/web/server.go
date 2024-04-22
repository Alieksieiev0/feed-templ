package web

import (
	"net/http"

	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/Alieksieiev0/feed-templ/internal/transport/web/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
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

	ws.app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.Dir("./"),
		PathPrefix: "static",
		Browse:     true,
	}))

	ws.pageHandlers(serv)
	ws.feedHandlers(serv)
	ws.userHandlers(serv)
	ws.notificationHandlers(serv)
	ws.authHandlers(serv)
	ws.app.Use(NotFoundMiddleware)
	return ws.app.Listen(ws.addr)
}

func (ws *WebServer) pageHandlers(serv services.FeedService) {
	ws.app.Get("/", handlers.HomePageHandler(serv))
	ws.app.Get("/search", handlers.SearchPage)
	ws.app.Get("/signup", handlers.SignupPage)
	ws.app.Get("/signin", handlers.SigninPage)
}

func (ws *WebServer) feedHandlers(serv services.FeedService) {
	ws.app.Get("/posts", handlers.GetPostsHandler(serv))
	ws.app.Post("/posts", handlers.CreatePostHandler(serv))
	ws.app.Post("/subscribe/:id", handlers.SubscribeHandler(serv))
}

func (ws *WebServer) userHandlers(serv services.UserService) {
	ws.app.Get("/users", handlers.GetUsersHandler(serv))
}

func (ws *WebServer) notificationHandlers(serv services.NotificationServices) {
	ws.app.Get("/notifications", handlers.GetNotificationsHandler(serv))
	ws.app.Get("/notifications/listen", handlers.ListenHandler(serv))
}

func (ws *WebServer) authHandlers(serv services.AuthService) {
	ws.app.Post("/signup", handlers.SignupHandler(serv))
	ws.app.Post("/signin", handlers.SigninHandler(serv))
}
