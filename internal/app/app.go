package app

import (
	"flag"
	"log"

	"github.com/Alieksieiev0/feed-templ/internal/services"
	"github.com/Alieksieiev0/feed-templ/internal/transport/web"
	"github.com/gofiber/fiber/v2"
)

func Run() {

	var (
		webServerAddr = flag.String("web-server", ":3003", "listen address of web server")
		apiAddr       = flag.String(
			"api",
			"http://localhost:8080",
			"api address",
		)
	)

	s := web.NewWebServer(fiber.New(), *webServerAddr)
	serv := services.NewService(*apiAddr)
	err := s.Start(serv)
	if err != nil {
		log.Fatal(err)
	}
}
